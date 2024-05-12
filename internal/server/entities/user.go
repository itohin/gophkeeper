package entities

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               uuid.UUID
	Email            string
	Password         string
	VerificationCode string
	VerifiedAt       sql.NullTime
	Version          int
}

func NewUser(id [16]byte, email, password, otp string) *User {
	return &User{
		ID:               id,
		Email:            email,
		Password:         password,
		VerificationCode: otp,
	}
}

func (u *User) Verify(otp string) error {
	if u.VerifiedAt.Valid {
		return errors.New("already verified")
	}
	if otp != u.VerificationCode {
		return errors.New("wrong verification code")
	}
	u.VerifiedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	return nil
}

func (u *User) IsVerified() bool {
	return u.VerifiedAt.Valid
}
