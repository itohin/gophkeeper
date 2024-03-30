package postgres

import (
	"context"
	"errors"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/database"
	"time"
)

type SecretsRepository struct {
	db *database.PgxPoolDB
}

func NewSecretsRepository(db *database.PgxPoolDB) *SecretsRepository {
	return &SecretsRepository{db: db}
}

func (r *SecretsRepository) Get(ctx context.Context, userID string) ([]entities.SecretDTO, error) {
	secrets := make([]entities.SecretDTO, 0)
	query := `SELECT id, user_id, type, name, data, notes FROM secrets where user_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var secretItem entities.SecretDTO
		err = rows.Scan(&secretItem.ID, &secretItem.UserID, &secretItem.SecretType, &secretItem.Name, &secretItem.Data, &secretItem.Notes)
		if err != nil {
			return nil, err
		}
		secrets = append(secrets, secretItem)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return secrets, nil
}

func (r *SecretsRepository) Save(ctx context.Context, s entities.Secret) error {
	query := `
		INSERT INTO secrets (
		    id, user_id, type, name, data, notes, created_at, updated_at
		) VALUES (
		    $1, $2, $3, $4, $5, $6, $7, $8
		)
		ON CONFLICT(id) DO UPDATE set user_id = $2, type = $3, name = $4, data = $5, notes = $6, updated_at = $8
		RETURNING secrets.id
	`
	result, err := r.db.Pool.Exec(
		ctx, query, s.ID, s.UserID, s.SecretType, s.Name, s.Data, s.Notes, time.Now(), time.Now(),
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
