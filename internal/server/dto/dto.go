package dto

type SecretDTO struct {
	ID         string
	Name       string
	SecretType uint32
	Notes      string
	Data       []byte
	UserID     string
}
