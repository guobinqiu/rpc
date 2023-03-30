package rpc

import (
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

func (s *Userservice) Add(a int, b int) int {
	return a + b
}

func (s *Userservice) GrowUpPointer(u *user) *user {
	u.Age += 1
	return u
}

func (s *Userservice) GrowUpStruct(u user) user {
	u.Age += 1
	return u
}

func TestGetUserById(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:7890")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	out, _ := client.Call("UserService", "GetUserById", []interface{}{1})
	t.Log(out)

	client.Close()
	listener.Close()
}

func TestGetUserByIdInvalidNum(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	_, err := client.Call("UserService", "GetUserById", []interface{}{1, 2})
	t.Log(err)

	client.Close()
	listener.Close()
}

func TestGetUserByIdInvalidType(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	_, err := client.Call("UserService", "GetUserById", []interface{}{"1"})
	t.Log(err)

	client.Close()
	listener.Close()
}

func TestServiceNotFound(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	_, err := client.Call("UserServicee", "GetUserById", []interface{}{"1"})
	t.Log(err)

	client.Close()
	listener.Close()
}

func TestMethodNotFound(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	_, err := client.Call("UserService", "GetUserByIdd", []interface{}{"1"})
	t.Log(err)

	client.Close()
	listener.Close()
}

func TestAdd(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", listener.Addr().String())
	out, _ := client.Call("UserService", "Add", []interface{}{1, 2})
	t.Log(out)

	client.Close()
	listener.Close()
}

func TestGrowUpPointer(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	u := user{
		Name: "Guobin",
		Age:  100,
		Address: address{
			HomeAddr:   "aaaaa",
			OfficeAddr: "bbbbb",
		},
	}

	client, _ := Dial("tcp", listener.Addr().String())
	out, _ := client.Call("UserService", "GrowUpPointer", []interface{}{&u})
	t.Log(out)

	client.Close()
	listener.Close()
}

func TestGrowUpStruct(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	listener, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	u := user{
		Name: "Guobin",
		Age:  100,
		Address: address{
			HomeAddr:   "aaaaa",
			OfficeAddr: "bbbbb",
		},
	}

	client, _ := Dial("tcp", listener.Addr().String())
	out, _ := client.Call("UserService", "GrowUpStruct", []interface{}{u})
	t.Log(out)

	client.Close()
	listener.Close()
}
