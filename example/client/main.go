package main

import (
	"fmt"

	"github.com/guobinqiu/rpc/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":3456")
	if err != nil {
		panic(err)
	}

	out, _ := client.Call("UserService", "Add", []interface{}{1, 2})
	fmt.Println(out.ToInt(0))

	client.Close()
}
