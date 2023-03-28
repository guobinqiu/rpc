package rpc

import (
	"encoding/json"
	"errors"
	"io"
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

func Call(conn io.ReadWriteCloser, serviceName, methodName string, inArgs []any) (*Out, error) {
	var p param
	encoder := json.NewEncoder(conn)
	decoder := json.NewDecoder(conn)

	if err := encoder.Encode(param{
		ServiceName: serviceName,
		MethodName:  methodName,
		InArgs:      inArgs,
	}); err != nil {
		return nil, err
	}

	if err := decoder.Decode(&p); err != nil {
		return nil, err
	}

	if p.Error != "" {
		return nil, errors.New(p.Error)
	}

	return &Out{p.OutArgs}, nil
}

type Server struct {
	services map[string]any
	mu       *sync.Mutex
}

func NewServer() *Server {
	return &Server{services: make(map[string]any), mu: new(sync.Mutex)}
}

func (s *Server) Register(srv any, name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services[name] = srv
}

func (s *Server) HandleConn(conn io.ReadWriteCloser) {
	var p param
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)

	for {
		decoder.Decode(&p)

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

		typeUnmatch := false

		var inValues []reflect.Value

		for i, arg := range p.InArgs {
			if reflect.ValueOf(arg).Type().ConvertibleTo(mtype.In(i + 1)) {
				inValues = append(inValues, reflect.ValueOf(arg).Convert(mtype.In(i+1)))
			} else if reflect.ValueOf(arg).Type().Kind() == reflect.Map {
				t := mtype.In(i + 1).Elem()
				if t.Kind() == reflect.Ptr {
					t = t.Elem()
				}
				vv := reflect.New(t)
				vvv := reflect.Indirect(vv)
				if !canMapToStruct(arg.(map[string]interface{}), vvv) {
					typeUnmatch = true
					break
				}
				inValues = append(inValues, vv)
			} else {
				typeUnmatch = true
				break
			}
		}

		if typeUnmatch {
			p.Error = "参数类型不匹配"
			encoder.Encode(&p)
			continue
		}

		outValues := reflect.ValueOf(s.services[p.ServiceName]).MethodByName(p.MethodName).Call(inValues)
		for _, i := range outValues {
			p.OutArgs = append(p.OutArgs, i.Interface())
		}

		encoder.Encode(&p)
	}
}

func canMapToStruct(arg map[string]interface{}, vv reflect.Value) bool {
	for k, v := range arg {
		structFieldValue := vv.FieldByName(k)
		if !structFieldValue.IsValid() {
			return false
		}
		if reflect.ValueOf(v).Type().ConvertibleTo(structFieldValue.Type()) {
			vvv := reflect.ValueOf(v).Convert(structFieldValue.Type())
			structFieldValue.Set(vvv)
		} else if structFieldValue.Kind() == reflect.Struct {
			return canMapToStruct(v.(map[string]interface{}), structFieldValue)
		}
	}
	return true
}
