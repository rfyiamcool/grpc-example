package main

import (
	"context"
	"flag"
	"io"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

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

	// create header
	md := metadata.Pairs("my-req-key1", "haha", "my-req-key2", "hello")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// req
	grpcClient := pb.NewUserServiceClient(conn)
	req := pb.UserRequest{ID: 1}
	stream, err := grpcClient.GetUserInfo(ctx, &req)
	if err != nil {
		log.Fatalf("recevie resp error: %v", err)
	}

	header, err := stream.Header()
	if err == nil {
		log.Printf("[RECEIVED HEADER]: %v\n", header)
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
