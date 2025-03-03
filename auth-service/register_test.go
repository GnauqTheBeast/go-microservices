package main

import (
	pb "auth-service/proto/auth"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	// Register user via auth-service
	req := &pb.RegisterRequest{
		Email:    "quang@example.com",
		Password: "secure123",
	}
	res, err := client.Register(context.Background(), req)
	if err != nil {
		log.Fatalf("Registration failed: %v", err)
	}

	fmt.Println("✅ Response:", res.Message)
}
