package main

import (
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/server-side-streaming/proto"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco", Age: 13},
	1: {Name: "vscode", Age: 70},
	2: {Name: "vim", Age: 75},
	3: {Name: "hhkb", Age: 62},
	4: {Name: "rfyiamcool", Age: 22},
	5: {Name: "rui", Age: 23},
}

type serverSideStreamServer struct{}

func (s *serverSideStreamServer) GetUserInfo(req *pb.UserRequest, stream pb.UserService_GetUserInfoServer) error {
	for _, user := range users {
		stream.Send(&user)
	}
	log.Printf("[RECEIVED REQUEST]: %v\n", req)
	return nil
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

	pb.RegisterUserServiceServer(grpcServer, &serverSideStreamServer{})
	grpcServer.Serve(listener)
}
