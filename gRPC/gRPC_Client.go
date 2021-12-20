package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial(":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer conn.Close()

	c := pb.NewGrpcSendMsgClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	var s string
	fmt.Print("입력하세요 : ")
	fmt.Scanln(&s)

	r, err := c.SendMsg(ctx, &pb.MessageRequest{Message: s})
	if err != nil {
		log.Fatalf("could not send : %v", err)
	}
	log.Printf("Sending : %s", r.GetMessage())
}
