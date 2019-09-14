package mastgrpc

import (
	"time"

	"google.golang.org/grpc"
)


// defaultTarget temporary
var defaultBuilder = &GRPCBuilder{
	Network: "tcp",
	Addr:"127.0.0.1:21001",

	ClientDialDeadline: time.Second,
	ClientCompresser: "gzip",


}

// GRPCBuilder - instance of mast
type GRPCBuilder struct {
	// Network is grpc listen network,default value is tcp
	Network string `dsn:"network"`
	// Addr is grpc listen addr,default value is 0.0.0.0:9000
	Addr string `dsn:"address"`

	ClientDialDeadline time.Duration `dsn:"client-side"`
	ClientCompresser   string

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
