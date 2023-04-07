package repo

import (
	"context"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
)

type EventRepo interface {
	LockEvents(ctx context.Context, n uint64) ([]model.UserEvent, error)
	UnlockEvents(ctx context.Context, eventID uint64) error
	RemoveEvents(ctx context.Context, eventID uint64) error
}
