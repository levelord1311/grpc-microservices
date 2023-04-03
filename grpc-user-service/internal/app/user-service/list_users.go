package user_service

import (
	"context"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) ListUsers(
	ctx context.Context,
	req *pb.ListUsersRequest,
) (*pb.ListUsersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	users, err := i.service.ListUsers(ctx, req.GetLimit(), req.GetCursorId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if users == nil {
		return nil, status.Error(codes.NotFound, "users not found")
	}

	pbUsers := make([]*pb.User, len(users))
	for i, user := range users {
		pbUsers[i] = userToPb(&user)
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
	}, nil
}
