package mast

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
)

const (
	// DefaultRPCPath used by HandleHTTP
	DefaultRPCPath = "/_goRPC_"
	// DefaultDebugPath DebugPath
	DefaultDebugPath = "/debug/rpc"
)

// Builder could build for given parameters to make
// rpc and grpc client.
type Builder interface {
	BuildClient(rpcname string) (interface{}, error)
	BuildServer(rpcname string) (interface{}, error)
	BuildResolver(rpcname string) (interface{}, error)
}

// chose one ✅
// server {
// BuilderOptions
// rpc.Server
// }
//  when call server.server.Accept(server.opts)

//chose two ❌
// rpc.Server{}
// when call server.Accept(bopts)

// BuilderOptions means how to build rpc
// provide some BuilderOptions.
type BuilderOptions struct {
	// network & address & port
	network string
	address string
	// rpc BuilderOptions
	rpcmode uint8                  // 0 use tcp,1 use http, 2 use json
	rcvrs   map[string]interface{} //receiver of methods for service

	// grpc-go BuilderOptions

}

func defaultBuildOptions() *BuilderOptions {
	return &BuilderOptions{
		network: "tcp",
		address: "127.0.0.1:0",
		rpcmode: 0,
	}
}

// RPCServer new a rpc server to serve conn by
// builder_options.network,address,rpcmode,rcvrs
func (bopts *BuilderOptions) RPCServer() *rpc.Server {
	s := rpc.NewServer()

	//register methods
	for rcvrName, rcvr := range bopts.rcvrs {
		if rcvrName == "" {
			s.Register(rcvr)
			continue
		}
		s.RegisterName(rcvrName, rcvr)
	}

	ln, serverAddr := func() (net.Listener, string) {
		l, e := net.Listen(bopts.network, bopts.address)
		if e != nil {
			log.Fatal("net.Listen ", bopts.network, bopts.address)
		}
		return l, l.Addr().String()
	}()

	switch bopts.rpcmode {
	case 0:
		log.Println("RPC server listening on", serverAddr)
		go s.Accept(ln)

	case 1:
		if bopts.network != "tcp" && bopts.network != "tcp6" {
			log.Fatal("cannot start http server by", bopts.network)
		}
		s.HandleHTTP(DefaultRPCPath, DefaultDebugPath)

		ts := &httptest.Server{
			Listener: ln,
			Config:   &http.Server{Handler: nil},
		}
		ts.Start()

		log.Println("Test HTTP RPC server listening on", serverAddr)

	case 2:
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Fatal("rpc.Serve: accept:", err.Error())
				}
				go jsonrpc.ServeConn(conn)
			}
		}()

		log.Println("JSON server listening on", serverAddr)
	}

	return s
}

func (bopts *BuilderOptions) rpcclient() *rpc.Client {
	return &rpc.Client{}
}
