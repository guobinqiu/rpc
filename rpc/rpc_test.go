package rpc

import (
	"fmt"
	"net"
	"testing"
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

func (s *Userservice) GrowUp2(u user) user {
	u.Age += 1
	return u
}

func TestGetUserById(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, _ := listener.Accept()
			go server.ServeConn(conn)
		}
	}()

	conn, _ := net.Dial("tcp", listener.Addr().String())
	out, err := Call(conn, "UserService", "GetUserById", []interface{}{1})
	t.Log(err == nil)
	fmt.Println(out)

	conn.Close()
}

func TestGetUserByIdInvalidNum(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, _ := listener.Accept()
			go server.ServeConn(conn)
		}
	}()

	conn, _ := net.Dial("tcp", listener.Addr().String())
	_, err := Call(conn, "UserService", "GetUserById", []interface{}{1, 2})
	t.Log(err != nil)
	fmt.Println(err)

	conn.Close()
}

func TestGetUserByIdInvalidType(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, _ := listener.Accept()
			go server.ServeConn(conn)
		}
	}()

	conn, _ := net.Dial("tcp", listener.Addr().String())
	_, err := Call(conn, "UserService", "GetUserById", []interface{}{"1"})
	t.Log(err != nil)
	fmt.Println(err)

	conn.Close()
}