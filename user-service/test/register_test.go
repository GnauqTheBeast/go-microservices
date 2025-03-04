package test

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"user-service/proto/pb"
)

func TestRegister(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewAuthServiceClient(conn)

	req := &pb.RegisterRequest{
		Email:    "quang1234@example.com",
		Password: "secure123",
	}

	res, err := client.Register(context.Background(), req)
	if err != nil {
		fmt.Printf("Registration failed: %v", err)
	}

	fmt.Println("✅ Response from user-service:", res.Message)
}
