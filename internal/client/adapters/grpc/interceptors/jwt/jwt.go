package jwt

import (
	"context"
	"github.com/itohin/gophkeeper/internal/client/entities"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func UnaryClientInterceptor(token *entities.Token, fingerPrint string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req interface{},
		reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption) error {

		if needToSkip(method) {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		if token.IsExpired() {
			err := token.Refresh(ctx, fingerPrint)
			if err != nil {
				return err
			}
		}

		authCtx := metadata.AppendToOutgoingContext(ctx, "Authorization", "Bearer "+token.AccessToken)
		return invoker(authCtx, method, req, reply, cc, opts...)
	}
}

var authRoutes = map[string]struct{}{
	"/gophkeeper.Auth/Refresh":  {},
	"/gophkeeper.Auth/Login":    {},
	"/gophkeeper.Auth/Register": {},
	"/gophkeeper.Auth/Verify":   {},
	"/gophkeeper.Auth/Logout":   {},
}

func needToSkip(method string) bool {
	_, ok := authRoutes[method]
	return ok
}
