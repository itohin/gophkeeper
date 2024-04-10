package grpc

import (
	"context"
	"errors"
	"github.com/itohin/gophkeeper/internal/server/adapters/grpc/interceptors/jwt"
	errors2 "github.com/itohin/gophkeeper/pkg/errors"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
)

type JWTManager interface {
	GetClaims(tokenString string) (map[string]interface{}, error)
}

type Server struct {
	srv *grpc.Server
	log logger.Logger
}

func NewServer(auth Auth, secrets Secrets, log logger.Logger, jwtManager JWTManager) *Server {
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
		secrets: secrets,
		log:     log,
	})

	return &Server{
		srv: srv,
		log: log,
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", ":3200")
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
