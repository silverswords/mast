package mastgrpc

import (
	"log"
	"net"
	"reflect"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Server -
type Server struct {
	*grpc.Server
	lis net.Listener
}

// Server -
func (b *GRPCBuilder) Server() *Server {

	if len(b.unaryServerInterceptors) != 0 {
		b.sopts = append(b.sopts, middleware.WithUnaryServerChain(b.unaryServerInterceptors...))
	}

	if len(b.streamServerInterceptors) != 0 {
		b.sopts = append(b.sopts, middleware.WithStreamServerChain(b.streamServerInterceptors...))
	}

	lis, err := net.Listen("tcp", b.Addr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return &Server{
		grpc.NewServer(b.sopts...),
		lis,
	}
}

// Prepare -
func (s *Server) Prepare(registerFunc, service interface{}) {
	f := reflect.ValueOf(registerFunc)
	if f.Type().NumIn() != 2 {
		grpclog.Fatal("The number of params is not adapted.")
	}

	if f.Type().In(0) != reflect.TypeOf(s.Server) {
		grpclog.Fatal("registerFunc aren't for GRPCServer")
	}

	in := make([]reflect.Value, 2)
	in[0] = reflect.ValueOf(s.Server)
	in[1] = reflect.ValueOf(service)
	f.Call(in)
}

// Serve -
func (s *Server) Serve() error {
	return s.Server.Serve(s.lis)
}

// Stop -
func (s *Server) Stop() {
	s.Server.Stop()
}

//// InvokeOptions struct having information about microservice API call parameters
//type InvokeOptions struct {
//	Stream bool
//	// Transport Dial Timeout
//	DialTimeout time.Duration
//	// Request/Response timeout
//	RequestTimeout time.Duration
//	// end to end, Directly call
//	Endpoint string
//	// end to end, Directly call
//	Protocol string
//	Port     string
//	//loadbalancer stratery
//	//StrategyFunc loadbalancer.Strategy
//	StrategyFunc string
//	Filters      []string
//	URLPath      string
//	MethodType   string
//	// local data
//	Metadata map[string]interface{}
//	// tags for router
//	RouteTags utiltags.Tags
//}
