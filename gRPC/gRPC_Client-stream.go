package main

import (
	"context"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"io"
	"log"
	"math/rand"
	"time"

	streamPb "github.com/jaden7856/Golang_TCP-Socket/gRPC/streamProtoc"
	"google.golang.org/grpc"
)

var serviceConfig = `{
		"loadBalancingPolicy": "round_robin",
		"healthCheckConfig": {
			"serviceName": "PingPong"
		}
	}`

func sreamCallSendMsg(client streamPb.GRPCSendMsgClient) {
	stream, err := client.SendMsg(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
		return
	}

	var max int32
	ctx := stream.Context()
	done := make(chan bool)

	// random int Stream 전송
	go func() {
		for i := 1; i <= 50; i++ {
			// generates random number
			rnd := int32(rand.Intn(i))
			req := streamPb.MessageRequest{Num: rnd}
			if err := stream.Send(&req); err != nil {
				log.Fatalf("can not send %v", err)
			}
			log.Printf("%d sent", req.Num)
			time.Sleep(time.Millisecond * 100)
		}
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// receives data from stream
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			max = resp.Result
			log.Printf("new max %d received", max)
		}
	}()

	// goroutine closes done channel
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
	log.Printf("finished with max=%d", max)
}

func main() {
	options := []grpc.DialOption{
		// ClientConn에 대한 전송 보안을 비활성화
		grpc.WithInsecure(),
		// 연결이 작동될 때까지 Dial 호출자가 차단. 이것이 없으면 Dial은 즉시 반환되고 서버 연결은 백그라운드에서 발생합니다.
		grpc.WithBlock(),
		grpc.WithDefaultServiceConfig(serviceConfig),
	}

	conn, err := grpc.Dial(":8080", options...)

	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	stmClient := streamPb.NewGRPCSendMsgClient(conn)
	resp, err := healthpb.NewHealthClient(conn).Check(ctx, &healthpb.HealthCheckRequest{
		Service: "test",
	})

	if err != nil {
		log.Printf("can't connect grpc server: %v, code: %v\n", err, grpc.Code(err))
	} else {
		log.Printf("status: %s", resp.GetStatus().String())
	}

	sreamCallSendMsg(stmClient)
}
