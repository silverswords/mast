package mastgrpc

import (
	"google.golang.org/grpc"
)

const (
	// DefaultTarget temporary
	DefaultTarget = "127.0.0.1:21001"
)

//Interceptors
//Please send a PR to add new interceptors or middleware to this list
//
//Auth
//grpc_auth - a customizable (via AuthFunc) piece of auth middleware
//Logging
//grpc_ctxtags - a library that adds a Tag map to context, with data populated from request body
//grpc_zap - integration of zap logging library into gRPC handlers.
//grpc_logrus - integration of logrus logging library into gRPC handlers.
//Monitoring
//grpc_prometheus⚡ - Prometheus client-side and server-side monitoring middleware
//otgrpc⚡ - OpenTracing client-side and server-side interceptors
//grpc_opentracing - OpenTracing client-side and server-side interceptors with support for streaming and handler-returned tags
//Client
//grpc_retry - a generic gRPC response code retry mechanism, client-side middleware
//Server
//grpc_validator - codegen inbound message validation from .proto options
//grpc_recovery - turn panics into gRPC errors
//ratelimit - grpc rate limiting by your own limiter
//Status
//This code has been running in production since May 2016 as the basis of the gRPC micro services stack at Improbable.
//
//Additional tooling will be added, and contributions are welcome.
//

// GRPCBuilder - instance of mast
type GRPCBuilder struct {
	// Network is grpc listen network,default value is tcp
	Network string `dsn:"network"`
	// Addr is grpc listen addr,default value is 0.0.0.0:9000
	Addr string `dsn:"address"`


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
		Addr: DefaultTarget,
	}
}

// BuildOption configures how we set up the connection.
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

// WithAddr -
func WithAddr(addr string) BuildOption {
	return newFuncDialOption(func(b *GRPCBuilder) {
		b.Addr = addr
	})
}
