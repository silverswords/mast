package mast

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"testing"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

func listenTCP() (net.Listener, string) {
	l, e := net.Listen("tcp", "127.0.0.1:1234")
	if e != nil {
		log.Fatal("net.Listen tcp:0:", e)
	}
	return l, l.Addr().String()
}

func startTCPServer() {
	rpc.Register(new(Arith))
	rpc.RegisterName("net.rpc.Arith", new(Arith))

	var l net.Listener
	l, serverAddr := listenTCP()

	log.Println("listening on ", serverAddr)
	go rpc.Accept(l)
}

func startHTTPServer() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

func HTTPclient() {
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		log.Fatal("arith error:", err)
	}
	fmt.Printf("Arith: %d*%d=%d", args.A, args.B, reply)

	// Asynchronous call
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	fmt.Println(replyCall)
}

func TestTCP(t *testing.T) {
	rpc.Register(new(Arith))
	rpc.RegisterName("net.rpc.Arith", new(Arith))

	var l net.Listener
	l, serverAddr := listenTCP()

	t.Log("listening on ", serverAddr)
	go rpc.Accept(l)

	client, err := rpc.Dial("tcp", "127.0.0.1:1234")
	if err != nil {
		t.Error("dialing", err)
	}
	defer client.Close()

	// Synchronous call
	args := &Args{7, 8}
	var reply int
	err = client.Call("Arith.Multiply", args, &reply)
	if err != nil {
		t.Error("arith error:", err)
	}
	if args.A*args.B != reply {
		t.Errorf("Arith: %d*%d=%d", args.A, args.B, reply)
	}
	// Asynchronous call
	quotient := new(Quotient)
	divCall := client.Go("Arith.Divide", args, quotient, nil)
	replyCall := <-divCall.Done // will be equal to divCall
	// check errors, print, etc.
	t.Log(replyCall.Args, replyCall.Reply)
}
