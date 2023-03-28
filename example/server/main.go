package main

import (
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

type Userservice struct{}

func (s *Userservice) GetUserById(id int64) *user {
	return &user{
		ID:   id,
		Name: "Jack",
		Age:  100,
	}
}

func (s *Userservice) GetUserByName(name string) *user {
	return &user{
		ID:   2,
		Name: name,
		Age:  100,
	}
}

func (s *Userservice) Add(a int, b int) (int, bool) {
	return a + b, true
}

func (s *Userservice) GrowUp(u *user) *user {
	u.Age += 1
	return u
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
