package grpc

import (
	"context"
	"encoding/json"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Secrets interface {
	Store(ctx context.Context, secret *entities.Secret) error
}

type SecretsServer struct {
	pb.UnimplementedSecretsServer
	secrets Secrets
	log     logger.Logger
}

func (s *SecretsServer) CreateText(ctx context.Context, in *pb.CreateTextRequest) (*pb.CreateTextResponse, error) {
	err := s.secrets.Store(
		ctx,
		&entities.Secret{
			Name:       in.Name,
			Notes:      in.Notes,
			SecretType: entities.TextType,
			Data:       in.Text,
		},
	)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.CreateTextResponse{}, nil
}

func (s *SecretsServer) CreatePassword(ctx context.Context, in *pb.CreatePasswordRequest) (*pb.CreatePasswordResponse, error) {
	data, err := json.Marshal(in.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	err = s.secrets.Store(
		ctx,
		&entities.Secret{
			Name:       in.Name,
			Notes:      in.Notes,
			SecretType: entities.PasswordType,
			Data:       string(data),
		},
	)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.CreatePasswordResponse{}, nil
}
