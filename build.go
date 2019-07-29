package mast

import "net/rpc"

type Builder interface {
	Build(target Target, cc ClientConn, opts BuildOption) (mast, error)
}

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
