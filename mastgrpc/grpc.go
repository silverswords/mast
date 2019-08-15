package mastgrpc

import (
	"log"
	"net"
	"reflect"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"golang.org/x/oauth2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/oauth"
)

// type serverOptions struct {
// 	creds                 credentials.TransportCredentials
// 	codec                 baseCodec
// 	cp                    Compressor
// 	dc                    Decompressor
// 	unaryInt              UnaryServerInterceptor
// 	streamInt             StreamServerInterceptor
// 	inTapHandle           tap.ServerInHandle
// 	statsHandler          stats.Handler
// 	maxConcurrentStreams  uint32
// 	maxReceiveMessageSize int
// 	maxSendMessageSize    int
// 	unknownStreamDesc     *StreamDesc
// 	keepaliveParams       keepalive.ServerParameters
// 	keepalivePolicy       keepalive.EnforcementPolicy
// 	initialWindowSize     int32
// 	initialConnWindowSize int32
// 	writeBufferSize       int
// 	readBufferSize        int
// 	connectionTimeout     time.Duration
// 	maxHeaderListSize     *uint32
// }

type GrpcBuilder struct {
	address string
	// grpc-go grpcBuilder

	compressorName string // like gzip.Name
	// token should is for client,should handle on server
	token string

	// only for use same serverCert
	// todo: auto apply TLS
	// should use path to cert
	serverHostOverride string
	serverCert         string
	serverKey          string

	unaryServerInterceptors  []grpc.UnaryServerInterceptor
	streamServerInterceptors []grpc.StreamServerInterceptor
	unaryClientInterceptors  grpc.UnaryClientInterceptor
	streamClientInterceptors grpc.StreamClientInterceptor
}

// DefaultGRPCBuildOptions return grpcBuilder
// which realized Builder interface
func DefaultGRPCBuildOptions() *GrpcBuilder {
	return &GrpcBuilder{
		address: "127.0.0.1:20001",
	}
}

type Server struct {
	*grpc.Server
	*GrpcBuilder
}

type Client struct {
	*grpc.ClientConn
}

func (s *Server) Prepare(service, registerFunc interface{}) {
	f := reflect.ValueOf(registerFunc)
	if 2 != f.Type().NumIn() {
		log.Fatal("The number of params is not adapted.")
	}

	p := make([]reflect.Value, 2)
	p[0] = reflect.ValueOf(s.Server)
	p[1] = reflect.ValueOf(service)
	f.Call(p)
}

func (s *Server) Serve() error {
	lis, err := net.Listen("tcp", s.GrpcBuilder.address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	return s.Server.Serve(lis)
}

func (s *Server) Stop() error{return nil}
func (c *Client) Call() {

}

func (c *Client) Go() {

}
func (bopts *GrpcBuilder) Server() *Server {
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

	return &Server{
		grpc.NewServer(opts...),
		bopts,
	}

	// panic("Create GRPCServer Fail")
}

// Client return a ClientConn by DialOption
// then you need use pb.New[ServiceName]Client(yourClientConn)
// to Create client which could Call Service and use context
// Should： ClientConn should be closed by Close()
func (bopts *GrpcBuilder) Client() *Client {
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

	return &Client{
		client,
	}
	// panic("Create GRPCClient Fail")
}
