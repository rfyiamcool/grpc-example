package main

import (
	"context"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/rfyiamcool/grpc-example/simple/proto"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco", Age: 13},
	1: {Name: "vscode", Age: 70},
	2: {Name: "vim", Age: 75},
	3: {Name: "hhkb", Age: 62},
	4: {Name: "rfyiamcool", Age: 22},
	5: {Name: "rui", Age: 23},
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
	addr := "0.0.0.0:3001"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &simpleServer{})
	grpcServer.Serve(listener)
}
