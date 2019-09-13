# mast

[![Go Report Card](https://goreportcard.com/badge/github.com/silverswords/mast)](https://goreportcard.com/report/github.com/silverswords/mast)

[English](https://github.com/silverswords/mast/Readme) [Chinese](https://github.com/silverswords/mast/zh-cn)

mast is a builder for rpc client and server By use options to reduce complexity and power new people use gprc.

Service Mesh 模式的核心，其基本原理在于将客户端 SDK 剥离，以 Proxy 独立进程运行；目标是将原来存在于 SDK 中的各种能力下沉，为应用减负，以帮助应用云原生化。

## Usage

### Server
```go
import (
	"context"
	"log"

	"github.com/silverswords/mast/mastgrpc"
	pb "github.com/silverswords/mast/example/helloworld"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	b := mastgrpc.DefaultGRPCBuildOptions()
	
	s := b.Server()

	s.Prepare(pb.RegisterGreeterServer, &server{})

	if err := s.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
```

### Client.go
```go
import (
	"context"
	"log"
	"os"
	"time"

	"github.com/silverswords/mast/mastgrpc"

	pb "github.com/silverswords/mast/example/helloworld"
)

const (
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	b := mastgrpc.DefaultGRPCBuildOptions()

	conn, err := b.ClientConn()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
```

1. 决定 grpc mast 的基底
2. 需要添加哪些特性开始扩展 Mast 
3. 构建和 istio  promethues 的集成

- [x] Client-side Dial 

- [x] Server-side register service and start

### features
- [x] timeout-dial  - Client-side 

- [ ] Compress

- [ ] autoauth - interceptor
    - [ ] jwt
    - [ ] DIY

- [ ] TLS - Client & Server
    - [x] DialTLS
    - [ ] use grpc.WithTransportCredentials(tlsCreds)
    - [ ] ServerCert  -> tlsCreds :only provide this way
    - [ ] Cert expired & auto-acme

- [ ] metrics

- [ ] logging

- [ ] grpc_retry - Client-side
    - [ ] retry by errcode

- [ ] panic-recovery

- [ ] ratelimit 

- [ ] meatedata

- [ ] transport

- [ ] Trace

- [ ] healthy

- [ ] loadbalance

- [ ] registery 
    - [ ]consul 

- [ ] CI/CD

- [ ] telemetry

- [ ] route or  GRPCGateWay

- [ ] MQ

- [ ] Cache


