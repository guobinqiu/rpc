package main

import (
	"fmt"
	"net"

	"github.com/guobinqiu/rpc/rpc"
)

type user struct {
	ID      int64
	Name    string
	Age     int
	Address address
}

type address struct {
	HomeAddr   string
	OfficeAddr string
}

func main() {
	conn, err := net.Dial("tcp", ":3456")
	if err != nil {
		panic(err)
	}

	var u user
	var out *rpc.Out

	// out, err = rpc.Call(conn, "UserService", "GetUserById", []interface{}{1})
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(out.ToInterface(0, u).(user))
	// }

	// out, err = rpc.Call(conn, "UserService", "GetUserByName", []interface{}{"guobin"})
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(out.ToInterface(0, u).(user))
	// }

	// out, err = rpc.Call(conn, "UserService", "Add", []interface{}{1, 2})
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println(out.ToInt(0))
	// 	fmt.Println(out.ToBool(1))
	// }

	u = user{
		Name: "Guobin",
		Age:  100,
		Address: address{
			HomeAddr:   "aaaaa",
			OfficeAddr: "bbbbb",
		},
	}
	out, err = rpc.Call(conn, "UserService", "GrowUp", []interface{}{&u})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(out.Values())
	}
}
