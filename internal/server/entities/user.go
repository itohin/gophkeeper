package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID               uuid.UUID
	Email            string
	Password         string
	VerificationCode string
	CreatedAt        time.Time
}
