package mast


// Builder could build for given parameters to make
// rpc and grpc client and server
type Builder interface {
	Dial() Client
	Server() Server
}

// ConfigBuilder it's useful for grpc to build Client from ClientConn and PB file (or named registerFunc)
type ConfigBuilder interface {
	Builder
	Config(*ConfigBuilder)
}

// Options configure Builder's options
type Options func(*Builder) error

// Server could listen and serve
type Server interface {
	Prepare(info, registerFunc interface{})
	Start()
	Stop() error
}

// Client supposed Synchronous and Asynchronous
type Client interface {
	Close() error
	ReloadConfigs(Options)
	GetOptions() Options
}
