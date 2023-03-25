package rpc

import (
	"encoding/json"
	"io"
	"reflect"
)

type param struct {
	ServiceName string
	MethodName  string
	InArgs      []interface{}
	OutArgs     []interface{}
}

func Call(conn io.ReadWriteCloser, serviceName, methodName string, inArgs []interface{}) *Out {
	encoder := json.NewEncoder(conn)
	encoder.Encode(param{
		ServiceName: serviceName,
		MethodName:  methodName,
		InArgs:      inArgs,
	})

	var p param
	decoder := json.NewDecoder(conn)
	decoder.Decode(&p)
	return &Out{p.OutArgs}
}

type Server struct {
	services map[string]interface{}
}

func NewServer() *Server {
	return &Server{services: make(map[string]interface{})}
}

func (s *Server) Register(srv interface{}, name string) {
	s.services[name] = srv
}

func (s *Server) HandleConn(conn io.ReadWriteCloser) {
	var p param
	decoder := json.NewDecoder(conn)
	encoder := json.NewEncoder(conn)
	for {
		decoder.Decode(&p)

		m, _ := reflect.TypeOf(s.services[p.ServiceName]).MethodByName(p.MethodName)
		mtype := m.Type
		if len(p.InArgs) != mtype.NumIn()-1 {
			p.OutArgs = append(p.OutArgs, "参数个数不匹配")
			encoder.Encode(&p)
			return
		}

		var inValues []reflect.Value
		for i, arg := range p.InArgs {
			if !reflect.ValueOf(arg).Type().ConvertibleTo(mtype.In(i + 1)) {
				p.OutArgs = append(p.OutArgs, "参数类型不匹配")
				encoder.Encode(&p)
				return
			}
			inValues = append(inValues, reflect.ValueOf(arg).Convert(mtype.In(i+1)))
		}

		outValues := reflect.ValueOf(s.services[p.ServiceName]).MethodByName(p.MethodName).Call(inValues)
		for _, i := range outValues {
			p.OutArgs = append(p.OutArgs, i.Interface())
		}

		encoder.Encode(&p)
	}
}
