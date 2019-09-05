// Copyright (C) 2019 Abser Ari
//
// This file is part of mast.

package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/silverswords/mast/mastgrpc"

	pb "github.com/silverswords/mast/unittest/helloworld"
)

const (
	defaultName = "world"
)

func main() {
	// Set up a connection to the server.
	b := mastgrpc.DefaultGRPCBuildOptions()

	Client := b.Client(pb.NewGreeterClient)
	Client.(pb.GreeterClient).SayHello(

		)
	conn, err := b.Dial()
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Greeting: %s", r.Message)
}
