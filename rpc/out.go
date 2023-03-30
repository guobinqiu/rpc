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
	err := json.Unmarshal(b, vv.Interface())
	if err != nil {
		panic("deconstruct err:" + err.Error())
	}
	return vv
}

func (o *Out) ToInt64(index int) int64 {
	t := reflect.TypeOf(0)
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().Int()
}

func (o *Out) ToBool(index int) bool {
	t := reflect.TypeOf(true)
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().Bool()
}

func (o *Out) ToFloat64(index int) float64 {
	t := reflect.TypeOf(0.0)
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().Float()
}

func (o *Out) ToString(index int) string {
	t := reflect.TypeOf("")
	v := o.deconstruct(t, o.outArgs[index])
	return v.Elem().String()
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
