package user_service

import (
	"context"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

func (i *Implementation) RemoveUser(
	ctx context.Context,
	req *pb.RemoveUserRequest,
) (*pb.RemoveUserResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := i.service.RemoveUser(ctx, req.GetId())
	if err != nil {
		log.Println("RemoveUser() error:", err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveUserResponse{}, nil
}
