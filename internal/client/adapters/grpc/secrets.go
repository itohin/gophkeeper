package grpc

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
	pb "github.com/itohin/gophkeeper/proto"
)

func (c *Client) CreateText(ctx context.Context, secret *entities.Secret, text string) error {
	_, err := c.secrets.CreateText(ctx, &pb.CreateTextRequest{
		Name:  secret.Name,
		Text:  text,
		Notes: secret.Notes,
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) CreatePassword(ctx context.Context, secret *entities.Secret, password *entities.Password) error {
	_, err := c.secrets.CreatePassword(ctx, &pb.CreatePasswordRequest{
		Name: secret.Name,
		Password: &pb.Password{
			Login:    password.Login,
			Password: password.Password,
		},
		Notes: secret.Notes,
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}
