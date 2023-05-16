package main

import (
	"net"

	"github.com/guobinqiu/rpc"
)

type Userservice struct{}

func (s *Userservice) Add(a, b int) int {
	return a + b
}

func main() {
	server := rpc.NewServer()
	server.Register(new(Userservice), "UserService")

	listener, err := net.Listen("tcp", ":3456")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		go server.ServeConn(conn)
	}
}
