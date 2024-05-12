package grpc

import (
	"context"

	pb "github.com/itohin/gophkeeper/proto"
)

func (c *Client) Register(ctx context.Context, email, password string) error {
	_, err := c.auth.Register(ctx, &pb.RegisterRequest{
		Email:    email,
		Password: password,
	})
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) Verify(ctx context.Context, email, otp string) (string, error) {
	var userID string
	t, err := c.auth.Verify(ctx, &pb.VerifyRequest{
		Email:       email,
		Otp:         otp,
		Fingerprint: c.fingerPrint,
	})
	if err != nil {
		return userID, handleError(err)
	}
	err = c.token.Update(t.Token.AccessToken, t.Token.RefreshToken)
	if err != nil {
		return userID, handleError(err)
	}
	return c.token.UserID, nil
}

func (c *Client) Login(ctx context.Context, email, password string) (string, error) {
	var userID string
	t, err := c.auth.Login(ctx, &pb.LoginRequest{
		Email:       email,
		Password:    password,
		Fingerprint: c.fingerPrint,
	})
	if err != nil {
		return userID, handleError(err)
	}
	err = c.token.Update(t.Token.AccessToken, t.Token.RefreshToken)
	if err != nil {
		return userID, handleError(err)
	}
	return c.token.UserID, nil
}

func (c *Client) Logout(ctx context.Context) error {
	if c.token.AccessToken == "" {
		return nil
	}
	_, err := c.auth.Logout(ctx, &pb.LogoutRequest{
		SessionId: c.token.RefreshToken,
	})
	if err != nil {
		return handleError(err)
	}
	close(c.shutdownCh)
	return nil
}
