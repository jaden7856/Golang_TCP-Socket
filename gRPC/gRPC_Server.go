package main

import (
	pb "Socket_gRPC/gRPC"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

type grpcServer struct {
	pb.GrpcSendMsgServer
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterRouteGuideServer(grpcServer, newServer())
	grpcServer.Serve(lis)
}
