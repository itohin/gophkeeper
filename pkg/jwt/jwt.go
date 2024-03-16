package jwt

import (
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type JWTGOManager struct {
	signature  string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewJWTGOManager(signature string, accessTTL, refreshTTL time.Duration) (*JWTGOManager, error) {
	if signature == "" {
		return nil, errors.New("jwt manager error: empty signature")
	}
	return &JWTGOManager{
		signature:  signature,
		accessTTL:  accessTTL,
		refreshTTL: refreshTTL,
	}, nil
}

func (j *JWTGOManager) MakeJWT(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(j.accessTTL).Unix(),
		Subject:   userID,
	})

	return token.SignedString([]byte(j.signature))
}

func (j *JWTGOManager) MakeRefreshExpiration() time.Time {
	return time.Now().Add(j.refreshTTL)
}

func (j *JWTGOManager) GetClaims(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected jwt signing method: %v", token.Header["alg"])
		}
		return []byte(j.signature), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("invalid jwt token")
	}
	return claims, nil
}
