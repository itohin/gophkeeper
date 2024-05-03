package storage

import (
	"context"
	"fmt"
	"sync"

	"github.com/itohin/gophkeeper/internal/client/entities"
)

type MemoryStorage struct {
	secrets map[string]*entities.Secret
	mx      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		secrets: make(map[string]*entities.Secret),
	}
}

func (m *MemoryStorage) SaveSecrets(ctx context.Context, secrets map[string]*entities.Secret) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	for k, v := range secrets {
		m.secrets[k] = v
	}
	return nil
}

func (m *MemoryStorage) SaveSecret(ctx context.Context, secret *entities.Secret) error {
	m.mx.Lock()
	defer m.mx.Unlock()

	m.secrets[secret.ID] = secret

	return nil
}

func (m *MemoryStorage) GetSecrets(ctx context.Context) (map[string]*entities.Secret, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	s := make(map[string]*entities.Secret, len(m.secrets))
	for id, v := range m.secrets {
		s[id] = v
	}

	return s, nil
}

func (m *MemoryStorage) GetSecret(ctx context.Context, id string) (*entities.Secret, error) {
	m.mx.RLock()
	defer m.mx.RUnlock()

	s, ok := m.secrets[id]
	if !ok {
		return nil, fmt.Errorf("secret ID %v not found", id)
	}
	return s, nil
}

func (m *MemoryStorage) DeleteSecret(ctx context.Context, id string) error {
	m.mx.RLock()
	defer m.mx.RUnlock()

	delete(m.secrets, id)

	return nil
}
