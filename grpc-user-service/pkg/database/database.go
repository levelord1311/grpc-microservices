package database

import (
	"context"
	sq "github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"time"
)

type Config interface {
	GetDSN() string
	GetMaxOpenConns() int
	GetMaxIdleConns() int
	GetConnMaxIdleTime() time.Duration
	GetConnMaxLifetime() time.Duration
}

var StatementBuilder = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)

func ConnectDB(ctx context.Context, cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("pgx", cfg.GetDSN())

	if err != nil {
		err = errors.Wrap(err, "sql.Open()")
		return nil, err
	}

	db.SetMaxOpenConns(cfg.GetMaxOpenConns())
	db.SetMaxIdleConns(cfg.GetMaxIdleConns())
	db.SetConnMaxIdleTime(cfg.GetConnMaxIdleTime())
	db.SetConnMaxLifetime(cfg.GetConnMaxLifetime())

	if err = db.PingContext(ctx); err != nil {
		err = errors.Wrap(err, "database.PingContext()")
		return nil, err
	}

	return db, nil
}
