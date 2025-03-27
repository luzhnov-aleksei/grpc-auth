package service

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpc-auth/internal/proto"
	"net"
)

// Реализация сервиса AuthService
type AuthService struct {
	pb.UnimplementedAuthServiceServer
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	var users UserRequest
	if _, exists := users[req.Username]; exists {
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}
	users[req.Username] = req.Password
	return &pb.RegisterResponse{Message: "User registered successfully"}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	fmt.Println("Received request for user:", req.Username)
	return &pb.LoginResponse{Name: "John Doe", Age: 30}, nil
}
