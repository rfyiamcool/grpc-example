package main

import (
	"context"
	"flag"
	"io"
	"log"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/server-side-streaming/proto"
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
	req := pb.UserRequest{ID: 1}
	stream, err := grpcClient.GetUserInfo(context.Background(), &req)
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatalf("receive error: %v", err)
		}

		log.Printf("[RECEIVED RESPONSE]: %v\n", resp)
	}
}
