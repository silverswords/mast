package mast

import "log"

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
	// network & address & port
	network string
	address string
	// rpc BuilderOptions
	rpcmode  uint8                  // 0 use tcp,1 use http, 2 use json
	rcvrs    map[string]interface{} //receiver of methods for service
	httppath string
	// grpc-go BuilderOptions

}

// Mast is instance of BuilderOptions to build client and server
type Mast struct {
	BuilderOptions
}

// BuildClient realize interface to gen Client
func (m *Mast) BuildClient(rpcname string) interface{} {
	switch rpcname {
	case "RPC":
		return m.BuilderOptions.RPCClient()
	case "gRPC":
	}

	log.Fatal("Unknown Error on Creating Client")
	return nil
}

// BuildServer realize interface to gen Server
func (m *Mast) BuildServer(rpcname string) interface{} {
	switch rpcname {
	case "RPC":
		return m.BuilderOptions.RPCServer()
	case "gRPC":
	}

	log.Fatal("Unknown Error on Creating Client")
	return nil
}
