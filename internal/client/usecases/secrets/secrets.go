package secrets

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
)

type Client interface {
	SearchSecrets(ctx context.Context) (map[string]*entities.Secret, error)
	GetSecret(ctx context.Context, id string) (*entities.Secret, error)
	CreateSecret(ctx context.Context, s *entities.Secret) error
	DeleteSecret(ctx context.Context, s *entities.Secret) error
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

func (s *SecretsUseCase) CreateSecret(ctx context.Context, secret *entities.Secret) error {
	err := s.client.CreateSecret(ctx, secret)
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

func (s *SecretsUseCase) DeleteSecret(ctx context.Context, id string) error {
	secret, err := s.storage.GetSecret(ctx, id)
	if err != nil {
		return err
	}
	return s.client.DeleteSecret(ctx, secret)
}
