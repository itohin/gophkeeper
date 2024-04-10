package secrets

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

type Client interface {
	CreateText(ctx context.Context, secret *entities.Secret, text string) error
	CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error
	SearchSecrets(ctx context.Context) (map[string]*entities.Secret, error)
	GetSecret(ctx context.Context, id string) (*entities.Secret, error)
}

type Storage interface {
	SaveSecrets(ctx context.Context, secrets map[string]*entities.Secret) error
	SaveSecret(ctx context.Context, secret *entities.Secret) error
	GetSecrets(ctx context.Context) (map[string]*entities.Secret, error)
	GetSecret(ctx context.Context, id string) (*entities.Secret, error)
}

type SecretsUseCase struct {
	client  Client
	storage Storage
}

func NewSecrets(client Client, storage Storage) *SecretsUseCase {
	return &SecretsUseCase{
		client:  client,
		storage: storage,
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

func (s *SecretsUseCase) GetSecrets(ctx context.Context) (map[string]*entities.Secret, error) {
	return s.storage.GetSecrets(ctx)
}

func (s *SecretsUseCase) GetSecret(ctx context.Context, id string) (*entities.Secret, error) {
	return s.storage.GetSecret(ctx, id)
}

func (s *SecretsUseCase) SaveSecret(ctx context.Context, secret *entities.Secret) error {
	return s.storage.SaveSecret(ctx, secret)
}

func (s *SecretsUseCase) SyncSecrets(ctx context.Context) error {
	secrets, err := s.client.SearchSecrets(ctx)
	if err != nil {
		return err
	}
	return s.storage.SaveSecrets(context.Background(), secrets)
}
