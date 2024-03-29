package user_service

import (
	"context"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) CreateUser(
	ctx context.Context,
	req *pb.CreateUserRequest,
) (*pb.CreateUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	user := reqToUser(req)

	id, err := i.service.CreateUser(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateUserResponse{
		Id: id,
	}, nil
}

func reqToUser(req *pb.CreateUserRequest) *model.User {
	return &model.User{
		Username: req.Username,
		Email:    req.Email,
	}
}
