package mast

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

// BuilderOptions means how to build rpc
// provide some BuilderOptions.
type BuilderOptions struct {
	// address & port
	address string

	// gorpc BuilderOptions
	rpcmode  uint8                  // 0 use tcp,1 use http, 2 use json
	rcvrs    map[string]interface{} //receiver of methods for service
	httppath string
}
