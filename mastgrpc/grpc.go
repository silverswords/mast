package mastgrpc

import (
	"google.golang.org/grpc"
)

const (
	DefaultTarget = "127.0.0.1:21001"
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

type GRPCBuilder struct {
	// Network is grpc listen network,default value is tcp
	Network string `dsn:"network"`
	// Addr is grpc listen addr,default value is 0.0.0.0:9000
	Addr string `dsn:"address"`

	// temporary use bool now
	secureConfig bool

	sopts []grpc.ServerOption
	dopts []grpc.DialOption

	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor

	unaryClientInterceptors  []grpc.UnaryClientInterceptor
	streamClientInterceptors []grpc.StreamClientInterceptor
}

//// ServerConfig is rpc server conf.
//type ServerConfig struct {
//	// Timeout is context timeout for per rpc call.
//	Timeout xtime.Duration `dsn:"query.timeout"`
//	// IdleTimeout is a duration for the amount of time after which an idle connection would be closed by sending a GoAway.
//	// Idleness duration is defined since the most recent time the number of outstanding RPCs became zero or the connection establishment.
//	IdleTimeout xtime.Duration `dsn:"query.idleTimeout"`
//	// MaxLifeTime is a duration for the maximum amount of time a connection may exist before it will be closed by sending a GoAway.
//	// A random jitter of +/-10% will be added to MaxConnectionAge to spread out connection storms.
//	MaxLifeTime xtime.Duration `dsn:"query.maxLife"`
//	// ForceCloseWait is an additive period after MaxLifeTime after which the connection will be forcibly closed.
//	ForceCloseWait xtime.Duration `dsn:"query.closeWait"`
//	// KeepAliveInterval is after a duration of this time if the server doesn't see any activity it pings the client to see if the transport is still alive.
//	KeepAliveInterval xtime.Duration `dsn:"query.keepaliveInterval"`
//	// KeepAliveTimeout  is After having pinged for keepalive check, the server waits for a duration of Timeout and if no activity is seen even after that
//	// the connection is closed.
//	KeepAliveTimeout xtime.Duration `dsn:"query.keepaliveTimeout"`
//	// LogFlag to control log behaviour. e.g. LogFlag: warden.LogFlagDisableLog.
//	// Disable: 1 DisableArgs: 2 DisableInfo: 4
//	LogFlag int8 `dsn:"query.logFlag"`
//}
// DefaultGRPCBuildOptions return GRPCBuilder
// which realized Builder interface
func DefaultGRPCBuildOptions() *GRPCBuilder {
	return &GRPCBuilder{
		Addr: DefaultTarget,
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

func WithAddr(addr string) BuildOption {
	return newFuncDialOption(func(b *GRPCBuilder) {
		b.Addr = addr
	})
}
