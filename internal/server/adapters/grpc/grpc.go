package grpc

import (
	"errors"
	"github.com/itohin/gophkeeper/internal/server/usecases"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type Server struct {
	srv *grpc.Server
}

func NewServer(auth Auth) *Server {
	srv := grpc.NewServer()
	pb.RegisterAuthServer(srv, &AuthServer{
		auth: auth,
	})

	return &Server{srv: srv}
}

func getErrorCode(err error) codes.Code {
	var invalidArgument *usecases.InvalidArgumentError
	if errors.As(err, &invalidArgument) {
		return codes.InvalidArgument
	}

	return codes.Internal
}
