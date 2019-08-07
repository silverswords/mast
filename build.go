package mast

import (
	"log"

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

// BuildClient realize interface to gen Client
func (m *Mast) BuildClient(rpcname string) interface{} {
	switch rpcname {
	case "RPC":
		return m.BuilderOptions.RPCClient()
	case "GRPC":
		return m.BuilderOptions.GRPCClient()
	}

	log.Fatal("Unknown Error on Creating Client")
	return nil
}

// BuildServer realize interface to gen Server
func (m *Mast) BuildServer(rpcname string) interface{} {
	switch rpcname {
	case "RPC":
		return m.BuilderOptions.RPCServer()
	case "GRPC":
		return m.BuilderOptions.GRPCServer()
	}

	log.Fatal("Unknown Error on Creating Server")
	return nil
}
