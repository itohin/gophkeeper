package entities

import "time"

type Token struct {
	accessToken  string
	refreshToken string
	expiration   time.Time
}
