package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxPoolTransaction struct {
	db *pgxpool.Pool
}

func NewPgxTransaction(db *pgxpool.Pool) *PgxPoolTransaction {
	return &PgxPoolTransaction{db: db}
}

func (t *PgxPoolTransaction) Transaction(ctx context.Context, f func() error) error {
	tx, err := t.db.Begin(ctx)
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	err = f()
	if err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}
	return nil
}
