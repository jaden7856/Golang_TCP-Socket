package main

import (
	"context"
	streamPb "github.com/jaden7856/Golang_TCP-Socket/gRPC/streamProtoc"
	"google.golang.org/grpc"
	"io"
	"log"
	"math/rand"
	"time"
)

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

	// third goroutine closes done channel
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
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()

	stmClient := streamPb.NewGRPCSendMsgClient(conn)

	sreamCallSendMsg(stmClient)
}
