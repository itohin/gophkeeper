package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"os"
	"path/filepath"
)

type PgxPoolDB struct {
	Pool *pgxpool.Pool
}

func NewPgxPoolDB(ctx context.Context, dsn, migrationsPath string) (*PgxPoolDB, error) {
	db := &PgxPoolDB{}
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

	if err = migrate(pool, migrationsPath); err != nil {
		return nil, err
	}

	db.Pool = pool

	return db, nil
}

func migrate(pool *pgxpool.Pool, migrationsPath string) error {
	migrationsPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return err
	}
	goose.SetBaseFS(os.DirFS(migrationsPath))

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	sqlDB := stdlib.OpenDBFromPool(pool)
	if err := goose.Up(sqlDB, "."); err != nil {
		return err
	}
	return nil
}
