package auth

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

type Auth interface {
	Login(ctx context.Context, login, password string) error
	Register(ctx context.Context, login, password string) error
}

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
