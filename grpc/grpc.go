package grpc

import (
	"github.com/silverswords/mast"
	"log"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

func defaultGRPCBuildOptions() *mast.BuilderOptions {
	return &mast.BuilderOptions{
		address: mast.DefaultAddress,
	}
}

// chose 🥁
// user register themself and call server's serve by s.Serve

// GRPCServer return a Server make by grpc.ServerOption，
// then you need use pb.Register[ServiceName]Server(yourServerName,yourRealizeServer)
func (bopts *mast.BuilderOptions) GRPCServer() *grpc.Server {
	var opts []grpc.ServerOption

	if bopts.serverCert != "" && bopts.serverKey != "" {
		creds, err := credentials.NewServerTLSFromFile(bopts.serverCert, bopts.serverKey)
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		opts = append(opts, grpc.Creds(creds))
	}

	if len(bopts.unaryServerInterceptors) != 0 {
		opts = append(opts, middleware.WithUnaryServerChain(bopts.unaryServerInterceptors...))
	}

	if len(bopts.streamServerInterceptors) != 0 {
		opts = append(opts, middleware.WithStreamServerChain(bopts.streamServerInterceptors...))
	}

	return grpc.NewServer(opts...)

	// panic("Create GRPCServer Fail")
}

// GRPCClient return a ClientConn by DialOption
// then you need use pb.New[ServiceName]Client(yourClientConn)
// to Create client which could Call Service and use context
// Should： ClientConn should be closed by Close()
func (bopts *mast.BuilderOptions) GRPCClient() *grpc.ClientConn {
	var opts []grpc.DialOption

	switch {
	case bopts.serverHostOverride != "" && bopts.serverCert != "":
		creds, err := credentials.NewClientTLSFromFile(bopts.serverCert, bopts.serverHostOverride)
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

	if len(bopts.unaryServerInterceptors) != 0 {
		opts = append(opts, grpc.WithUnaryInterceptor(bopts.unaryClientInterceptors))
	}

	if len(bopts.streamServerInterceptors) != 0 {
		opts = append(opts, grpc.WithStreamInterceptor(bopts.streamClientInterceptors))
	}

	client, err := grpc.Dial(bopts.address, opts...)
	if err != nil {
		log.Fatal("cannot connect ", bopts.address)
	}

	return client
	// panic("Create GRPCClient Fail")
}
