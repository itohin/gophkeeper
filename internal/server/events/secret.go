package events

import "github.com/itohin/gophkeeper/internal/server/dto"

const (
	TypeCreated = iota + 1
	TypeUpdated
	TypeDeleted
)

type SecretEvent struct {
	EventType int
	Secret    dto.SecretDTO
}
