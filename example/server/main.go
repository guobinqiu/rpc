package main

import (
	"net"

	"github.com/guobinqiu/rpc/rpc"
)

type user struct {
	ID   int64
	Name string
	Age  int
}

type Userservice struct{}

func (u *Userservice) GetUserById(id int64) *user {
	return &user{
		ID:   id,
		Name: "Jack",
		Age:  100,
	}
}

func (u *Userservice) GetUserByName(name string) *user {
	return &user{
		ID:   2,
		Name: name,
		Age:  100,
	}
}

func (u *Userservice) Add(a int, b int) (int, bool) {
	return a + b, true
}

func main() {
	server := rpc.NewServer()
	server.Register(new(Userservice), "UserService")

	l, _ := net.Listen("tcp", ":3456")
	for {
		conn, _ := l.Accept()
		go server.HandleConn(conn)
	}
}
