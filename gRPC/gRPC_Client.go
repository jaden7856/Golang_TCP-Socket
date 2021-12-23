package main

import (
	"context"
	"log"
	"runtime"
	"time"

	pb "github.com/jaden7856/Golang_TCP-Socket/gRPC/protoc"
	"google.golang.org/grpc"
)

var (
	msg1 = "client-1 입니다."
	msg2 = "client-2 입니다."
	msg3 = "client-3 입니다."
)

func callSendMsg(c pb.GRPCSendMsgClient, msg string) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	r, err := c.SendMsg(ctx, &pb.MessageRequest{Message: msg})
	if err != nil {
		log.Fatalf("could not send : %v", err)
	}
	log.Printf("Sending : %s", r.GetMessage())
}

func main() {
	conn, err := grpc.Dial(":8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()

	msgc := pb.NewGRPCSendMsgClient(conn)

	go callSendMsg(msgc, msg1)
	runtime.Gosched()
	go callSendMsg(msgc, msg2)
	runtime.Gosched()
	callSendMsg(msgc, msg3)

}
