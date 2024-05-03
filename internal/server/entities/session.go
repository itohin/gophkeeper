package entities

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	FingerPrint string
	ExpiresAt   time.Time
}

func NewSession(id, userId [16]byte, fingerprint string, expires time.Time) *Session {
	return &Session{
		ID:          id,
		UserID:      userId,
		FingerPrint: fingerprint,
		ExpiresAt:   expires,
	}
}

func (s *Session) IsExpired() bool {
	return s.ExpiresAt.Before(time.Now())
}
