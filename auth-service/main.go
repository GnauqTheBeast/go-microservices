package main

import (
	"fmt"
	"net"

	pb "auth-service/proto"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello world")
	err := InitDB()
	if err != nil {
		fmt.Println(err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println(err)
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &AuthServiceServer{})

	fmt.Println("auth service listening on port 50051")
	if err := s.Serve(lis); err != nil {
		fmt.Println(err)
	}
}
