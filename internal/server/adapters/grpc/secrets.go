package grpc

import (
	"context"

	"github.com/itohin/gophkeeper/internal/server/entities"
	"github.com/itohin/gophkeeper/pkg/events"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Secrets interface {
	Save(ctx context.Context, secret *entities.Secret) (*events.SecretDTO, error)
	GetUserSecrets(ctx context.Context, userID string) ([]events.SecretDTO, error)
	GetUserSecret(ctx context.Context, userID, secretID string) (events.SecretDTO, error)
	DeleteUserSecret(ctx context.Context, secret *entities.Secret) (*events.SecretDTO, error)
}

type SecretsServer struct {
	pb.UnimplementedSecretsServer
	secrets  Secrets
	hydrator SecretHydrator
	log      logger.Logger
}

func (s *SecretsServer) Search(ctx context.Context, in *pb.SearchRequest) (*pb.SearchResponse, error) {
	userSecrets, err := s.secrets.GetUserSecrets(ctx, ctx.Value("user_id").(string))
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	var secrets []*pb.Secret
	for _, v := range userSecrets {
		secret, err := s.hydrator.ToProto(&v)
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

func (s *SecretsServer) Get(ctx context.Context, in *pb.GetRequest) (*pb.GetResponse, error) {
	sDTO, err := s.secrets.GetUserSecret(ctx, ctx.Value("user_id").(string), in.Id)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	secret, err := s.hydrator.ToProto(&sDTO)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetResponse{
		Secret: secret,
	}, nil
}

func (s *SecretsServer) Create(ctx context.Context, in *pb.CreateRequest) (*pb.CreateResponse, error) {
	secret, err := s.hydrator.FromProto(in.Secret, ctx.Value("user_id").(string))
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	savedSecret, err := s.secrets.Save(ctx, secret)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.CreateResponse{
		Id: savedSecret.ID,
	}, nil
}

func (s *SecretsServer) Delete(ctx context.Context, in *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	secret, err := s.hydrator.FromProto(in.Secret, ctx.Value("user_id").(string))
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	_, err = s.secrets.DeleteUserSecret(ctx, secret)
	if err != nil {
		s.log.Error(err)
		return nil, status.Error(getErrorCode(err), err.Error())
	}
	return &pb.DeleteResponse{}, err
}
