package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"

	streamPb "github.com/jaden7856/go-tcp_grpc-server-client/gRPC/streamProtoc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	stPort = flag.Int("port", 8080, "the port to serve on")
	sleep  = flag.Duration("sleep", time.Second*5, "duration between changes in health")

	serviceName = "test"
	isHealth    = false
)

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
	grpcServer := grpc.NewServer(
	//grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
	//	grpc_recovery.StreamServerInterceptor(),
	//)),
	)

	healthCheck := health.NewServer()
	healthpb.RegisterHealthServer(grpcServer, healthCheck)
	streamPb.RegisterGRPCSendMsgServer(grpcServer, &stServer{})

	// 정상적으로 연결이 된 상태
	healthCheck.SetServingStatus(serviceName, healthpb.HealthCheckResponse_SERVING)

	log.Printf("start gRPC server on 8080 port")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}

}
