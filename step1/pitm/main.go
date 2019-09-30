/*
 *
 * Copyright 2015 gRPC authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

//go:generate protoc -I ../helloworld --go_out=plugins=grpc:../helloworld ../helloworld/helloworld.proto

// Package main implements a server for Greeter service.
package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/kcollasarundell/grpc-demo-me/step1/helloworld"
	"google.golang.org/grpc"
)

const (
	listen   = ":50052"
	upstream = "localhost:50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	conn, err := grpc.Dial(upstream, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	in.Name = "bob"
	defer conn.Close()
	c := pb.NewGreeterClient(conn)
	out, err := c.SayHello(ctx, in)
	out.Message = fmt.Sprintf("%v, THis is a changed message", out.Message)
	log.Printf("received request for %s, and authed = %v", in.Name, out.GetAuthed())
	return out, err
}

func main() {

	lis, err := net.Listen("tcp", listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	conn, err := grpc.Dial(upstream, grpc.WithInsecure())
	if err != nil {
		conn.Close()
	}

	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

}
