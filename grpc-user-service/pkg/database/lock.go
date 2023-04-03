package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type lockType int

var ErrFailedToAcquireLock = errors.New("failed to acquire lock")

const (
	_ lockType = iota
	LockTypeEvent
)

func AcquireTryLock(ctx context.Context, tx *sqlx.Tx, lockID lockType) error {
	var isAcquired bool
	err := tx.GetContext(ctx, &isAcquired, fmt.Sprintf("select pg_try_advisory_xact_lock(%d)", lockID))
	if err != nil {
		return errors.Wrap(err, "tx.GetContext()")
	}
	if !isAcquired {
		return ErrFailedToAcquireLock
	}
	return nil
}
