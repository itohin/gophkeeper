package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/database"
	"github.com/itohin/gophkeeper/pkg/events"
)

type SecretsRepository struct {
	db *database.PgxPoolDB
}

func NewSecretsRepository(db *database.PgxPoolDB) *SecretsRepository {
	return &SecretsRepository{db: db}
}

func (r *SecretsRepository) GetUserSecrets(ctx context.Context, userID string) ([]events.SecretDTO, error) {
	secrets := make([]events.SecretDTO, 0)
	query := `SELECT id, user_id, type, name, data, notes FROM secrets WHERE user_id = $1 AND deleted_at IS NULL`
	rows, err := r.db.Pool.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query select secrets: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var secretItem events.SecretDTO
		err = rows.Scan(&secretItem.ID, &secretItem.UserID, &secretItem.SecretType, &secretItem.Name, &secretItem.Data, &secretItem.Notes)
		if err != nil {
			return nil, fmt.Errorf("failed to scan secret row: %v", err)
		}
		secrets = append(secrets, secretItem)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("secrets rows error: %v", err)
	}
	return secrets, nil
}

func (r *SecretsRepository) GetUserSecret(ctx context.Context, userID, secretID string) (events.SecretDTO, error) {
	var s events.SecretDTO
	query := `SELECT id, user_id, type, name, data, notes FROM secrets WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	err := r.db.Pool.QueryRow(ctx, query, secretID, userID).Scan(&s.ID, &s.UserID, &s.SecretType, &s.Name, &s.Data, &s.Notes)
	if err != nil {
		return s, fmt.Errorf("failed to get secret row: %v", err)
	}
	return s, nil
}

func (r *SecretsRepository) Save(ctx context.Context, s entities.Secret) (*events.SecretDTO, error) {
	var sDTO events.SecretDTO
	query := `
		INSERT INTO secrets (
		    id, user_id, type, name, data, notes, created_at, updated_at, deleted_at
		) VALUES (
		    $1, $2, $3, $4, $5, $6, $7, $8, $9
		)
		ON CONFLICT(id, user_id) DO UPDATE set name = $4, data = $5, notes = $6, updated_at = $8, deleted_at = $9
		RETURNING secrets.id, secrets.user_id, secrets.type, secrets.name, secrets.data, secrets.notes
	`

	err := r.db.Pool.QueryRow(
		ctx, query, s.ID, s.UserID, s.SecretType, s.Name, s.Data, s.Notes, time.Now(), time.Now(), s.DeletedAt,
	).Scan(
		&sDTO.ID, &sDTO.UserID, &sDTO.SecretType, &sDTO.Name, &sDTO.Data, &sDTO.Notes,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to save secret row: %v", err)
	}

	return &sDTO, nil
}
