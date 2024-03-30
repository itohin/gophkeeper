package storage

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"sync"
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
	s, ok := m.secrets[id]
	if !ok {
		return nil, fmt.Errorf("secret ID %v not found", id)
	}
	return s, nil
}
