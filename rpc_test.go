package rpc

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

type user struct {
	ID             int64
	Name           string
	Age            int
	Address        address
	HobbiesSlice   []string
	HobbiesArr     [3]string
	SliceStruct    []user
	SlicePtrStruct []*user
	PtrSliceStruct *[]user
	PtrArrayStruct *[3]user
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

func (s *Userservice) TestTime(t time.Time) time.Time {
	return t.Add(time.Hour)
}

func (s *Userservice) TestTimePtr(t *time.Time) *time.Time {
	tt := t.Add(time.Hour)
	return &tt
}

func (s *Userservice) EmptyIn() string {
	return "guobin"
}

func (s *Userservice) EmptyOut(name string) {
}

func (s *Userservice) EmptyInAndOut() {
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
	out, err := client.Call("UserService", "GetUserById", []interface{}{1})
	if err != nil {
		t.Error(err)
	}
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
	out, err := client.Call("UserService", "GetUserById", []interface{}{1, 2})
	if err == nil {
		t.Error(out)
	}
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
	out, err := client.Call("UserService", "GetUserById", []interface{}{"1"})
	if err == nil {
		t.Error(out)
	}
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
	out, err := client.Call("UserServicee", "GetUserById", []interface{}{"1"})
	if err == nil {
		t.Error(out)
	}
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
	out, err := client.Call("UserService", "GetUserByIdd", []interface{}{"1"})
	if err == nil {
		t.Error(out)
	}
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
	out, err := client.Call("UserService", "Add", []interface{}{1, 2})
	if err != nil {
		t.Error(err)
	}
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
	out, err := client.Call("UserService", "GrowUpPointer", []interface{}{&u})
	if err != nil {
		t.Error(err)
	}
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
	out, err := client.Call("UserService", "GrowUpStruct", []interface{}{u})
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
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
	if err == nil {
		t.Log(out)
	}
	t.Log(err)

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
	if err != nil {
		t.Error(err)
	}
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
	if err != nil {
		t.Error(err)
	}
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
		SliceStruct: []user{
			{Name: "c"},
			{Name: "d"},
		},
		SlicePtrStruct: []*user{
			{Name: "a"},
			{Name: "b"},
		},
		PtrSliceStruct: &[]user{
			{Name: "eee"},
			{Name: "fff"},
		},
		PtrArrayStruct: &[3]user{
			{Name: "x"},
			{Name: "y"},
			{Name: "z"},
		},
	}

	client, _ := Dial("tcp", l.Addr().String())
	out, err := client.Call("UserService", "GrowUpPointer", []interface{}{&u})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	b, _ := json.Marshal(out.Get(0))
	json.Unmarshal(b, &u)

	t.Log(u.SliceStruct[0].Name)
	t.Log(u.SlicePtrStruct[0].Name)
	t.Log((*u.PtrSliceStruct)[0].Name)
	t.Log((*u.PtrArrayStruct)[0].Name)

	client.Close()
	l.Close()
}

func TestTestTime(t *testing.T) {
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
	out, err := client.Call("UserService", "TestTime", []interface{}{time.Now()})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	client.Close()
	l.Close()
}

func TestTestTimePtr(t *testing.T) {
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
	tt := time.Now()
	out, err := client.Call("UserService", "TestTimePtr", []interface{}{&tt})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	client.Close()
	l.Close()
}

func TestEmptyIn(t *testing.T) {
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
	out, err := client.Call("UserService", "EmptyIn", []interface{}{})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	client.Close()
	l.Close()
}

func TestEmptyOut(t *testing.T) {
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
	out, err := client.Call("UserService", "EmptyOut", []interface{}{"guobin"})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	client.Close()
	l.Close()
}

func TestEmptyInAndOut(t *testing.T) {
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
	out, err := client.Call("UserService", "EmptyInAndOut", []interface{}{})
	if err != nil {
		t.Error(err)
	}
	t.Log(out)

	client.Close()
	l.Close()
}

func TestConcurrency(t *testing.T) {
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
		Age:  0,
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			client, _ := Dial("tcp", l.Addr().String())
			out, err := client.Call("UserService", "GrowUpStruct", []interface{}{u})
			if err != nil {
				t.Error(err)
			}
			t.Log(out)
			client.Close()
			wg.Done()
		}()
	}
	wg.Wait()

	l.Close()
}
