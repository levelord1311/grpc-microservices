package user_service

import (
	"context"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (i *Implementation) DescribeUser(
	ctx context.Context,
	req *pb.DescribeUserRequest,
) (*pb.DescribeUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	user, err := i.service.DescribeUser(ctx, req.GetId())
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
