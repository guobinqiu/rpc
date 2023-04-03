package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"testing"
)

type user struct {
	ID           int64
	Name         string
	Age          int
	Address      address
	HobbiesSlice []string
	HobbiesArr   [3]string
	SonsStruct   []user
	SonsPtr      []*user
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

func (s *Userservice) SumUserAgePointer(users []*user) int {
	sum := 0
	for _, u := range users {
		sum += u.Age
	}
	return sum
}

func (s *Userservice) SumUserAgeStruct(users []user) int {
	sum := 0
	for _, u := range users {
		sum += u.Age
	}
	return sum
}

func (s *Userservice) TestFunc(f func(int, int) int) {
	fmt.Println("a+b=", f(1, 2))
}

func (s *Userservice) TestArrStruct(users [3]user) int {
	sum := 0
	for _, u := range users {
		sum += u.Age
	}
	return sum
}

func (s *Userservice) TestArrPointer(users [3]*user) int {
	sum := 0
	for _, u := range users {
		sum += u.Age
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

func TestSumUserAgePointer(t *testing.T) {
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
	out, err := client.Call("UserService", "SumUserAgePointer", []interface{}{[]*user{
		{Age: 1},
		{Age: 2},
		{Age: 3},
	},
	})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestSumUserAgeStruct(t *testing.T) {
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
	out, err := client.Call("UserService", "SumUserAgeStruct", []interface{}{[]*user{
		{Age: 1},
		{Age: 2},
		{Age: 3},
	},
	})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestTestFunc(t *testing.T) {
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
	out, err := client.Call("UserService", "TestFunc", []interface{}{func(a, b int) int {
		return a + b
	}})

	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestTestArrStruct(t *testing.T) {
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
	out, err := client.Call("UserService", "TestArrStruct", []interface{}{[3]user{
		{Age: 1},
		{Age: 2},
		{Age: 3},
	},
	})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestTestArrPointer(t *testing.T) {
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
	out, err := client.Call("UserService", "TestArrPointer", []interface{}{[3]*user{
		{Age: 1},
		{Age: 2},
		{Age: 3},
	},
	})
	t.Log(err)
	t.Log(out)

	client.Close()
	l.Close()
}

func TestArrInsideStruct(t *testing.T) {
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
		HobbiesSlice: []string{
			"football",
			"basketball",
		},
		HobbiesArr: [3]string{
			"football",
			"basketball",
		},
		SonsPtr: []*user{
			{Name: "a"},
			{Name: "b"},
		},
		SonsStruct: []user{
			{Name: "aa"},
			{Name: "bb"},
		},
	}

	client, _ := Dial("tcp", l.Addr().String())
	out, _ := client.Call("UserService", "GrowUpPointer", []interface{}{&u})

	b, _ := json.Marshal(out.Get(0))
	json.Unmarshal(b, &u)

	t.Log(u.SonsPtr[0].Name)
	t.Log(u.SonsStruct[0].Name)

	client.Close()
	l.Close()
}
