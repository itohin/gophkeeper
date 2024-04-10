package entities

import (
	"context"
	"github.com/itohin/gophkeeper/proto"
	pb "github.com/itohin/gophkeeper/proto"
	"time"
)

type JWTManager interface {
	GetClaims(tokenString string) (map[string]interface{}, error)
}

type Token struct {
	AccessToken  string
	RefreshToken string
	Expiration   int64
	UserID       string
	client       proto.AuthClient
	jwt          JWTManager
}

func NewToken(jwt JWTManager) *Token {
	return &Token{jwt: jwt}
}

func (t *Token) SetClient(client proto.AuthClient) {
	t.client = client
}

func (t *Token) IsExpired() bool {
	return t.Expiration < time.Now().Unix()
}

func (t *Token) Refresh(ctx context.Context, fingerPrint string) error {
	r, err := t.client.Refresh(ctx, &pb.RefreshRequest{
		SessionId:   t.RefreshToken,
		Fingerprint: fingerPrint,
	})
	if err != nil {
		return err
	}

	return t.Update(r.Token.AccessToken, r.Token.RefreshToken)
}

func (t *Token) Update(at, rt string) error {
	claims, err := t.jwt.GetClaims(at)
	if err != nil {
		return err
	}
	t.AccessToken = at
	t.RefreshToken = rt
	t.Expiration = int64(claims["exp"].(float64))
	t.UserID = claims["sub"].(string)

	return nil
}
