package database

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PgxPoolDB is a pool of connections to postgres db
type PgxPoolDB struct {
	Pool *pgxpool.Pool
}

// NewPgxPoolDB is a PgxPoolDB constructor that creates metrics table if not exists
func NewPgxPoolDB(ctx context.Context, dsn string) (*PgxPoolDB, error) {
	if dsn == "" {
		return nil, errors.New("no database dsn specified")
	}
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, err
	}

	//TODO migration

	if err != nil {
		return nil, err
	}

	return &PgxPoolDB{
		Pool: pool,
	}, nil
}
