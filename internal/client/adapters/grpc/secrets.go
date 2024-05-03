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
	secret, err := c.secretsHydrator.FromProto(s.Secret)
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
		secret, err := c.secretsHydrator.FromProto(v)
		if err != nil {
			return nil, handleError(err)
		}
		secrets[v.Id] = secret
	}
	return secrets, nil
}

func (c *Client) CreateSecret(ctx context.Context, secret *entities.Secret) error {
	ps, err := c.secretsHydrator.ToProto(secret)
	if err != nil {
		return fmt.Errorf("failed convert secret to proto: %v", err)
	}
	_, err = c.secrets.Create(ctx, &pb.CreateRequest{
		Secret: ps,
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) DeleteSecret(ctx context.Context, secret *entities.Secret) error {
	ps, err := c.secretsHydrator.ToProto(secret)
	if err != nil {
		return fmt.Errorf("failed convert secret to proto: %v", err)
	}
	_, err = c.secrets.Delete(ctx, &pb.DeleteRequest{
		Secret: ps,
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}
