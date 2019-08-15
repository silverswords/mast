package rpc

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

func defaultRPCBuildOptions() *BuilderOptions {
	return &BuilderOptions{
		address: DefaultAddress,
		rpcmode: DefaultNetwork,
		rcvrs:   make(map[string]interface{}),
	}
}

// RPCServer new a rpc server to serve conn by
// builder_options.address,rpcmode,rcvrs
func (bopts *BuilderOptions) RPCServer() *rpc.Server {
	s := rpc.NewServer()

	//register methods
	for rcvrName, rcvr := range bopts.rcvrs {
		err := s.RegisterName(rcvrName, rcvr)
		if err != nil {
			log.Println("register ", rcvrName, " ", err.Error())
		}
	}

	ln, serverAddr := func() (net.Listener, string) {
		l, e := net.Listen("tcp", bopts.address)
		if e != nil {
			log.Fatal("net.Listen on tcp: ", bopts.address)
		}
		return l, l.Addr().String()
	}()

	switch bopts.rpcmode {
	case TCP:
		log.Println("TCPRPC server listening on", serverAddr)
		go s.Accept(ln)

	case HTTP:
		if bopts.httppath != "" {
			s.HandleHTTP(bopts.httppath, "/debug"+bopts.httppath)
		} else {
			s.HandleHTTP(DefaultRPCPath, DefaultDebugPath)
		}

		ts := &httptest.Server{
			Listener: ln,
			Config:   &http.Server{Handler: nil},
		}
		ts.Start()

		log.Println("Test HTTP RPC server listening on", serverAddr)

	case JSON:
		go func() {
			for {
				conn, err := ln.Accept()
				if err != nil {
					log.Fatal("rpc.Serve: accept:", err.Error())
				}
				go s.ServeCodec(jsonrpc.NewServerCodec(conn))
			}
		}()
		log.Println("JSON server listening on", serverAddr)
	}

	return s
}

// RPCClient should call defer Client.Close() to exit graceful
// When new a client ,could use it Call() and Go() method.
// if nil ,means it's wrong to create a client.
func (bopts *BuilderOptions) RPCClient() *rpc.Client {
	switch bopts.rpcmode {
	case TCP:
		client, err := rpc.Dial("tcp", bopts.address)
		if err != nil {
			log.Fatal("Client Dial TCP error:", err.Error())
		}
		return client

	case HTTP:
		var client *rpc.Client
		var err error
		if bopts.httppath == "" {
			client, err = rpc.DialHTTP("tcp", bopts.address)
		} else {
			client, err = rpc.DialHTTPPath("tcp", bopts.address, bopts.httppath)
		}
		if err != nil {
			log.Fatal("[error] :" + err.Error())
		}
		return client

	case JSON:
		client, err := jsonrpc.Dial("tcp", bopts.address)
		if err != nil {
			log.Fatal("[error]: ", err.Error())
		}

		return client
	}

	panic("Create Client Fail")
}
