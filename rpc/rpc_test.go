package rpc

import (
	"fmt"
	"github.com/silverswords/mast"
	"log"
	"net"
	"net/http"

	"testing"
)

func TestTCPBuilder(t *testing.T) {
	bopts := defaultRPCBuildOptions()

	bopts.rcvrs["Arith"] = new(mast.Arith)
	bopts.RPCServer()

	args := &mast.Args{7, 8}
	var reply int
	err := bopts.RPCClient().Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if args.B*args.A != reply {
		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
	}
	t.Logf("Arith: %d*%d=%d", args.A, args.B, reply)
}

func TestHTTPBuilder(t *testing.T) {
	bopts := defaultRPCBuildOptions()
	bopts.rpcmode = 1

	bopts.rcvrs["Arith"] = new(mast.Arith)
	bopts.RPCServer()

	args := &mast.Args{7, 8}
	var reply int
	err := bopts.RPCClient().Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if args.B*args.A != reply {
		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
	}
	t.Logf("Arith: %d*%d=%d", args.A, args.B, reply)
}

func TestJSONBuilder(t *testing.T) {
	bopts := defaultRPCBuildOptions()
	bopts.rpcmode = 2

	bopts.rcvrs["Arith"] = new(mast.Arith)
	bopts.RPCServer()

	args := &mast.Args{7, 8}
	var reply int
	err := bopts.RPCClient().Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal(err)
	}
	if args.B*args.A != reply {
		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
	}
	t.Logf("Arith: %d*%d=%d", args.A, args.B, reply)
}

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("net.Listen tcp:0:", e)
	}
	return l, l.Addr().String()
}

func startTCPServer() {
	Register(new(mast.Arith))
	RegisterName("net.rpc.Arith", new(mast.Arith))

	var l net.Listener
	l, serverAddr := listenTCP()

	log.Println("listening on ", serverAddr)
	go Accept(l)
}

func startHTTPServer() {
	arith := new(mast.Arith)
	Register(arith)
	HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func HTTPclient() {
	client, err := DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &mast.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	// Asynchronous call
	quotient := new(mast.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	fmt.Println(replyCall)
}

func TestTCP(t *testing.T) {
	Register(new(mast.Arith))
	RegisterName("net.rpc.Arith", new(mast.Arith))

	var l net.Listener
	l, serverAddr := listenTCP()

	t.Log("listening on ", serverAddr)
	go Accept(l)

	client, err := Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Error("dialing", err)
	}
	defer client.Close()

	// Synchronous call
	args := &mast.Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		t.Error("arith error:", err)
	}
	if args.A*args.B != reply {
		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
	}
	// Asynchronous call
	quotient := new(mast.Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	t.Log(replyCall.Args, replyCall.Reply)
}
