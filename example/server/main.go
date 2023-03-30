package main

import (
	"net"

	"github.com/guobinqiu/rpc/rpc"
)

type Userservice struct{}

func (s *Userservice) Add(a int, b int) (int, bool) {
	return a + b, true
}

func main() {
	server := rpc.NewServer()
	server.Register(new(Userservice), "UserService")

	l, err := net.Listen("tcp", ":3456")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			continue
		}
		go server.ServeConn(conn)
	}
}
