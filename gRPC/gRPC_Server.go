package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net"

	pb "github.com/jaden7856/Golang_TCP-Socket/gRPC/protoc"
	"google.golang.org/grpc"
)

var port = flag.Int("port", 8080, "the port to serve on")

type Server struct {
	pb.UnimplementedGRPCSendMsgServer
}

func (s *Server) SendMsg(ctx context.Context, in *pb.MessageRequest) (*pb.MessageReply, error) {
	log.Printf("Received : %v", in.GetMessage())
	return &pb.MessageReply{Message: in.GetMessage()}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()
	pb.RegisterGRPCSendMsgServer(grpcServer, &Server{})

	log.Printf("start gRPC server on 8080 port")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
