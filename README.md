RPC
---
不喜欢go内置的rpc写法，这个rpc希望达到的目的是你不用为了符合go内置的rpc规范而修改你现在的代码，普通的方法就可以同时当作rpc的方法被调用。

### Run test

```
go clean -testcache && go test -v ./rpc
```
