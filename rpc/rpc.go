package rpc

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"reflect"
	"sync"
	"time"
)

type param struct {
	ServiceName string
	MethodName  string
	InArgs      []any
	OutArgs     []any
	Error       string
}

type Client struct {
	encoder *json.Encoder
	decoder *json.Decoder
	conn    net.Conn
}

func Dial(network, address string) (*Client, error) {
	conn, err := net.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &Client{
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
		conn:    conn,
	}, nil
}

func (c *Client) Call(serviceName, methodName string, inArgs []any) (*Out, error) {
	var p param

	for _, arg := range inArgs {
		if reflect.TypeOf(arg).Kind() == reflect.Func {
			return nil, errors.New("不支持函数类型")
		}
	}

	if err := c.encoder.Encode(param{
		ServiceName: serviceName,
		MethodName:  methodName,
		InArgs:      inArgs,
	}); err != nil {
		return nil, err
	}

	if err := c.decoder.Decode(&p); err != nil {
		return nil, err
	}

	if p.Error != "" {
		return nil, errors.New(p.Error)
	}

	return &Out{p.OutArgs}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetConn() net.Conn {
	return c.conn
}

type Server struct {
	services map[string]any
	mu       *sync.Mutex
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]any),
		mu:       new(sync.Mutex),
	}
}

func (s *Server) Register(srv any, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[name] = srv
}

func (s *Server) ServeConn(conn net.Conn) {
	defer conn.Close()

	var p param
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		if err := decoder.Decode(&p); err == io.EOF {
			break
		}

		_, ok := s.services[p.ServiceName]
		if !ok {
			p.Error = "服务没找到"
			encoder.Encode(&p)
			continue
		}

		m, b := reflect.TypeOf(s.services[p.ServiceName]).MethodByName(p.MethodName)
		if !b {
			p.Error = "方法没找到"
			encoder.Encode(&p)
			continue
		}

		mtype := m.Type
		if len(p.InArgs) != mtype.NumIn()-1 {
			p.Error = "参数个数不匹配"
			encoder.Encode(&p)
			continue
		}

		inValues, matched := s.match(p, mtype)
		if !matched {
			p.Error = "参数类型不匹配"
			encoder.Encode(&p)
			continue
		}

		outValues := reflect.ValueOf(s.services[p.ServiceName]).MethodByName(p.MethodName).Call(inValues)
		for _, v := range outValues {
			p.OutArgs = append(p.OutArgs, v.Interface())
		}

		encoder.Encode(&p)
	}
}

func (s *Server) match(p param, mtype reflect.Type) ([]reflect.Value, bool) {
	var inValues []reflect.Value
	for i, arg := range p.InArgs {
		t := mtype.In(i + 1)
		if t == reflect.TypeOf(&time.Time{}) {
			t, err := time.Parse(time.RFC3339, arg.(string))
			if err != nil {
				return nil, false
			}
			inValues = append(inValues, reflect.ValueOf(&t))
		} else if t == reflect.TypeOf(time.Time{}) {
			t, err := time.Parse(time.RFC3339, arg.(string))
			if err != nil {
				return nil, false
			}
			inValues = append(inValues, reflect.ValueOf(t))
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			v := reflect.New(t.Elem())
			if !s.mapToStruct(arg.(map[string]any), v.Elem()) {
				return nil, false
			}
			inValues = append(inValues, v)
		} else if t.Kind() == reflect.Struct {
			v := reflect.New(t)
			if !s.mapToStruct(arg.(map[string]any), v.Elem()) {
				return nil, false
			}
			inValues = append(inValues, v.Elem())
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Slice {
			v := reflect.New(reflect.SliceOf(t.Elem().Elem()))
			if !s.copySlice(arg.([]any), v.Elem(), t.Elem().Elem()) {
				return nil, false
			}
			inValues = append(inValues, v)
		} else if t.Kind() == reflect.Slice {
			v := reflect.New(reflect.SliceOf(t.Elem()))
			if !s.copySlice(arg.([]any), v.Elem(), t.Elem()) {
				return nil, false
			}
			inValues = append(inValues, v.Elem())
		} else if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Array {
			v := reflect.New(reflect.ArrayOf(t.Elem().Len(), t.Elem().Elem()))
			if !s.copyArray(arg.([]any), v.Elem(), t.Elem().Elem()) {
				return nil, false
			}
			inValues = append(inValues, v)
		} else if t.Kind() == reflect.Array {
			v := reflect.New(reflect.ArrayOf(t.Len(), t.Elem()))
			if !s.copyArray(arg.([]any), v.Elem(), t.Elem()) {
				return nil, false
			}
			inValues = append(inValues, v.Elem())
		} else if reflect.ValueOf(arg).Type().ConvertibleTo(t) {
			inValues = append(inValues, reflect.ValueOf(arg).Convert(t))
		} else {
			return nil, false
		}
	}
	return inValues, true
}

func (s *Server) mapToStruct(arg map[string]any, v reflect.Value) bool {
	for key, value := range arg {
		structFieldValue := v.FieldByName(key)
		if !structFieldValue.IsValid() {
			return false
		}
		if structFieldValue.Kind() == reflect.Struct {
			if !s.mapToStruct(value.(map[string]any), structFieldValue) {
				return false
			}
		} else if structFieldValue.Kind() == reflect.Slice {
			if value == nil {
				value = make([]any, 0)
			}
			if !s.copySlice(value.([]any), structFieldValue, structFieldValue.Type().Elem()) {
				return false
			}
		} else if structFieldValue.Kind() == reflect.Array {
			if !s.copyArray(value.([]any), structFieldValue, structFieldValue.Type().Elem()) {
				return false
			}
		} else if structFieldValue.Kind() == reflect.Ptr && structFieldValue.Type().Elem().Kind() == reflect.Slice {
			if value == nil {
				value = make([]any, 0)
			}
			vv := reflect.New(reflect.SliceOf(structFieldValue.Type().Elem().Elem()))
			if !s.copySlice(value.([]any), vv.Elem(), structFieldValue.Type().Elem().Elem()) {
				return false
			}
			structFieldValue.Set(vv)
		} else if structFieldValue.Kind() == reflect.Ptr && structFieldValue.Type().Elem().Kind() == reflect.Array {
			if value == nil {
				value = make([]any, 0)
			}
			vv := reflect.New(reflect.ArrayOf(structFieldValue.Type().Elem().Len(), structFieldValue.Type().Elem().Elem()))
			if !s.copyArray(value.([]any), vv.Elem(), structFieldValue.Type().Elem().Elem()) {
				return false
			}
			structFieldValue.Set(vv)
		} else if value != nil && reflect.ValueOf(value).Type().ConvertibleTo(structFieldValue.Type()) {
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		} else {
			return false
		}
	}
	return true
}

func (s *Server) copySlice(arg []any, v reflect.Value, t reflect.Type) bool {
	for _, value := range arg {
		if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			vv := reflect.New(t.Elem())
			if !s.mapToStruct(value.(map[string]any), vv.Elem()) {
				return false
			}
			v.Set(reflect.Append(v, vv))
		} else if t.Kind() == reflect.Struct {
			vv := reflect.New(t)
			if !s.mapToStruct(value.(map[string]any), vv.Elem()) {
				return false
			}
			v.Set(reflect.Append(v, vv.Elem()))
		} else if value != nil && reflect.ValueOf(value).Type().ConvertibleTo(t) {
			v.Set(reflect.Append(v, reflect.ValueOf(value).Convert(t)))
		} else {
			return false
		}
	}
	return true
}

func (s *Server) copyArray(arg []any, v reflect.Value, t reflect.Type) bool {
	for i, value := range arg {
		if t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct {
			vv := reflect.New(t.Elem())
			if !s.mapToStruct(value.(map[string]any), vv.Elem()) {
				return false
			}
			v.Index(i).Set(vv)
		} else if t.Kind() == reflect.Struct {
			vv := reflect.New(t)
			if !s.mapToStruct(value.(map[string]any), vv.Elem()) {
				return false
			}
			v.Index(i).Set(vv.Elem())
		} else if value != nil && reflect.ValueOf(value).Type().ConvertibleTo(t) {
			v.Index(i).Set(reflect.ValueOf(value).Convert(t))
		} else {
			return false
		}
	}
	return true
}
