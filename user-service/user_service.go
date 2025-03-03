package main

import (
	"context"
	"log"

	pb "user-service/proto/user"
)

type UserServiceServer struct {
	pb.UnimplementedUserServiceServer
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
