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

func (c *Client) Verify(ctx context.Context, email, otp string) error {
	t, err := c.auth.Verify(ctx, &pb.VerifyRequest{
		Email:       email,
		Otp:         otp,
		Fingerprint: c.fingerPrint,
	})
	if err != nil {
		return handleError(err)
	}
	err = c.token.Update(t.Token.AccessToken, t.Token.RefreshToken)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) Login(ctx context.Context, email, password string) error {
	t, err := c.auth.Login(ctx, &pb.LoginRequest{
		Email:       email,
		Password:    password,
		Fingerprint: c.fingerPrint,
	})
	if err != nil {
		return handleError(err)
	}
	err = c.token.Update(t.Token.AccessToken, t.Token.RefreshToken)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) refresh(ctx context.Context) error {
	t, err := c.auth.Refresh(ctx, &pb.RefreshRequest{
		SessionId:   c.token.RefreshToken,
		Fingerprint: c.fingerPrint,
	})
	if err != nil {
		return handleError(err)
	}
	err = c.token.Update(t.Token.AccessToken, t.Token.RefreshToken)
	if err != nil {
		return handleError(err)
	}
	return nil
}

func (c *Client) Logout(ctx context.Context) error {
	_, err := c.auth.Logout(ctx, &pb.LogoutRequest{
		SessionId: c.token.RefreshToken,
	})
	if err != nil {
		return handleError(err)
	}
	close(c.shutdownCh)
	return nil
}
