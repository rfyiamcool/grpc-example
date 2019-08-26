package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/rfyiamcool/grpc-example/bidirectional-streaming/proto"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco", Age: 13},
	1: {Name: "vscode", Age: 70},
	2: {Name: "vim", Age: 75},
	3: {Name: "hhkb", Age: 62},
	4: {Name: "rfyiamcool", Age: 22},
	5: {Name: "rui", Age: 23},
}

type bidirectionalStreamServer struct{}

func (s *bidirectionalStreamServer) GetUserInfo(stream pb.UserService_GetUserInfoServer) error {
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}

		if err != nil {
			return err
		}

		u, ok := users[req.ID]
		if !ok {
			u = pb.UserResponse{Name: "unknow", Age: 99}
		}
		err = stream.Send(&u)
		if err != nil {
			return err
		}

		log.Printf("[RECEVIED REQUEST]: %v\n", req)
	}
}

var (
	port = flag.Int("port", 3001, "listening port")
)

func main() {
	flag.Parse()
	addr := fmt.Sprintf("0.0.0.0:%d", *port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	} else {
		log.Println("server listen: ", addr)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, &bidirectionalStreamServer{})
	grpcServer.Serve(listener)
}
