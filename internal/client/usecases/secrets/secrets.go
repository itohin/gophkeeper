package secrets

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

type Client interface {
	CreateText(ctx context.Context, secret *entities.Secret, text string) error
	CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error
}

type SecretsUseCase struct {
	client Client
}

func NewSecrets(client Client) *SecretsUseCase {
	return &SecretsUseCase{
		client: client,
	}
}

func (s *SecretsUseCase) CreateText(ctx context.Context, secret *entities.Secret, text string) error {
	err := s.client.CreateText(ctx, secret, text)
	if err != nil {
		return err
	}
	return nil
}

func (s *SecretsUseCase) CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error {
	err := s.client.CreatePassword(ctx, secret, password)
	if err != nil {
		return err
	}
	return nil
}
