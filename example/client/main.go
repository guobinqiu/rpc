package main

import (
	"fmt"

	"github.com/guobinqiu/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":3456")
	if err != nil {
		panic(err)
	}

	out, _ := client.Call("UserService", "Add", []interface{}{1, 2})
	fmt.Println(out.Get(0))

	client.Close()
}
