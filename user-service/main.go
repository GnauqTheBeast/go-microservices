package main

import (
	"fmt"
	"net"

	"google.golang.org/grpc"
	"user-service/proto/pb"
)

func main() {
	fmt.Println("user-service")
	err := InitDB()
	if err != nil {
		fmt.Printf("init db err: %v\n", err)
	}

	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Printf("lis err: %v\n", err)
	}

	authAddr := "localhost:50051"

	s := grpc.NewServer()

	userServiceServer, err := NewUserServiceServer(authAddr)
	if err != nil {
		fmt.Printf("NewUserServiceServer err: %v\n", err)
	}

	pb.RegisterUserServiceServer(s, userServiceServer)

	fmt.Println("user-service served at port 50052.")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("serve err: %v\n", err)
	}
}
