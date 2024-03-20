package secrets

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/server/entities"
)

type SecretsStorage interface {
	Save(ctx context.Context, user entities.Secret) error
}

type SecretsUseCase struct {
}

func NewSecretsUseCase() *SecretsUseCase {
	return &SecretsUseCase{}
}

func (s *SecretsUseCase) Store(ctx context.Context, secret *entities.Secret) error {
	fmt.Println("store secret", secret)
	return nil
}
