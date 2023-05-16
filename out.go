package rpc

type Out struct {
	outArgs []any
}

func (o *Out) Len() int {
	return len(o.outArgs)
}

func (o *Out) Get(index int) any {
	return o.outArgs[index]
}
