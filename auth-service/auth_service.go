package main

import (
	"context"
	"database/sql"

	pb "auth-service/proto/auth"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceServer struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthServiceServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("INSERT INTO users (email, password) VALUES ($1, $2)", req.Email, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *AuthServiceServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	var storedPassword string
	err := db.QueryRow("SELECT password FROM users WHERE email=$1", req.Email).Scan(&storedPassword)
	if err == sql.ErrNoRows {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(req.Password))
	if err != nil {
		return nil, err
	}

	token := "fake-jwt-token"
	return &pb.LoginResponse{Token: token}, nil
}
