package test

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"testing"
	"user-service/proto/pb"
)

func TestLogin(t *testing.T) {
	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	req := &pb.LoginRequest{
		Email:    "quang1234@example.com",
		Password: "secure123",
	}

	_, _ = client.GetUser(context.Background(), req)

	//fmt.Println("✅ Response from user-service:", res.Message, "Token:", res.Token)
}
