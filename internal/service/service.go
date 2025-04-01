package service

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "grpc-auth/internal/protos"
	"grpc-auth/internal/repo"
)

// Реализация сервиса AuthService
type AuthService struct {
	pb.UnimplementedAuthServiceServer
	repo repo.Repository
	log  *zap.SugaredLogger
}

// NewAuthService - конструктор для AuthService
func NewAuthService(r repo.Repository, l *zap.SugaredLogger) *AuthService {
	return &AuthService{
		repo: r,
		log:  l,
	}
}

func (s *AuthService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	exists, err := s.repo.CheckUserExists(ctx, req.Username)
	if err != nil {
		s.log.Errorf("Failed to check user existence for username '%s': %v", req.Username, err)
		return nil, status.Errorf(codes.Internal, "failed to check user existence: %v", err)
	}
	if exists {
		s.log.Warnf("Registration attempt for existing username: '%s'", req.Username)
		return nil, status.Errorf(codes.AlreadyExists, "user already exists")
	}

	hashPass, err := bcrypt.GenerateFromPassword([]byte(req.GetPassword()), bcrypt.DefaultCost)

	if err != nil {
		s.log.Error("failed to hash password", zap.Error(err))
		return nil, status.Error(codes.InvalidArgument, "failed to generate hash")
	}
	s.log.Infof("started registering a user")

	errReg := s.repo.RegisterUser(ctx, repo.User{
		Email:     req.GetEmail(),
		Username:  req.GetUsername(),
		PassHash:  string(hashPass),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
	})
	if errReg != nil {
		s.log.Errorf("failed to register user: %v", errReg)
		return nil, status.Error(codes.Internal, "failed to register user")
	}
	s.log.Infof("User registered successfully: %s", req.Username)

	return &pb.RegisterResponse{Message: fmt.Sprintf("User registered successfully: %s", req.Username)}, nil
}

func (s *AuthService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := s.repo.GetUser(ctx, req.GetUsername())
	if err != nil {
		s.log.Errorf("User not found: %s", err)
		return nil, status.Error(codes.NotFound, "User not found")
	}

	errPass := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(req.GetPassword()))
	if errPass != nil {
		s.log.Errorf("invalid password for user: %s", user.Username)
		return nil, status.Error(codes.Unauthenticated, "invalid password")
	}

	errTime := s.repo.UpdateLoginTime(ctx, req.GetUsername())
	if errTime != nil {
		s.log.Errorf("Failed to update login time: %s", errTime)
	}

	s.log.Infof("Login successfully: %s", req.Username)
	return &pb.LoginResponse{Token: "token123(temp)"}, nil
}
