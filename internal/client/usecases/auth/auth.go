package auth

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

type AuthUseCase struct {
	token *entities.Token
}

func NewAuth() *AuthUseCase {
	return &AuthUseCase{}
}

func (a *AuthUseCase) Login(ctx context.Context, login, password string) error {
	return nil
}

func (a *AuthUseCase) Register(ctx context.Context, login, password string) error {
	return nil
}

func (a *AuthUseCase) Verify(ctx context.Context, login, otp string) error {
	return nil
}
