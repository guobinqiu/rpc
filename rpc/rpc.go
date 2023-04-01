package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"sync"
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

		fmt.Println(p.InArgs)

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
		if reflect.ValueOf(arg).Type().ConvertibleTo(t) {
			inValues = append(inValues, reflect.ValueOf(arg).Convert(t))
		} else if reflect.ValueOf(arg).Type().Kind() == reflect.Map {
			tt := t
			if t.Kind() == reflect.Ptr {
				tt = t.Elem()
			}
			v := reflect.New(tt)
			if !s.mapToStruct(arg.(map[string]any), reflect.Indirect(v)) {
				return nil, false
			}
			if t.Kind() == reflect.Struct {
				v = v.Elem()
			}
			inValues = append(inValues, v)
		} else if reflect.ValueOf(arg).Type().Kind() == reflect.Slice {
			tt := t
			if t.Kind() == reflect.Ptr {
				tt = t.Elem()
			}
			v := reflect.New(tt)
			if !s.copySlice(arg.([]any), reflect.Indirect(v), tt.Elem()) {
				return nil, false
			}
			if t.Kind() != reflect.Pointer {
				v = v.Elem()
			}
			inValues = append(inValues, v)
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
		if reflect.ValueOf(value).Type().ConvertibleTo(structFieldValue.Type()) {
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		} else if structFieldValue.Kind() == reflect.Struct {
			return s.mapToStruct(value.(map[string]any), structFieldValue)
		} else {
			return false
		}
	}
	return true
}

func (s *Server) mapToStruct2(arg map[string]any, v reflect.Value) bool {
	for key, value := range arg {
		structFieldValue := v.FieldByName(key)
		if !structFieldValue.IsValid() {
			return false
		}
		if reflect.ValueOf(value).Type().ConvertibleTo(structFieldValue.Type()) {
			structFieldValue.Set(reflect.ValueOf(value).Convert(structFieldValue.Type()))
		} else if structFieldValue.Kind() == reflect.Struct {
			return s.mapToStruct(value.(map[string]any), structFieldValue)
		} else {
			return false
		}
	}
	return true
}

func (s *Server) copySlice(arg []any, v reflect.Value, t reflect.Type) bool {
	for _, value := range arg {
		if reflect.ValueOf(value).Type().ConvertibleTo(t) {
			v.Set(reflect.Append(v, reflect.ValueOf(value).Convert(t)))
		} else if reflect.ValueOf(value).Kind() == reflect.Map {
			tt := t
			if t.Kind() == reflect.Ptr {
				tt = t.Elem()
			}
			v2 := reflect.New(tt)
			s.mapToStruct(value.(map[string]any), reflect.Indirect(v2))
			if t.Kind() != reflect.Pointer {
				v2 = v2.Elem()
			}
			v.Set(reflect.Append(v, v2))
		} else {
			return false
		}
	}
	return true
}
