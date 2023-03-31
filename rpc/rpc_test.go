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
		Name: "Guobin",
		Age:  40,
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

func (s *Userservice) Sum(nums []int) int {
	sum := 0
	for _, num := range nums {
		sum += num
	}
	return sum
}

func (s *Userservice) SumPointer(nums *[]int) int {
	sum := 0
	for _, num := range *nums {
		sum += num
	}
	return sum
}

func TestGetUserById(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	out, _ := client.Call("UserService", "GetUserById", []interface{}{1})
	t.Log(out)

	client.Close()
	l.Close()
}

func TestGetUserByIdInvalidNum(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	_, err := client.Call("UserService", "GetUserById", []interface{}{1, 2})
	t.Log(err)

	client.Close()
	l.Close()
}

func TestGetUserByIdInvalidType(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	_, err := client.Call("UserService", "GetUserById", []interface{}{"1"})
	t.Log(err)

	client.Close()
	l.Close()
}

func TestServiceNotFound(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	_, err := client.Call("UserServicee", "GetUserById", []interface{}{"1"})
	t.Log(err)

	client.Close()
	l.Close()
}

func TestMethodNotFound(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	_, err := client.Call("UserService", "GetUserByIdd", []interface{}{"1"})
	t.Log(err)

	client.Close()
	l.Close()
}

func TestAdd(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	out, _ := client.Call("UserService", "Add", []interface{}{1, 2})
	t.Log(out)

	client.Close()
	l.Close()
}

func TestGrowUpPointer(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	u := user{
		Name: "Guobin",
		Age:  40,
		Address: address{
			HomeAddr:   "aaaaa",
			OfficeAddr: "bbbbb",
		},
	}

	client, _ := Dial("tcp", l.Addr().String())
	out, _ := client.Call("UserService", "GrowUpPointer", []interface{}{&u})
	t.Log(out)

	client.Close()
	l.Close()
}

func TestGrowUpStruct(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	u := user{
		Name: "Guobin",
		Age:  40,
		Address: address{
			HomeAddr:   "aaaaa",
			OfficeAddr: "bbbbb",
		},
	}

	client, _ := Dial("tcp", l.Addr().String())
	out, _ := client.Call("UserService", "GrowUpStruct", []interface{}{u})
	t.Log(out)

	client.Close()
	l.Close()
}

func TestSum(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	out, err := client.Call("UserService", "Sum", []interface{}{[]int{1, 2, 3}})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestSumPointer(t *testing.T) {
	server := NewServer()
	server.Register(new(Userservice), "UserService")
	l, _ := net.Listen("tcp", "127.0.0.1:0")
	go func() {
		for {
			conn, err := l.Accept()
			if err != nil {
				return
			}
			go server.ServeConn(conn)
		}
	}()

	client, _ := Dial("tcp", l.Addr().String())
	out, err := client.Call("UserService", "SumPointer", []interface{}{&[]int{1, 2, 3}})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}
