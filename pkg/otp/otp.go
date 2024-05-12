package otp

import "github.com/xlzd/gotp"

type GOTPGenerator struct {
	RandomSecretLength int
}

func NewGOTPGenerator(secretLength int) *GOTPGenerator {
	return &GOTPGenerator{
		RandomSecretLength: secretLength,
	}
}

func (g *GOTPGenerator) RandomSecret() (string, error) {
	return gotp.RandomSecret(g.RandomSecretLength), nil
}
