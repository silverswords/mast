package mast

import (
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc"
	"net/rpc/jsonrpc"
)

func defaultRPCBuildOptions() *BuilderOptions {
	return &BuilderOptions{
		network: "tcp",
		address: "127.0.0.1:21001",
		rpcmode: 0,
		rcvrs:   make(map[string]interface{}),
	}
}

// RPCServer new a rpc server to serve conn by
// builder_options.network,address,rpcmode,rcvrs
func (bopts *BuilderOptions) RPCServer() *rpc.Server {
	s := rpc.NewServer()

	//register methods
	for rcvrName, rcvr := range bopts.rcvrs {
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
	case TCP:
		log.Println("TCPRPC server listening on", serverAddr)
		go s.Accept(ln)

	case HTTP:
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

	case JSON:
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

// RPCClient should call defer Client.Close() to exit graceful
// When new a client ,could use it Call() and Go() method.
// if nil ,means it's wrong to create a client.
func (bopts *BuilderOptions) RPCClient() *rpc.Client {
	switch bopts.rpcmode {
	case TCP:
		client, err := rpc.Dial(bopts.network, bopts.address)
		if err != nil {
			log.Fatal("Client Dial TCP error:", err.Error())
		}
		return client
	case HTTP:
		var client *rpc.Client
		var err error
		if bopts.httppath == "" {
			client, err = rpc.DialHTTP("tcp", bopts.network+bopts.address)
		} else {
			client, err = rpc.DialHTTPPath("tcp", bopts.network+bopts.address, bopts.httppath)
		}
		if err != nil {
			log.Fatal("dialing http use ", bopts.network, err.Error())
		}
		return client
	case JSON:
		client, err := jsonrpc.Dial(bopts.network, bopts.address)
		if err != nil {
			log.Fatal("dialing JSON use ", bopts.network, " error:", err.Error())
		}
		return client
	}

	panic("Create Client Fail")
}
