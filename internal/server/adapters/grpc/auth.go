package grpc

import (
	"context"

	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Register(ctx context.Context, email, password string) error
	Verify(ctx context.Context, email, otp, fingerprint string) (*entities.Token, error)
	Login(ctx context.Context, email, password, fingerprint string) (*entities.Token, error)
	Refresh(ctx context.Context, sessionID, fingerprint string) (*entities.Token, error)
	Logout(ctx context.Context, sessionID string) error
}

type AuthServer struct {
	pb.UnimplementedAuthServer
	auth Auth
	log  logger.Logger
}

func (a *AuthServer) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	err := a.auth.Register(ctx, in.Email, in.Password)
	if err != nil {
		a.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.RegisterResponse{}, nil
}

func (a *AuthServer) Verify(ctx context.Context, in *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	token, err := a.auth.Verify(ctx, in.Email, in.Otp, in.Fingerprint)
	if err != nil {
		a.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.VerifyResponse{
		Token: &pb.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (a *AuthServer) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	token, err := a.auth.Login(ctx, in.Email, in.Password, in.Fingerprint)
	if err != nil {
		a.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.LoginResponse{
		Token: &pb.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (a *AuthServer) Refresh(ctx context.Context, in *pb.RefreshRequest) (*pb.RefreshResponse, error) {
	token, err := a.auth.Refresh(ctx, in.SessionId, in.Fingerprint)
	if err != nil {
		a.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.RefreshResponse{
		Token: &pb.Token{
			AccessToken:  token.AccessToken,
			RefreshToken: token.RefreshToken,
		},
	}, nil
}

func (a *AuthServer) Logout(ctx context.Context, in *pb.LogoutRequest) (*pb.LogoutResponse, error) {
	err := a.auth.Logout(ctx, in.SessionId)
	if err != nil {
		a.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.LogoutResponse{}, nil
}
