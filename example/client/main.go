package main

import (
	"fmt"
	"net"

	"github.com/guobinqiu/rpc/rpc"
)

type user struct {
	ID   int64
	Name string
	Age  int
}

func main() {
	conn, err := net.Dial("tcp", ":3456")
	if err != nil {
		panic(err)
	}

	var u user
	var out *rpc.Out

	out = rpc.Call(conn, "UserService", "GetUserById", []interface{}{1})
	fmt.Println(out.Interface(0, u).(user))

	out = rpc.Call(conn, "UserService", "GetUserByName", []interface{}{"guobin"})
	fmt.Println(out)

	out = rpc.Call(conn, "UserService", "Add", []interface{}{1, 2})
	fmt.Println(out)
}
