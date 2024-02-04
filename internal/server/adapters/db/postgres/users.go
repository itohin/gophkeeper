package postgres

import (
	"context"
	"github.com/itohin/gophkeeper/pkg/database"
)

type UsersRepository struct {
	db *database.PgxPoolDB
}

func NewUsersRepository(db *database.PgxPoolDB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Store(ctx context.Context, login, password string) error {
	return nil
}
