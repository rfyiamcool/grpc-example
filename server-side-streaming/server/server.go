package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	pb "github.com/rfyiamcool/grpc-example/server-side-streaming/proto"
)

var users = map[int32]pb.UserResponse{
	0: {Name: "filco", Age: 13},
	1: {Name: "vscode", Age: 70},
	2: {Name: "vim", Age: 75},
	3: {Name: "hhkb", Age: 62},
	4: {Name: "rfyiamcool", Age: 22},
	5: {Name: "rui", Age: 23},
	6: {Name: "emacs", Age: 25},
	7: {Name: "leopold", Age: 21},
}

type serverSideStreamServer struct{}

func (s *serverSideStreamServer) GetUserInfo(req *pb.UserRequest, stream pb.UserService_GetUserInfoServer) error {
	reqHeader, _ := metadata.FromIncomingContext(stream.Context())
	log.Printf("req header: %+v \n", reqHeader)

	gen := func() metadata.MD {
		md := metadata.Pairs("server-key1", "haha", "server-key2", fmt.Sprintf("%s", time.Now().String()))
		return md
	}

	err := stream.SendHeader(gen())
	if err != nil {
		log.Printf("send header err: %+v", err.Error())
		return err
	}

	queue := make(chan pb.UserResponse, 10)
	go func() {
		for _, user := range users {
			queue <- user
			time.Sleep(1 * time.Second)
		}
	}()

	running := true
	for running {
		select {
		case mesg := <-queue:
			stream.Send(&mesg)
			// md := gen()
			// stream.SendHeader(md)

		case <-stream.Context().Done():
			log.Println("client active closed")
			running = false
		}
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
