package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/bidirectional-streaming/proto"
)

var (
	addr = flag.String("addr", "0.0.0.0:3001", "connect server address")
)

func main() {
	flag.Parse()
	conn, err := grpc.Dial(*addr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()

	grpcClient := pb.NewUserServiceClient(conn)
	stream, err := grpcClient.GetUserInfo(context.Background())
	if err != nil {
		log.Fatalf("receive stream error: %v", err)
	}

	var i int32
	for i = 1; i < 100; i++ {
		stream.Send(&pb.UserRequest{ID: i})
		resp, err := stream.Recv()
		if err != nil {
			log.Fatalf("resp error: %v", err)
		}

		log.Printf("[RECEIVED RESPONSE]: %v\n", resp)
		time.Sleep(1 * time.Second)
	}
}
