package auth

import (
	"context"
	"errors"
	"github.com/itohin/gophkeeper/internal/server/entities"
)

type Mailer interface {
	SendMail(to []string, message string) error
}

type UsersStorage interface {
	Save(ctx context.Context, user entities.User) (entities.User, error)
	FindByEmail(ctx context.Context, email string) (entities.User, error)
}

type PasswordHasher interface {
	HashPassword(password string) (string, error)
	IsValidPasswordHash(password, hash string) bool
}

type OTPGenerator interface {
	RandomSecret() (string, error)
}

type UUIDGenerator interface {
	Generate() ([16]byte, error)
}

type AuthUseCase struct {
	hash   PasswordHasher
	uuid   UUIDGenerator
	otp    OTPGenerator
	repo   UsersStorage
	mailer Mailer
}

func NewAuthUseCase(
	hash PasswordHasher,
	uuid UUIDGenerator,
	otp OTPGenerator,
	repo UsersStorage,
	mailer Mailer,
) *AuthUseCase {
	return &AuthUseCase{
		hash:   hash,
		uuid:   uuid,
		otp:    otp,
		repo:   repo,
		mailer: mailer,
	}
}

func (a *AuthUseCase) Login(ctx context.Context, email, password string) (entities.Token, error) {
	token := entities.Token{}

	user, err := a.repo.FindByEmail(ctx, email)
	if err != nil {
		return token, err
	}

	if !a.hash.IsValidPasswordHash(user.Password, password) {
		return token, errors.New("wrong credentials")
	}

	return token, nil
}

func (a *AuthUseCase) Register(ctx context.Context, email, password string) error {

	passwordHash, err := a.hash.HashPassword(password)
	if err != nil {
		return err
	}
	uuid, err := a.uuid.Generate()
	if err != nil {
		return err
	}

	otp, err := a.otp.RandomSecret()
	if err != nil {
		return err
	}
	user := entities.User{
		ID:               uuid,
		Email:            email,
		Password:         passwordHash,
		VerificationCode: otp,
	}

	user, err = a.repo.Save(ctx, user)
	if err != nil {
		return err
	}

	//TODO: async send && retry
	err = a.mailer.SendMail(
		[]string{email},
		"Для подтверждения адреса электронной почты в сервисе gophkeeper, введите пожалуйста код подтверждения: "+otp,
	)
	if err != nil {
		return err
	}

	return nil
}
