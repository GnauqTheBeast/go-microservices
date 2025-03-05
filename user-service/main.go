package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"time"
	"user-service/proto/pb"
)

const (
	authAddr = "localhost:50051"
)

func main() {
	fmt.Println("user-service")
	err := InitDB()
	if err != nil {
		fmt.Printf("init db err: %v\n", err)
	}

	//go startGRPC()

	conn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("failed to connect to auth service", err)
	}

	authClient := pb.NewAuthServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	req := &pb.LoginRequest{
		Email:    "quang@example.com",
		Password: "secure123",
	}

	login, err := authClient.Login(ctx, req)
	if err != nil {
		fmt.Printf("Login error: %v\n", err)
	}
	fmt.Println(login)
}

func startGRPC() {
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		fmt.Printf("lis err: %v\n", err)
	}

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
