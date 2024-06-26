package auth

import (
	"context"
)

type Client interface {
	Register(ctx context.Context, email, password string) error
	Verify(ctx context.Context, email, otp string) (string, error)
	Login(ctx context.Context, email, password string) (string, error)
	Logout(ctx context.Context) error
}

type ServerListener interface {
	Listen(ctx context.Context, userID string) error
}

type AuthUseCase struct {
	client Client
	authCh chan string
}

func NewAuth(client Client, authCh chan string) *AuthUseCase {
	return &AuthUseCase{
		client: client,
		authCh: authCh,
	}
}

func (a *AuthUseCase) Register(ctx context.Context, login, password string) error {
	return a.client.Register(ctx, login, password)
}

func (a *AuthUseCase) Verify(ctx context.Context, login, otp string) error {
	userID, err := a.client.Verify(ctx, login, otp)
	if err != nil {
		return err
	}

	a.authCh <- userID
	return nil
}

func (a *AuthUseCase) Login(ctx context.Context, login, password string) error {
	userID, err := a.client.Login(ctx, login, password)
	if err != nil {
		return err
	}

	a.authCh <- userID
	return nil
}

func (a *AuthUseCase) Logout(ctx context.Context) error {
	return a.client.Logout(ctx)
}
