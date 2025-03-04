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

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &UserServiceServer{})

	fmt.Println("user-service served at port 50052.")
	if err := s.Serve(lis); err != nil {
		fmt.Printf("serve err: %v\n", err)
	}
}
