package rpc

import (
	"encoding/json"
	"reflect"
)

type Out struct {
	outArgs []any
}

func (o *Out) deconstruct(t reflect.Type, v any) reflect.Value {
	vv := reflect.New(t)
	b, _ := json.Marshal(v)
	json.Unmarshal(b, vv.Interface())
	return vv
}

func (o *Out) ToInt(index int) int {
	t := reflect.TypeOf(1)
	v := o.deconstruct(t, o.outArgs[index])
	return int(v.Elem().Int())
}

func (o *Out) ToBool(index int) bool {
	t := reflect.TypeOf(true)
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().Bool()
}

func (o *Out) ToInterface(index int, customType any) any {
	t := reflect.TypeOf(customType)
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().Interface()
}

func (o *Out) Len() int {
	return len(o.outArgs)
}

func (o *Out) Values() []any {
	return o.outArgs
}
