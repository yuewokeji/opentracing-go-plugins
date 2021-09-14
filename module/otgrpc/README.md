# otgrpc

----

OpenTracing plugin for [grpc](https://github.com/grpc/grpc).

## Usage

---

### Server

```go
package main

import (
	"github.com/yuewokeji/opentracing-go-plugins/module/otgrpc"
	"google.golang.org/grpc"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := grpc.NewServer(
		grpc.UnaryInterceptor(otgrpc.NewUnaryServerInterceptor()),
	)

	if err := server.Serve(lis); err != nil {
		panic(err)
	}
}

```

### Client

```go
package main

import (
	"github.com/yuewokeji/opentracing-go-plugins/module/otgrpc"
	"google.golang.org/grpc"
)

func main() {
	// ...
	conn, err := grpc.Dial(":8080", grpc.WithUnaryInterceptor(
		otgrpc.NewUnaryClientInterceptor(),
	))

	// do something
}

```