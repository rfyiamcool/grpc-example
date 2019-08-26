package main

import (
	"context"
	"flag"
	"log"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/client-side-streaming/proto"
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
		log.Fatalf("grpc getUserInfo error: %v", err)
	}

	var i int32
	for i = 1; i < 4; i++ {
		err := stream.Send(&pb.UserRequest{ID: i})
		if err != nil {
			log.Fatalf("send error: %v", err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	log.Printf("[RECEIVED RESPONSE]: %v\n", resp) // 输出响应
}
