package mast

// Builder could build for given parameters to make
// rpc and grpc client and server
type Builder interface {
	Client() interface{}
	Server() interface{}
}
type Mast struct{
	// all option
	BuildOptions
	grpc.ServerOption
	map[string]func (*grpc.ServerOption)()
}

type Server interface{
	Listen()
	Serve()
}
type Client interface{
	Call()
	Go()
}
func (m *Mast) Client() Client{
	return nil
}

func (m *Mast) Server() Server{
	return nil
}

func (m *Mast) Config(opt string){

}