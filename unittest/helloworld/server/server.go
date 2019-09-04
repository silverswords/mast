// Copyright (C) 2019 Abser Ari
//
// This file is part of mast.
//
// mast is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// mast is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with mast.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"context"
	"log"

	pb "github.com/silverswords/mast/example/helloworld"
	"github.com/silverswords/mast/mastgrpc"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received: %v", in.Name)
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	b := mastgrpc.DefaultGRPCBuildOptions()
	s := b.Server()

	s.Prepare(pb.RegisterGreeterServer, &server{})

	if err := s.Serve(); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
