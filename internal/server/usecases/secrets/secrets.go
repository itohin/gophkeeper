package secrets

import (
	"context"
	"github.com/google/uuid"
	"github.com/itohin/gophkeeper/internal/server/entities"
)

type SecretsStorage interface {
	Save(ctx context.Context, secret entities.Secret) error
	Get(ctx context.Context, userID string) ([]entities.SecretDTO, error)
}
type UUIDGenerator interface {
	Generate() ([16]byte, error)
}

type SecretsUseCase struct {
	uuid UUIDGenerator
	repo SecretsStorage
}

func NewSecretsUseCase(uuid UUIDGenerator, repo SecretsStorage) *SecretsUseCase {
	return &SecretsUseCase{
		uuid: uuid,
		repo: repo,
	}
}

func (s *SecretsUseCase) GetUserSecrets(ctx context.Context, userID string) ([]entities.SecretDTO, error) {
	return s.repo.Get(ctx, userID)
}

func (s *SecretsUseCase) Save(ctx context.Context, secret *entities.Secret) (*entities.Secret, error) {
	var secretID uuid.UUID
	secretID, err := s.uuid.Generate()
	if err != nil {
		return nil, err
	}
	secret.ID = secretID.String()
	err = s.repo.Save(ctx, *secret)
	if err != nil {
		return nil, err
	}
	return secret, nil
}
