package main

import (
	"context"
	"github.com/luckybet100/protodeps/samples/gitdeps/gen/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"log"
	"net"
)

type impl struct {
	proto.UnimplementedHelloServer
}

func (*impl) Hello(_ context.Context, req *proto.HelloRequest) (*proto.HelloReply, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "request validation failed")
	}
	return &proto.HelloReply{
		Message: "Hello, " + req.Name + "!",
	}, nil
}

func main() {
	listenAddr := "127.0.0.1:8000"
	listener, err := net.Listen("tcp", listenAddr)
	if err != nil {
		log.Fatalln(err)
	}
	grpcServer := grpc.NewServer()
	proto.RegisterHelloServer(grpcServer, &impl{})
	reflection.Register(grpcServer)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalln(err)
	}
}
