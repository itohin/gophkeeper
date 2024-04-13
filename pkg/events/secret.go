package events

const (
	TypeCreated = iota + 1
	TypeUpdated
	TypeDeleted
)

type SecretEvent struct {
	EventType int
	Secret    SecretDTO
}

type SecretDTO struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       []byte
	UserID     string
}
