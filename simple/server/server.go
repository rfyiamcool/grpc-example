package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "github.com/rfyiamcool/grpc-example/simple/proto"
)

var (
	addr = "0.0.0.0:3001"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco......", Age: 13},
	1: {Name: "vscode.....", Age: 70},
	2: {Name: "vim......", Age: 75},
	3: {Name: "hhkb......", Age: 62},
}

type simpleServer struct{}

func (s *simpleServer) GetUserInfo(ctx context.Context, req *pb.UserRequest) (resp *pb.UserResponse, err error) {
	reqHeader, _ := metadata.FromIncomingContext(ctx)
	user, ok := users[req.ID]
	if !ok {
		resp = &pb.UserResponse{
			Name: "null",
			Age:  0,
		}
		return resp, grpc.Errorf(999, "not found user id")
	}

	resp = &user

	// Create and send header.
	header := metadata.New(
		map[string]string{"my-resp-location": "beijing", "my-resp-ts": time.Now().String()},
	)
	grpc.SendHeader(ctx, header)

	log.Printf("recv request: %+v\n", req.ID)
	log.Printf("recv header: %+v\n", reqHeader)

	return
}

func main() {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}

	creds, err := credentials.NewServerTLSFromFile("cert/server.cert", "cert/server.key")
	if err != nil {
		log.Fatalf("tls error: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.Creds(creds))
	pb.RegisterUserServiceServer(grpcServer, &simpleServer{})
	grpcServer.Serve(listener)
}
