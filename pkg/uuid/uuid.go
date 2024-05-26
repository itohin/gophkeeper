package uuid

import "github.com/google/uuid"

type GoogleUUIDGenerator struct{}

func NewGoogleUUIDGenerator() *GoogleUUIDGenerator {
	return &GoogleUUIDGenerator{}
}

func (g *GoogleUUIDGenerator) Generate() ([16]byte, error) {
	return uuid.NewV7()
}
