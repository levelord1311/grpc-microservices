package user_service

import (
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service/user"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
)

type Implementation struct {
	pb.UnimplementedUserServiceServer
	service user.Service
}

func NewUserService(userService user.Service) pb.UserServiceServer {
	return &Implementation{service: userService}
}
