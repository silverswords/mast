package mast

import "net/rpc"

// Builder could build for given parameters to make
// rpc and grpc client.
type Builder interface {
	BuildClient(rpcname string) (interface{}, error)
	BuildServer(rpcname string) (interface{}, error)
	BuildResolver(rpcname string) (interface{}, error)
}

// chose one
// server {
// rpc.Server
// }

//chose two
// rpc.Server{}

// BuilderOptions means how to build rpc
// provide some BuilderOptions.
type BuilderOptions struct {
	// rpc BuilderOptions

	// grpc-go BuilderOptions

}

func (bopts *BuilderOptions) rpcclient() *rpc.Client {
	return &rpc.Client{}
}
func (bopts *BuilderOptions) rpcServer() *rpc.Server {
	return &rpc.Server{}
}
