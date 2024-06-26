package grpc

import (
	"context"
	"errors"
	"net"

	"github.com/itohin/gophkeeper/internal/server/adapters/grpc/interceptors/jwt"
	"github.com/itohin/gophkeeper/internal/server/config"
	"github.com/itohin/gophkeeper/internal/server/entities"
	errors2 "github.com/itohin/gophkeeper/pkg/errors"
	"github.com/itohin/gophkeeper/pkg/events"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type SecretHydrator interface {
	FromProto(in *pb.Secret, userID string) (*entities.Secret, error)
	ToProto(in *events.SecretDTO) (*pb.Secret, error)
}

type JWTManager interface {
	GetClaims(tokenString string) (map[string]interface{}, error)
}

type Server struct {
	srv *grpc.Server
	log logger.Logger
	cfg *config.AppConfig
}

func NewServer(
	auth Auth,
	secrets Secrets,
	log logger.Logger,
	jwtManager JWTManager,
	hydrator SecretHydrator,
	cfg *config.AppConfig,
) *Server {
	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			jwt.UnaryServerInterceptor(jwtManager.GetClaims),
		),
	)
	pb.RegisterAuthServer(srv, &AuthServer{
		auth: auth,
		log:  log,
	})
	pb.RegisterSecretsServer(srv, &SecretsServer{
		secrets:  secrets,
		hydrator: hydrator,
		log:      log,
	})

	return &Server{
		srv: srv,
		log: log,
		cfg: cfg,
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", s.cfg.GRPC.Address)
	if err != nil {
		s.log.Errorf("server error: %v", err)
		return err
	}
	s.log.Info("server started")
	if err := s.srv.Serve(listen); err != nil {
		return err
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) {
	s.srv.GracefulStop()
}

func getErrorCode(err error) codes.Code {
	var invalidArgument *errors2.InvalidArgumentError
	if errors.As(err, &invalidArgument) {
		return codes.InvalidArgument
	}

	return codes.Internal
}
