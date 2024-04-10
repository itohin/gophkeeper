package grpc

import (
	"context"
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	pb "github.com/itohin/gophkeeper/proto"
)

func (c *Client) GetSecret(ctx context.Context, id string) (*entities.Secret, error) {
	s, err := c.secrets.Get(ctx, &pb.GetRequest{
		Id: id,
	})
	if err != nil {
		return nil, handleError(err)
	}
	secret, err := c.buildSecret(s.Secret)
	if err != nil {
		return nil, handleError(err)
	}
	return secret, nil
}

func (c *Client) SearchSecrets(ctx context.Context) (map[string]*entities.Secret, error) {
	s, err := c.secrets.Search(ctx, &pb.SearchRequest{})
	if err != nil {
		return nil, handleError(err)
	}
	secrets := make(map[string]*entities.Secret, len(s.Secrets))
	for _, v := range s.Secrets {
		secret, err := c.buildSecret(v)
		if err != nil {
			return nil, handleError(err)
		}
		secrets[v.Id] = secret
	}
	return secrets, nil
}

func (c *Client) buildSecret(v *pb.Secret) (*entities.Secret, error) {
	secret := &entities.Secret{
		ID:         v.Id,
		Name:       v.Name,
		SecretType: v.SecretType,
		Notes:      v.Notes,
	}
	switch d := v.Data.(type) {
	case *pb.Secret_Password:
		secret.Data = entities.Password{
			Login:    d.Password.Login,
			Password: d.Password.Password,
		}
	case *pb.Secret_Text:
		secret.Data = d.Text
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}
	return secret, nil
}

func (c *Client) CreateText(ctx context.Context, secret *entities.Secret, text string) error {
	_, err := c.secrets.Create(ctx, &pb.CreateRequest{
		Secret: &pb.Secret{
			Name:       secret.Name,
			SecretType: secret.SecretType,
			Notes:      secret.Notes,
			Data: &pb.Secret_Text{
				Text: text,
			},
		},
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error {
	_, err := c.secrets.Create(ctx, &pb.CreateRequest{
		Secret: &pb.Secret{
			Name:       secret.Name,
			SecretType: secret.SecretType,
			Notes:      secret.Notes,
			Data: &pb.Secret_Password{
				Password: &pb.Password{
					Login:    password.Login,
					Password: password.Password,
				},
			},
		},
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}
