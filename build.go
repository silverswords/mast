package mast

// Builder could build for given parameters to make
// rpc and grpc client and server
type Builder interface {
	Client() interface{}
	Server() interface{}
}

// Mast is Builder
// could use namespace like grpc.client.gzip
type Mast struct {
	// all option
	BuilderOptions
}

// Server could listen and serve
type Server interface {
	Listen()
	Serve()
}

// Client supposed Synchronous and Asynchronous
type Client interface {
	Call()
	Go()
}

// Client return Client
func (m *Mast) Client() Client {
	return nil
}

// Server return server
func (m *Mast) Server() Server {
	return nil
}

// Config configure options in rpc
func (m *Mast) Config(opt string) {

}
