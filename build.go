package mast

import (
	"github.com/silverswords/mast/rpc"
	"google.golang.org/grpc"
)

const (
	// DefaultRPCPath used by HandleHTTP
	DefaultRPCPath = "/_goRPC_"
	// DefaultDebugPath DebugPath
	DefaultDebugPath = "/debug/rpc"

	// TCP is defaultBuildOptions for communicate
	TCP = 0
	// HTTP should use path like DefaultDebugPath DefaultRPCPath
	HTTP = 1
	// JSON is change default codec for server and client
	JSON = 2

	// DefaultNetwork 0 means TCP Client and Server
	DefaultNetwork = TCP
	// DefaultAddress means itself
	DefaultAddress = "127.0.0.1:21001"
)

// Builder could build for given parameters to make
// rpc and grpc client.
type Builder interface {
	BuildClient(rpcname string) interface{}
	BuildServer(rpcname string) interface{}
}

// BuilderOptions means how to build rpc
// provide some BuilderOptions.
type BuilderOptions struct {
	// address & port
	address string

	// gorpc BuilderOptions
	rpcmode  uint8                  // 0 use tcp,1 use http, 2 use json
	rcvrs    map[string]interface{} //receiver of methods for service
	httppath string

	// grpc-go BuilderOptions

	// grpc.CallOption could use by DialOption
	// grpc.DialOption
	// timeout        time.Duration // timeout * time.Second

	// compress string // grpc.UseCompressor(gzip.Name)

	compressorName string // like gzip.Name
	// token should is for client,should handle on server
	token string

	//only for use same serverCert
	// todo: auto apply TLS
	serverHostOverride string
	serverCert         string
	serverKey          string

	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor
	unaryClientInterceptors  grpc.UnaryClientInterceptor
	streamClientInterceptors grpc.StreamClientInterceptor
}

// Mast is instance of BuilderOptions to build client and server
type Mast struct {
	*BuilderOptions
}

// BuildRPCClient return *rpc.Client
func (m *Mast) BuildRPCClient() *rpc.Client {
	return m.BuilderOptions.RPCClient()
}

// BuildGRPCClient return *grpc.ClientConn and
// is supposed to Call register from PB file
func (m *Mast) BuildGRPCClient() *grpc.ClientConn {
	return m.BuilderOptions.GRPCClient()
}

// BuildRPCServer return *rpc.Server
func (m *Mast) BuildRPCServer() *rpc.Server {
	return m.BuilderOptions.RPCServer()
}

// BuildGRPCServer return *grpc.Server and
// is supposed to registe service to serve
func (m *Mast) BuildGRPCServer() *grpc.Server {
	return m.BuilderOptions.GRPCServer()
}
