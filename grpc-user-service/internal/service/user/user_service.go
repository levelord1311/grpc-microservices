package user

import (
	"context"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
)

type Repo interface {
	CreateUser(ctx context.Context, user *model.User) (uint64, error)
	DescribeUser(ctx context.Context, id uint64) (*model.User, error)
	ListUsers(ctx context.Context, limit, cursor uint64) ([]model.User, error)
	RemoveUser(ctx context.Context, id uint64) error
}

type Service struct {
	repo Repo
}

func NewService(repo Repo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) CreateUser(ctx context.Context, user *model.User) (uint64, error) {
	return s.repo.CreateUser(ctx, user)
}

func (s *Service) DescribeUser(ctx context.Context, id uint64) (*model.User, error) {
	return s.repo.DescribeUser(ctx, id)
}

func (s *Service) ListUsers(ctx context.Context, limit, cursor uint64) ([]model.User, error) {
	return s.repo.ListUsers(ctx, limit, cursor)
}

func (s *Service) RemoveUser(ctx context.Context, id uint64) error {
	return s.repo.RemoveUser(ctx, id)
}
