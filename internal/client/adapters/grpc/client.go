package grpc

import (
	"fmt"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"github.com/itohin/gophkeeper/pkg/errors"
	pb "github.com/itohin/gophkeeper/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Client struct {
	conn        *grpc.ClientConn
	auth        pb.AuthClient
	shutdownCh  chan struct{}
	token       *entities.Token
	fingerPrint string
}

func NewClient(fingerPrint string, shutdownCh chan struct{}) (*Client, error) {
	conn, err := grpc.Dial(":3200", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:        conn,
		auth:        pb.NewAuthClient(conn),
		shutdownCh:  shutdownCh,
		fingerPrint: fingerPrint,
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
