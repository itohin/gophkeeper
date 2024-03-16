package postgres

import (
	"context"
	"errors"
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

func (r *UsersRepository) Save(ctx context.Context, u entities.User) error {
	query := `
		INSERT INTO users (
		    id, email, password, verification_code, verified_at, version, created_at, updated_at
		) VALUES (
		    $1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT(id) DO UPDATE 
		    set email = $2, 
		    password = $3, 
		    verification_code = $4, 
		    verified_at = $5, 
		    updated_at = $8,
		    version = $9
		WHERE users.version = $10
		RETURNING users.id, users.email
	`
	result, err := r.db.Pool.Exec(
		ctx, query, u.ID, u.Email, u.Password, u.VerificationCode, u.VerifiedAt, 1, time.Now(), time.Now(), u.Version+1, u.Version,
	)
	if err != nil {
		return err
	}
	count := result.RowsAffected()
	if count < 1 {
		return errors.New("failed to update record in database")
	}
	return nil
}

func (r *UsersRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, email, password, verification_code, verified_at, version from users where email = $1`
	err := r.db.Pool.QueryRow(ctx, query, email).
		Scan(&user.ID, &user.Email, &user.Password, &user.VerificationCode, &user.VerifiedAt, &user.Version)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UsersRepository) FindByID(ctx context.Context, id string) (*entities.User, error) {
	var user entities.User
	query := `SELECT id, email, created_at from users where id = $1`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&user.ID, &user.Email)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
