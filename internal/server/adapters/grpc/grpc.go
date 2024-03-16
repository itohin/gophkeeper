package grpc

import (
	"context"
	"errors"
	errors2 "github.com/itohin/gophkeeper/pkg/errors"
	"github.com/itohin/gophkeeper/pkg/logger"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
)

type Server struct {
	srv *grpc.Server
	log logger.Logger
}

func NewServer(auth Auth, log logger.Logger) *Server {
	srv := grpc.NewServer()
	pb.RegisterAuthServer(srv, &AuthServer{
		auth: auth,
		log:  log,
	})

	return &Server{
		srv: srv,
		log: log,
	}
}

func (s *Server) Start() error {
	listen, err := net.Listen("tcp", ":3200")
	if err != nil {
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
