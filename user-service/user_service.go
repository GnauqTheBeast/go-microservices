package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"

	pb "user-service/proto/pb"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
	authClient pb.AuthServiceClient
}

func NewUserServiceServer(authAddr string) (*UserServiceServer, error) {
	conn, err := grpc.NewClient(authAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to auth service: %v", err)
	}

	return &UserServiceServer{
		authClient: pb.NewAuthServiceClient(conn),
	}, nil
}

func (s *UserServiceServer) GetUser(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	var email, name string
	err := db.QueryRow("SELECT email, name FROM users WHERE id=$1", req.UserId).Scan(&email, &name)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &pb.UserResponse{Email: email, Name: name}, nil
}

func (s *UserServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	authRes, err := s.authClient.Login(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("auth login failed: %v", err)
	}

	return &pb.LoginResponse{
		Token: authRes.Token,
	}, nil
}
