package rpc

import (
	"github.com/silverswords/mast"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/rpc/jsonrpc"
)

func defaultRPCBuildOptions() *mast.BuilderOptions {
	return &mast.BuilderOptions{
		address: mast.DefaultAddress,
		rpcmode: mast.DefaultNetwork,
		rcvrs:   make(map[string]interface{}),
	}
}

// RPCServer new a rpc server to serve conn by
// builder_options.address,rpcmode,rcvrs
func (bopts *mast.BuilderOptions) RPCServer() *Server {
	s := NewServer()

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
	case mast.TCP:
		log.Println("TCPRPC server listening on", serverAddr)
		go s.Accept(ln)

	case mast.HTTP:
		if bopts.httppath != "" {
			s.HandleHTTP(bopts.httppath, "/debug"+bopts.httppath)
		} else {
			s.HandleHTTP(mast.DefaultRPCPath, mast.DefaultDebugPath)
		}

		ts := &httptest.Server{
			Listener: ln,
			Config:   &http.Server{Handler: nil},
		}
		ts.Start()

		log.Println("Test HTTP RPC server listening on", serverAddr)

	case mast.JSON:
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
func (bopts *mast.BuilderOptions) RPCClient() *Client {
	switch bopts.rpcmode {
	case mast.TCP:
		client, err := Dial("tcp", bopts.address)
		if err != nil {
			log.Fatal("Client Dial TCP error:", err.Error())
		}
		return client

	case mast.HTTP:
		var client *Client
		var err error
		if bopts.httppath == "" {
			client, err = DialHTTP("tcp", bopts.address)
		} else {
			client, err = DialHTTPPath("tcp", bopts.address, bopts.httppath)
		}
		if err != nil {
			log.Fatal("[error] :" + err.Error())
		}
		return client

	case mast.JSON:
		client, err := jsonrpc.Dial("tcp", bopts.address)
		if err != nil {
			log.Fatal("[error]: ", err.Error())
		}

		return client
	}

	panic("Create Client Fail")
}
