package entities

import "time"

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiration   time.Time
}
