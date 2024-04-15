package grpc

import (
	"fmt"
	ji "github.com/itohin/gophkeeper/internal/client/adapters/grpc/interceptors/jwt"
	"github.com/itohin/gophkeeper/internal/client/adapters/storage"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/errors"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn            *grpc.ClientConn
	auth            pb.AuthClient
	secrets         pb.SecretsClient
	shutdownCh      chan struct{}
	token           *entities.Token
	fingerPrint     string
	secretsHydrator *storage.SecretsHydrator
}

func NewClient(
	fingerPrint string,
	token *entities.Token,
	shutdownCh chan struct{},
	secretsHydrator *storage.SecretsHydrator,
) (*Client, error) {
	conn, err := grpc.Dial(
		":3200",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			ji.UnaryClientInterceptor(token, fingerPrint),
		),
		grpc.WithChainStreamInterceptor(
			ji.StreamClientInterceptor(token, fingerPrint),
		),
	)
	if err != nil {
		return nil, err
	}
	auth := pb.NewAuthClient(conn)
	token.SetClient(auth)
	return &Client{
		conn:            conn,
		auth:            auth,
		secrets:         pb.NewSecretsClient(conn),
		shutdownCh:      shutdownCh,
		token:           token,
		fingerPrint:     fingerPrint,
		secretsHydrator: secretsHydrator,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func handleError(err error) error {
	e, ok := status.FromError(err)
	if ok && e.Code() == codes.InvalidArgument {
		return errors.NewDomainError(
			fmt.Errorf("input error: %v", e.Message()),
		)
	}
	return errors.NewDomainError(
		fmt.Errorf("internal error: please try again later"),
	)
}
