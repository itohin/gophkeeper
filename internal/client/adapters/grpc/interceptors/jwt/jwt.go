package jwt

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"google.golang.org/grpc"
)

func UnaryClientInterceptor(token *entities.Token, fingerPrint string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		if token.IsExpired() {
			err := token.Refresh(ctx, fingerPrint)
			if err != nil {
				return err
			}
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
