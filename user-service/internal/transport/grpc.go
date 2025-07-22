package transport

import (
	"context"
	"time"

	userpb "github.com/best-microservice/common/protos/user"
	"github.com/best-microservice/user-service/internal/models"
	"github.com/best-microservice/user-service/internal/service"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UserServer struct {
	userpb.UnimplementedUserServiceServer
	service *service.UserService
}

func NewUserServer(service *service.UserService) *UserServer {
	return &UserServer{service: service}
}

func (s *UserServer) CreateUser(ctx context.Context, req *userpb.CreateUserRequest) (*userpb.UserResponse, error) {
	newUser := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	err := s.service.CreateUser(ctx, newUser)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &userpb.UserResponse{
		Id:        newUser.ID,
		Name:      newUser.Name,
		Email:     newUser.Email,
		CreatedAt: newUser.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserServer) GetUser(ctx context.Context, req *userpb.GetUserRequest) (*userpb.UserResponse, error) {
	user, err := s.service.GetUser(ctx, req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found: %v", err)
	}

	return &userpb.UserResponse{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *UserServer) AuthenticateUser(ctx context.Context, req *userpb.AuthRequest) (*userpb.AuthResponse, error) {
	user, err := s.service.Authenticate(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "authentication failed: %v", err)
	}

	// In a real implementation, generate a proper JWT token
	token := "generated-jwt-token"

	return &userpb.AuthResponse{
		Token: token,
		User: &userpb.UserResponse{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
		},
	}, nil
}
