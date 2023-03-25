package rpc

import (
	"encoding/json"
	"reflect"
)

type Out struct {
	OutArgs []any
}

func (o *Out) deconstruct(t reflect.Type, v any) reflect.Value {
	vv := reflect.New(t)
	b, _ := json.Marshal(v)
	json.Unmarshal(b, vv.Interface())
	return vv
}

func (o *Out) Int(index int) int {
	t := reflect.TypeOf(1)
	v := o.deconstruct(t, o.OutArgs[index])
	return int(v.Elem().Int())
}

func (o *Out) Bool(index int) bool {
	t := reflect.TypeOf(true)
	v := o.deconstruct(t, o.OutArgs[index])
	return v.Elem().Bool()
}

func (o *Out) Interface(index int, customType any) any {
	t := reflect.TypeOf(customType)
	v := o.deconstruct(t, o.OutArgs[index])
	return v.Elem().Interface()
}

func (o *Out) Len() int {
	return len(o.OutArgs)
}
