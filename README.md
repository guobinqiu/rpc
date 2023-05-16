# RPC

不喜欢go内置的rpc写法，这个rpc希望达到的目的是你不用为了符合go内置的rpc规范而修改你现在的代码，普通的方法就可以同时当作rpc的方法被调用。

## Usage

server

```
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
```

client

```
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
```

## Test

```
go clean -testcache && go test -v .
```

## License

MIT
