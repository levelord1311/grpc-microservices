package service

import (
	"context"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Repo interface {
	CreateUser(ctx context.Context) (*model.User, error)
	DescribeUser(ctx context.Context, id uint64) (*model.User, error)
	ListUsers(ctx context.Context) ([]*model.User, error)
	RemoveUser(ctx context.Context) error
}

type userService struct {
	pb.UnimplementedUserServiceServer
	repo Repo
}

func NewUserService(r Repo) pb.UserServiceServer {
	return &userService{repo: r}
}

func (s *userService) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.repo.CreateUser(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	pbUser := userToPb(user)
	return &pb.CreateUserResponse{
		User: pbUser,
	}, nil
}

func (s *userService) DescribeUser(
	ctx context.Context,
	req *pb.DescribeUserRequest,
) (*pb.DescribeUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := s.repo.DescribeUser(ctx, req.GetId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if user == nil {
		return nil, status.Error(codes.NotFound, "user not found")
	}

	pbUser := userToPb(user)
	return &pb.DescribeUserResponse{
		User: pbUser,
	}, nil
}

func (s *userService) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, err := s.repo.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if users == nil {
		return nil, status.Error(codes.NotFound, "users not found")
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = userToPb(user)
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
	}, nil
}

func (s *userService) RemoveUser(
	ctx context.Context,
	req *pb.RemoveUserRequest,
) (*pb.RemoveUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.repo.RemoveUser(ctx); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveUserResponse{
		Result: true,
	}, nil
}

func userToPb(u *model.User) *pb.User {
	if u == nil {
		return nil
	}
	return &pb.User{
		Id:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Name:     u.Name,
		Surname:  u.Surname,
	}
}
