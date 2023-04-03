package user_service

import (
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	pb "github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/user-service-api"
)

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
