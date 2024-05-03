package secrets

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/events"
)

type SecretsStorage interface {
	Save(ctx context.Context, secret entities.Secret) (*events.SecretDTO, error)
	GetUserSecrets(ctx context.Context, userID string) ([]events.SecretDTO, error)
	GetUserSecret(ctx context.Context, userID, secretID string) (events.SecretDTO, error)
}
type UUIDGenerator interface {
	Generate() ([16]byte, error)
}

type SecretsUseCase struct {
	uuid    UUIDGenerator
	repo    SecretsStorage
	eventCh chan *events.SecretEvent
}

func NewSecretsUseCase(uuid UUIDGenerator, repo SecretsStorage, eventCh chan *events.SecretEvent) *SecretsUseCase {
	return &SecretsUseCase{
		uuid:    uuid,
		repo:    repo,
		eventCh: eventCh,
	}
}

func (s *SecretsUseCase) GetUserSecrets(ctx context.Context, userID string) ([]events.SecretDTO, error) {
	return s.repo.GetUserSecrets(ctx, userID)
}

func (s *SecretsUseCase) GetUserSecret(ctx context.Context, userID, secretID string) (events.SecretDTO, error) {
	return s.repo.GetUserSecret(ctx, userID, secretID)
}

func (s *SecretsUseCase) DeleteUserSecret(ctx context.Context, secret *entities.Secret) (*events.SecretDTO, error) {
	secret.DeletedAt = sql.NullTime{
		Time:  time.Now(),
		Valid: true,
	}
	dto, err := s.repo.Save(ctx, *secret)

	if err == nil {
		s.sendEvent(dto, events.TypeDeleted)
	}

	return dto, err
}

func (s *SecretsUseCase) Save(ctx context.Context, secret *entities.Secret) (*events.SecretDTO, error) {
	var secretID uuid.UUID

	secretID, err := s.uuid.Generate()
	if err != nil {
		return nil, err
	}
	secret.ID = secretID.String()
	dto, err := s.repo.Save(ctx, *secret)
	if err != nil {
		return nil, err
	}
	s.sendEvent(dto, events.TypeCreated)

	return dto, nil
}

func (s *SecretsUseCase) sendEvent(dto *events.SecretDTO, eventType int) {
	ev := &events.SecretEvent{
		EventType: eventType,
		Secret:    dto,
	}

	s.eventCh <- ev
}
