package mastgrpc

import (
	"context"

	middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// ClientConn get grpc.ClientConn
func (b *GRPCBuilder) ClientConn() (*grpc.ClientConn, error) {
	return b.Dial()
}

// Dial return a ClientConn by DialOption
// then you need use pb.New[ServiceName]Client(yourClientConn)
// to Create client which could Call Service and use context
// Should： ClientConn should be closed by Close()
func (b *GRPCBuilder) Dial() (*grpc.ClientConn, error) {
	return b.dialContext(context.Background())
}

func (b *GRPCBuilder) dialContext(context context.Context) (*grpc.ClientConn, error) {
	b.dopts = append(b.dopts, grpc.WithInsecure())
	if len(b.unaryServerInterceptors) != 0 {
		b.dopts = append(b.dopts, grpc.WithUnaryInterceptor(middleware.ChainUnaryClient(b.unaryClientInterceptors...)))
	}

	if len(b.streamServerInterceptors) != 0 {
		b.dopts = append(b.dopts, grpc.WithStreamInterceptor(middleware.ChainStreamClient(b.streamClientInterceptors...)))
	}

	return grpc.DialContext(context, b.Addr, b.dopts...)
}

// DialTLS creates a client connection over tls transport with given serverCert and server's name.
func (b *GRPCBuilder) DialTLS(ctx context.Context, file string, name string) (conn *grpc.ClientConn, err error) {
	var creds credentials.TransportCredentials
	creds, err = credentials.NewClientTLSFromFile(file, name)
	if err != nil {
		return
	}
	b.dopts = append(b.dopts, grpc.WithTransportCredentials(creds))
	return b.dialContext(ctx)
}

// Client to call Service
type Client struct {
	ClientConn *grpc.ClientConn
	context.Context
}
