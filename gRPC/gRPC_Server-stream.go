package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	streamPb "github.com/jaden7856/Golang_TCP-Socket/gRPC/streamProtoc"
	"google.golang.org/grpc"
)

var stPort = flag.Int("port", 8080, "the port to serve on")

type stServer struct {
	streamPb.UnimplementedGRPCSendMsgServer
}

func (*stServer) SendMsg(msgServer streamPb.GRPCSendMsg_SendMsgServer) error {
	log.Println("start new server")

	var max int32
	ctx := msgServer.Context()

	for {

		//exit if context is done
		//or continue
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
		}

		// receive data from stream
		req, err := msgServer.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
		}

		// continue if number reveived from stream
		// less than max
		if req.Num <= max {
			continue
		}

		// update max and send it to stream
		max = req.Num
		resp := streamPb.MessageReply{Result: max}
		if err := msgServer.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Printf("send new max=%d", max)
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *stPort))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// gRPC 서버 생성
	grpcServer := grpc.NewServer()
	streamPb.RegisterGRPCSendMsgServer(grpcServer, &stServer{})

	log.Printf("start gRPC server on 8080 port")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
