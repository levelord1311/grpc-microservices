package repo

import (
	"context"
	"errors"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/service/user"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/database"
)

var (
	_ user.Repo = &repo{}
)

var (
	ErrNothingToLock = errors.New("no events to lock")
)

const (
	usersTable       = "users"
	usersEventsTable = "users_events"
)

type repo struct {
	db *sqlx.DB
}

func NewRepo(db *sqlx.DB) *repo {
	return &repo{
		db: db,
	}
}

func (r *repo) CreateUser(ctx context.Context, user *model.User) (id uint64, txErr error) {
	txErr = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		createdId, err := r.createUser(ctx, user)
		if err != nil {
			return err
		}
		id = createdId
		if err = r.createUserEvent(ctx, model.Created, id, user); err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return 0, txErr
	}
	return id, nil

}

func (r *repo) DescribeUser(ctx context.Context, id uint64) (*model.User, error) {
	sb := database.StatementBuilder.
		Select("*").
		From(usersTable).
		Where(sq.Eq{"id": id})

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}
	var user *model.User
	err = r.db.SelectContext(ctx, user, query, args...)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *repo) ListUsers(ctx context.Context, limit, cursor uint64) ([]model.User, error) {
	sb := database.StatementBuilder.
		Select("*").
		From(usersTable).
		Limit(limit).
		Where(sq.Lt{"id": cursor})

	query, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var res []model.User
	err = r.db.SelectContext(ctx, &res, query, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (r *repo) RemoveUser(ctx context.Context, id uint64) error {
	txErr := database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		if err := r.removeUser(ctx, id); err != nil {
			return err
		}
		if err := r.createUserEvent(ctx, model.Removed, id, nil); err != nil {
			return err
		}
		return nil
	})
	return txErr
}

func (r *repo) createUserEvent(ctx context.Context, eventType model.EventType, id uint64, user *model.User) error {
	sb := database.StatementBuilder.
		Insert(usersEventsTable).
		Columns("user_id", "event_type", "payload").
		Values(id, eventType, user).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	return err
}

func (r *repo) createUser(ctx context.Context, user *model.User) (id uint64, err error) {
	sb := database.StatementBuilder.
		Insert(usersTable).
		Columns("username", "email").
		Values(user.Username, user.Email).
		Suffix("RETURNING id")

	query, args, err := sb.ToSql()
	if err != nil {
		return 0, err
	}
	row := r.db.QueryRowContext(ctx, query, args...)
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *repo) removeUser(ctx context.Context, id uint64) error {
	sb := database.StatementBuilder.
		Update(usersTable).
		Set("removed", true).
		Where(sq.Eq{"id": id}).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
