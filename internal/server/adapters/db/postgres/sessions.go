package postgres

import (
	"context"
	"errors"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/database"
	"time"
)

type SessionsRepository struct {
	db *database.PgxPoolDB
}

func NewSessionsRepository(db *database.PgxPoolDB) *SessionsRepository {
	return &SessionsRepository{db: db}
}

func (r *SessionsRepository) Save(ctx context.Context, s entities.Session) error {
	query := `
		INSERT INTO sessions (
		    id, user_id, fingerprint, expires_at, created_at, updated_at
		) VALUES (
		    $1, $2, $3, $4, $5, $6
		)
		ON CONFLICT(id) DO UPDATE set user_id = $2, fingerprint = $3, expires_at = $4, updated_at = $6
		RETURNING sessions.id
	`
	result, err := r.db.Pool.Exec(
		ctx, query, s.ID, s.UserID, s.FingerPrint, s.ExpiresAt, time.Now(), time.Now(),
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

func (r *SessionsRepository) FindByID(ctx context.Context, id string) (*entities.Session, error) {
	var session entities.Session
	query := `SELECT id, user_id, fingerprint, expires_at from sessions where id = $1`
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(&session.ID, &session.UserID, &session.FingerPrint, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionsRepository) FindByFingerPrint(ctx context.Context, userId, fingerPrint string) (*entities.Session, error) {
	var session entities.Session
	query := `SELECT id, user_id, fingerprint, expires_at from sessions where user_id = $1 AND fingerprint $2`
	err := r.db.Pool.QueryRow(ctx, query, userId, fingerPrint).Scan(&session.ID, &session.UserID, &session.FingerPrint, &session.ExpiresAt)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func (r *SessionsRepository) DeleteByUserAndFingerPrint(ctx context.Context, userId, fingerPrint string) error {
	query := `DELETE from sessions where user_id = $1 and fingerprint = $2`
	_, err := r.db.Pool.Exec(ctx, query, userId, fingerPrint)
	if err != nil {
		return err
	}

	return nil
}

func (r *SessionsRepository) DeleteByID(ctx context.Context, sessionID string) error {
	query := `DELETE from sessions where id = $1`
	result, err := r.db.Pool.Exec(ctx, query, sessionID)
	if err != nil {
		return err
	}
	count := result.RowsAffected()
	if count < 1 {
		return errors.New("failed to delete record from database")
	}

	return nil
}
