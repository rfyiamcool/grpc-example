package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/rfyiamcool/grpc-example/simple/proto"
)

var (
	addr = flag.String("addr", "127.0.0.1:3001", "connect server address")
)

func main() {
	flag.Parse()

	creds, err := credentials.NewClientTLSFromFile("cert/server.cert", "xiaorui.cc") // xiaorui.cc is cert's "Common Name (fully qualified host name) "
	if err != nil {
		log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	}

	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("connect server error: %v", err)
	}
	defer conn.Close()

	// init grpc client
	grpcClient := pb.NewUserServiceClient(conn)

	// create header
	md := metadata.Pairs("my-req-key1", "haha", "my-req-key2", "hello")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	for index := 0; index < 10; index++ {
		var (
			respHeader, trailer metadata.MD
			req                 = pb.UserRequest{ID: int32(index)}
		)

		resp, err := grpcClient.GetUserInfo(
			ctx,
			&req,
			grpc.Header(&respHeader),
			grpc.Trailer(&trailer),
		)
		if err != nil {
			log.Printf("recevie resp error: %v\n", err)
		}

		log.Printf("recv resp: %v\n", resp)
		log.Printf("recv header: %+v\n", respHeader)
		log.Printf("recv trailer: %+v\n", trailer)

		time.Sleep(1 * time.Second)
	}
}
