package entities

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSession_IsExpired(t *testing.T) {
	tests := []struct {
		name      string
		ExpiresAt time.Time
		want      bool
	}{
		{
			name:      "expired",
			ExpiresAt: time.Now().Add(-1 * time.Second),
			want:      true,
		},
		{
			name:      "not expired",
			ExpiresAt: time.Now().Add(1 * time.Second),
			want:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Session{
				ExpiresAt: tt.ExpiresAt,
			}
			assert.Equalf(t, tt.want, s.IsExpired(), "IsExpired()")
		})
	}
}
