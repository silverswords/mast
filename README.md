# mast

[![Go Report Card](https://goreportcard.com/badge/github.com/silverswords/mast)](https://goreportcard.com/report/github.com/silverswords/mast)

[English](https://github.com/silverswords/mast/Readme) [Chinese](https://github.com/silverswords/mast/zh-cn)

mast is a builder for rpc client and server By use options to reduce complexity and power new people use gprc.

The core of the Service Mesh model is based on the basic principle of separating the client SDK and running it as a Proxy independent process; the goal is to sink the various capabilities that originally existed in the SDK to reduce the burden on the application to help the application cloud native.

mast is used for beginners to quickly obtain solutions for the grpc production environment and quickly learn to land the use of grpc
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

- [x] Client-side Dial 

- [x] Server-side register service and start

### features
- [ ] autoauth - interceptor
    - [ ] jwt
    - [ ] DIY
    - [ ] grpc_auth

- [ ] TLS - Client & Server
    - [x] DialTLS
    - [ ] use grpc.WithTransportCredentials(tlsCreds)
    - [ ] ServerCert  -> tlsCreds :only provide this way
    - [ ] Cert expired & auto-acme

- [ ] metrics 
    - [ ] grpc_prometheus
    - [ ] otgrpc
    
- [ ] Trace
    - [ ] grpc_opentracing
    
- [ ] logging
    - [ ] grpc_logrus
    - [ ] grpc_zap
    - [ ] grpc_ctxtags

- Client-side
- [ ] grpc_retry - Client-side
    - [ ] retry by errcode
    
- [x] timeout-dial  - Client-side 

- [x] Compress

- [ ] loadbalance

- Server-side
- [ ] panic-recovery : grpc_recovery

- [ ] grpc_validator

- [ ] ratelimit 

- [ ] transport

- [ ] healthy

- [ ] registery 
    - [ ]consul 

- [ ] telemetry

- [ ] route or  GRPCGateWay

### Stage

1. Complete Features
2. Could Build from Config file
3. Web Control and Automatic
