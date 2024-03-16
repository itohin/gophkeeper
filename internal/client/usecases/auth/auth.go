package auth

import (
	"context"
)

type Client interface {
	Register(ctx context.Context, email, password string) error
	Verify(ctx context.Context, email, otp string) error
	Login(ctx context.Context, email, password string) error
	Logout(ctx context.Context) error
}

type AuthUseCase struct {
	client Client
}

func NewAuth(client Client) *AuthUseCase {
	return &AuthUseCase{
		client: client,
	}
}

func (a *AuthUseCase) Register(ctx context.Context, login, password string) error {
	return a.client.Register(ctx, login, password)
}

func (a *AuthUseCase) Verify(ctx context.Context, login, otp string) error {
	err := a.client.Verify(ctx, login, otp)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) Login(ctx context.Context, login, password string) error {
	err := a.client.Login(ctx, login, password)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthUseCase) Logout(ctx context.Context) error {
	return a.client.Logout(ctx)
}
