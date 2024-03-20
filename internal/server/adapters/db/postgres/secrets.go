package postgres

import (
	"context"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/database"
)

type SecretsRepository struct {
	db *database.PgxPoolDB
}

func NewSecretsRepository(db *database.PgxPoolDB) *SecretsRepository {
	return &SecretsRepository{db: db}
}

func (r *SecretsRepository) Save(ctx context.Context, s entities.Secret) error {
	//query := `
	//	INSERT INTO secrets (
	//	    id, user_id, fingerprint, expires_at, created_at, updated_at
	//	) VALUES (
	//	    $1, $2, $3, $4, $5, $6
	//	)
	//	ON CONFLICT(id) DO UPDATE set user_id = $2, fingerprint = $3, expires_at = $4, updated_at = $6
	//	RETURNING sessions.id
	//`
	//result, err := r.db.Pool.Exec(
	//	ctx, query, s.ID, s.UserID, s.FingerPrint, s.ExpiresAt, time.Now(), time.Now(),
	//)
	//if err != nil {
	//	return err
	//}
	//count := result.RowsAffected()
	//if count < 1 {
	//	return errors.New("failed to update record in database")
	//}

	return nil
}
