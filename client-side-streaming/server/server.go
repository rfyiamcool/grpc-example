package main

import (
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/client-side-streaming/proto"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco", Age: 13},
	1: {Name: "vscode", Age: 70},
	2: {Name: "vim", Age: 75},
	3: {Name: "hhkb", Age: 62},
	4: {Name: "rfyiamcool", Age: 22},
	5: {Name: "rui", Age: 23},
}

type clientSideStreamServer struct{}

func (s *clientSideStreamServer) GetUserInfo(stream pb.UserService_GetUserInfoServer) error {
	var lastID int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			if u, ok := users[lastID]; ok {
				stream.SendAndClose(&u)
				return nil
			}
		}
		lastID = req.ID
		log.Printf("[RECEVIED REQUEST]: %v\n", req)
	}
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
	pb.RegisterUserServiceServer(grpcServer, &clientSideStreamServer{})
	grpcServer.Serve(listener)
}
