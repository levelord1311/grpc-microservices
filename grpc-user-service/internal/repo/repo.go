package repo

import (
	"context"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
)

type repo struct {
	db        *sqlx.DB
	batchSize uint64
}

func NewRepo(db *sqlx.DB, batchSize uint64) *repo {
	return &repo{
		db:        db,
		batchSize: batchSize,
	}
}

func (r *repo) CreateUser(ctx context.Context) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) DescribeUser(ctx context.Context, id uint64) (*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) ListUsers(ctx context.Context) ([]*model.User, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repo) RemoveUser(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
