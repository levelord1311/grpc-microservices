package repo

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/database"
)

const (
	usersTable       = "users"
	usersEventsTable = "users_events"
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

func (r *repo) CreateUser(ctx context.Context, user *model.User) (uint64, error) {
	sb := database.StatementBuilder.
		Insert(usersTable).
		Columns("username", "email").
		Values(user.Username, user.Email).
		Suffix("RETURNING id").
		RunWith(r.db)

	res, err := sb.ExecContext(ctx)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return uint64(id), nil
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

func (r *repo) Lock(ctx context.Context, n uint64) ([]model.UserEvent, error) {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}

	err = database.AcquireTryLock(ctx, tx, database.LockTypeEvent)
	if err != nil {
		return nil, err
	}

	sb := database.StatementBuilder.
		Select("*").
		From(usersEventsTable).
		OrderBy("created").
		Where(sq.Eq{"status": ""}).
		Limit(n)

	selectQuery, args, err := sb.ToSql()
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	var events []model.UserEvent

	if err = tx.SelectContext(ctx, events, selectQuery, args...); err != nil {
		tx.Rollback()
		return nil, err
	}

	ids := make([]uint64, len(events))
	for i, event := range events {
		ids[i] = event.ID
	}

	ub := database.StatementBuilder.
		Update(usersEventsTable).
		Set("status", model.Locked).
		Where(sq.Eq{"id": ids})

	updateQuery, args, err := ub.ToSql()
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	_, err = tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return events, nil
}

func (r *repo) Unlock(ctx context.Context, eventIDs []uint64) error {

	sb := database.StatementBuilder.
		Update(usersEventsTable).
		Set("status", "").
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) Remove(ctx context.Context, eventIDs []uint64) error {
	sb := database.StatementBuilder.
		Update(usersEventsTable).
		Set("status", "removed").
		Where(sq.Eq{"id": eventIDs}).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}
