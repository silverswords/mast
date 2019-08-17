package mastgrpc

import (
	"context"
	"log"
	"net"
	"reflect"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

// type serverOptions struct {
// 	creds                 credentials.TransportCredentials
// 	codec                 baseCodec
// 	cp                    Compressor
// 	dc                    Decompressor
// 	unaryInt              UnaryServerInterceptor
// 	streamInt             StreamServerInterceptor
// 	inTapHandle           tap.ServerInHandle
// 	statsHandler          stats.Handler
// 	maxConcurrentStreams  uint32
// 	maxReceiveMessageSize int
// 	maxSendMessageSize    int
// 	unknownStreamDesc     *StreamDesc
// 	keepaliveParams       keepalive.ServerParameters
// 	keepalivePolicy       keepalive.EnforcementPolicy
// 	initialWindowSize     int32
// 	initialConnWindowSize int32
// 	writeBufferSize       int
// 	readBufferSize        int
// 	connectionTimeout     time.Duration
// 	maxHeaderListSize     *uint32
// }

type GRPCBuilder struct {
	target string

	// temporary use bool now
	secureConfig bool

	sopts []grpc.ServerOption
	dopts []grpc.DialOption

	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor

	unaryClientInterceptors  []grpc.UnaryClientInterceptor
	streamClientInterceptors []grpc.StreamClientInterceptor
}

// DefaultGRPCBuildOptions return GRPCBuilder
// which realized Builder interface
func DefaultGRPCBuildOptions() *GRPCBuilder {
	return &GRPCBuilder{
		target: DefaultTarget,
	}
}

// DialOption configures how we set up the connection.
type BuildOption interface {
	apply(*GRPCBuilder)
}

// EmptyDialOption does not alter the dial configuration. It can be embedded in
// another structure to build custom dial options.
//
// This API is EXPERIMENTAL.
type EmptyDialOption struct{}

func (EmptyDialOption) apply(*GRPCBuilder) {}

// funcDialOption wraps a function that modifies dialOptions into an
// implementation of the DialOption interface.
type funcBuildOption struct {
	f func(*GRPCBuilder)
}

func (fdo *funcBuildOption) apply(do *GRPCBuilder) {
	fdo.f(do)
}

func newFuncDialOption(f func(*GRPCBuilder)) *funcBuildOption {
	return &funcBuildOption{
		f: f,
	}
}

func WithTarget(target string) BuildOption{
	return newFuncDialOption (func(b *GRPCBuilder){
		b.target = target
	})
}

type Server struct {
	*grpc.Server
	lis net.Listener
}

func (s *Server) Prepare(registerFunc, service interface{}) {
	f := reflect.ValueOf(registerFunc)
	if f.Type().NumIn() != 2 {
		grpclog.Fatal("The number of params is not adapted.")
	}

	if f.Type().In(0) != reflect.TypeOf(s.Server) {
		grpclog.Fatal("registerFunc aren't for GRPCServer")
	}

	p := make([]reflect.Value, 2)
	p[0] = reflect.ValueOf(s.Server)
	p[1] = reflect.ValueOf(service)
	f.Call(p)
}

func (s *Server) Serve() error {
	return s.Server.Serve(s.lis)
}

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

func (b *GRPCBuilder) Server() *Server {

	if len(b.unaryServerInterceptors) != 0 {
		b.sopts = append(b.sopts, middleware.WithUnaryServerChain(b.unaryServerInterceptors...))
	}

	if len(b.streamServerInterceptors) != 0 {
		b.sopts = append(b.sopts, middleware.WithStreamServerChain(b.streamServerInterceptors...))
	}

	lis, err := net.Listen("tcp", b.target)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	return &Server{
		grpc.NewServer(b.sopts...),
		lis,
	}
}

func (b *GRPCBuilder) Dial() (*grpc.ClientConn, error) {
	return b.dialContext(context.Background())
}

// Dial return a ClientConn by DialOption
// then you need use pb.New[ServiceName]Client(yourClientConn)
// to Create client which could Call Service and use context
// Should： ClientConn should be closed by Close()
func (b *GRPCBuilder) dialContext(context context.Context) (*grpc.ClientConn, error) {
	b.dopts = append(b.dopts, grpc.WithInsecure())
	if len(b.unaryServerInterceptors) != 0 {
		b.dopts = append(b.dopts, grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(b.unaryClientInterceptors...)))
	}

	if len(b.streamServerInterceptors) != 0 {
		b.dopts = append(b.dopts, grpc.WithStreamInterceptor(middleware.ChainStreamClient(b.streamClientInterceptors...)))
	}

	return grpc.DialContext(context, b.target, b.dopts...)
}

// DialTLS creates a client connection over tls transport with given serverCert and server's name.
func (b *GRPCBuilder) DialTLS(ctx context.Context, file string, name string) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		return
	}
	b.dopts = append(b.dopts, grpc.WithTransportCredentials(creds))
	return b.dialContext(ctx)
}

// it's not good enough for use this Client interface()
// when somebody just need simple , it will works
type Client struct {
	c interface{} // it's PB file's involve Client to call method
	context.Context
}

//func (c *Client) Call(methodName string) {
//	reflect.TypeOf(c.Client).MethodByName(methodName)
//}
