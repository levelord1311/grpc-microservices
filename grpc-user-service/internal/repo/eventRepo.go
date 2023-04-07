package repo

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/internal/model"
	"github.com/levelord1311/grpc-microservices/grpc-user-service/pkg/database"
	"github.com/pkg/errors"
	"log"
)

func (r *repo) LockEvents(ctx context.Context, n uint64) (events []model.UserEvent, txErr error) {
	txErr = database.WithTx(ctx, r.db, func(ctx context.Context, tx *sqlx.Tx) error {
		err := database.AcquireTryLock(ctx, tx, database.LockTypeEvent)
		if err != nil {
			return err
		}
		events, err = selectEventsForLock(ctx, tx, n)
		if err != nil {
			return err
		}

		ids := make([]uint64, len(events))
		for i, event := range events {
			ids[i] = event.ID
		}

		if err = lockSelectedIDs(ctx, tx, ids); err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return events, nil
}

func (r *repo) UnlockEvents(ctx context.Context, eventID uint64) error {
	log.Println("UnlockEvents()")
	sb := database.StatementBuilder.
		Update(usersEventsTable).
		Set("status", "").
		Where(sq.Eq{"id": eventID}).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (r *repo) RemoveEvents(ctx context.Context, eventID uint64) error {
	log.Println("RemoveEvents()")
	sb := database.StatementBuilder.
		Update(usersEventsTable).
		Set("deleted", "true").
		Where(sq.Eq{"id": eventID}).
		RunWith(r.db)

	_, err := sb.ExecContext(ctx)
	if err != nil {
		return err
	}

	return nil
}

func selectEventsForLock(ctx context.Context, tx *sqlx.Tx, n uint64) ([]model.UserEvent, error) {
	sb := database.StatementBuilder.
		Select("id", "user_id", "event_type", "payload").
		From(usersEventsTable).
		OrderBy("created").
		Where("locked IS NOT true").
		Limit(n)

	selectQuery, args, err := sb.ToSql()
	if err != nil {
		return nil, err
	}

	var events []model.UserEvent

	rows, err := tx.QueryxContext(ctx, selectQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "SelectContext()")
	}
	for rows.Next() {
		var event model.UserEvent
		if err = rows.StructScan(&event); err != nil {
			return nil, errors.Wrap(err, "rows.Scan()")
		}
		events = append(events, event)
	}
	if len(events) == 0 {
		return nil, ErrNothingToLock
	}

	return events, nil
}

func lockSelectedIDs(ctx context.Context, tx *sqlx.Tx, IDs []uint64) error {
	sb := database.StatementBuilder.
		Update(usersEventsTable).
		Set("locked", "TRUE").
		Where(sq.Eq{"id": IDs})

	updateQuery, args, err := sb.ToSql()
	if err != nil {
		return err
	}
	_, err = tx.ExecContext(ctx, updateQuery, args...)
	if err != nil {
		return err
	}
	return nil
}
