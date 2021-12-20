package main

import (
	"golang.org/x/net/context"
	"log"
	"net"

	pb "github.com/jaden7856/Golang_TCP-Socket/gRPC/message"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedGrpcSendMsgServer
}

func (s *Server) SendMsg(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	log.Printf("Received : %v", in.GetMessage())
	return &pb.MessageReply{Message: in.GetMessage()}, nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	pb.RegisterGrpcSendMsgServer(grpcServer, &Server{})

	log.Printf("start gRPC server on 8080 port")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
