# mast

[![Go Report Card](https://goreportcard.com/badge/github.com/silverswords/mast)](https://goreportcard.com/report/github.com/silverswords/mast)

mast is a builder for rpc client and server By use options to reduce complexity and power new people use gprc.

Service Mesh 模式的核心，其基本原理在于将客户端 SDK 剥离，以 Proxy 独立进程运行；目标是将原来存在于 SDK 中的各种能力下沉，为应用减负，以帮助应用云原生化。

## Usage

### Server
```go
import (
	"context"
	"log"

	"github.com/silverswords/mast/mastgrpc"
	pb example
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

	pb "github.com/silverswords/mast/unittest/helloworld"
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

第一是  决定 grpc mast 的基底
第二是 需要添加哪些特性开始扩展 Mast 
第三 构建和 istio  promethues 的集成
添加特性
思考一下下
- [ ] 超时dial 控制

- [ ] 压缩选项

- [ ] 认证选项

- [ ] TLS 加密

- [ ] 自动重试

- [ ] panic 恢复

- [ ] 流量控制：ratelimit 限制

- [ ] 统计

- [ ] 系统监控和日志

- [ ] 转发 transport

- [ ] 追踪 Trace

- [ ] 心跳健康检查

- [ ] 负载均衡

- [ ] CI/CD

- [ ] 遥测

- [ ] 路由支持 http 或 GRPCGateWay

- [ ] 消息

- [ ] 缓存


