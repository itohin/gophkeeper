package postgres

import (
	"context"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/database"
	"time"
)

type UsersRepository struct {
	db *database.PgxPoolDB
}

func NewUsersRepository(db *database.PgxPoolDB) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) Save(ctx context.Context, u entities.User) (entities.User, error) {
	var user entities.User
	query := `
		INSERT INTO users (id, email, password, verification_code, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING users.id, users.email
	`
	err := r.db.Pool.QueryRow(ctx, query, u.ID, u.Email, u.Password, u.VerificationCode, time.Now(), time.Now()).Scan(&user.ID, &user.Email)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (entities.User, error) {
	var user entities.User
	query := `SELECT id, email, created_at from users where email = $1`
	err := r.db.Pool.QueryRow(ctx, query, email).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (entities.User, error) {
	var user entities.User
	query := `SELECT id, email, created_at from users where id = $1`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, err
	}

	return user, nil
}
