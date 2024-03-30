package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Secrets interface {
	Save(ctx context.Context, secret *entities.Secret) (*entities.Secret, error)
	GetUserSecrets(ctx context.Context, userID string) ([]entities.SecretDTO, error)
}

type SecretsServer struct {
	pb.UnimplementedSecretsServer
	secrets Secrets
	log     logger.Logger
}

func (s *SecretsServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	userSecrets, err := s.secrets.GetUserSecrets(ctx, ctx.Value("user_id").(string))
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var secrets []*pb.Secret
	for _, v := range userSecrets {
		secret, err := s.buildSecret(&v)
		if err != nil {
			s.log.Error(err)
			return nil, status.Error(codes.Internal, err.Error())
		}

		secrets = append(secrets, secret)
	}

	return &pb.SearchResponse{
		Secrets: secrets,
	}, nil
}

func (s *SecretsServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	data, err := s.getData(in.Secret)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	savedSecret, err := s.secrets.Save(
		ctx,
		&entities.Secret{
			Name:       in.Secret.Name,
			Notes:      in.Secret.Notes,
			SecretType: in.Secret.SecretType,
			Data:       data,
			UserID:     ctx.Value("user_id").(string),
		},
	)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.CreateResponse{
		Id: savedSecret.ID,
	}, nil
}

func (s *SecretsServer) buildSecret(in *entities.SecretDTO) (*pb.Secret, error) {
	var t entities.Text
	var p entities.Password
	secret := pb.Secret{
		Id:         in.ID,
		Name:       in.Name,
		SecretType: in.SecretType,
		Notes:      in.Notes,
	}
	switch in.SecretType {
	case entities.TypeText:
		err := json.Unmarshal(in.Data, &t)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Text{Text: t.Text}
	case entities.TypePassword:
		err := json.Unmarshal(in.Data, &p)
		if err != nil {
			return nil, err
		}
		secret.Data = &pb.Secret_Password{
			Password: &pb.Password{Login: p.Login, Password: p.Password},
		}
	default:
		return nil, fmt.Errorf("unknown secret type")
	}
	return &secret, nil
}

func (s *SecretsServer) getData(in *pb.Secret) ([]byte, error) {
	switch d := in.Data.(type) {
	case *pb.Secret_Text:
		return json.Marshal(&entities.Text{
			Text: d.Text,
		})
	case *pb.Secret_Password:
		return json.Marshal(&entities.Password{
			Login:    d.Password.Login,
			Password: d.Password.Password,
		})
	default:
		return nil, fmt.Errorf("unknown secret data type")
	}
}
