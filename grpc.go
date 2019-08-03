package mast

import (
	"context"
	"log"
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func defaultGRPCBuildOptions() *BuilderOptions {
	return &BuilderOptions{
		address: DefaultAddress,
		rpcmode: DefaultNetwork,

		timeout: 10 * time.Second,
	}
}

// To-do :
// make pb.RegisterXXXServer() and pb.NewXXXClient() in GRPCServer and GRPCClient
// judge the first parameter whether type is *grpc.Server and cannot judge
// the second parameter.
// It's not same in client by using NewXXXClient(ClientConn) to call service which register
// on server.but it supposed to return a XXXClient, and it's impossible.
// consider if use reflect to do register

// Should recognize stream and unary !!
// so resolve about context should provide two methods or even four
// two on Client point , stream should invoke sendMessage and recvMessage.

// so two things
// one: register integrate in build
// two: make stream and unary call wraped and close graceful
// and then could use
// calloptions
// context to cancel
// server.Serve() automatic

// GRPCServer return a Server make by grpc.ServerOption，
// then you need use pb.Register[ServiceName]Server(yourServerName,yourRealizeServer)
func (bopts *BuilderOptions) GRPCServer() *grpc.Server {
	var opts []grpc.ServerOption

	if bopts.creds != "" && bopts.privitekey != "" && bopts.cert != "" {
		creds, err := credentials.NewServerTLSFromFile(bopts.privitekey, bopts.cert)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	return grpc.NewServer(opts...)

	// panic("Create GRPCServer Fail")
}

// GRPCClient return a ClientConn by DialOption
// then you need use pb.New[ServiceName]Client(yourClientConn)
// to Create client which could Call Service and use context
// Should： ClientConn should be closed by Close()
func (bopts *BuilderOptions) GRPCClient() *grpc.ClientConn {
	var opts []grpc.DialOption

	switch {
	case bopts.creds != "" && bopts.privitekey != "" && bopts.cert != "":
		creds, err := credentials.NewClientTLSFromFile(bopts.creds, "")
		if err != nil {
			log.Fatalf("failed to load credentials: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
		fallthrough
	case bopts.token != "":
		perRPC := oauth.NewOauthAccess(&oauth2.Token{AccessToken: bopts.token})
		opts = append(opts, grpc.WithPerRPCCredentials(perRPC))
	default:
		opts = append(opts, grpc.WithInsecure())
	}

	if bopts.compressorName != "" {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(bopts.compressorName)))
	}

	if bopts.timeout != 0 {
		_, _ = context.WithTimeout(context.Background(), bopts.timeout)
	}
	// if bopts.creds != "" && bopts.privitekey != "" && bopts.cert != "" {
	// 	creds, err := credentials.NewClientTLSFromFile(bopts.creds, "")
	// 	if err != nil {
	// 		log.Fatalf("failed to load credentials: %v", err)
	// 	}
	// 	opts = append(opts, grpc.WithTransportCredentials(creds))
	// } else if {
	// 	perRPC := oauth.NewOauthAccess(&oauth2.Token{AccessToken: bopts.token})
	// 	opts = append(opts, grpc.WithPerRPCCredentials(perRPC))
	// }else {
	// 	opts = append(opts, grpc.WithInsecure())
	// }

	client, err := grpc.Dial(bopts.address, opts...)
	if err != nil {
		log.Fatal("cannot connect ", bopts.address)
	}

	return client
	// panic("Create GRPCClient Fail")
}
