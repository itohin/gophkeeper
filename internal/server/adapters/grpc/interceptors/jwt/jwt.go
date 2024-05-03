package jwt

import (
	"context"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func UnaryServerInterceptor(f func(tokenString string) (map[string]interface{}, error)) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if needToSkip(info.FullMethod) {
			return handler(ctx, req)
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, status.Error(codes.PermissionDenied, "authorization denied")
		}
		jwtString, err := getJWTString(md)
		if err != nil {
			return nil, status.Error(codes.PermissionDenied, "authorization denied")
		}

		claims, err := f(jwtString)

		if int64(claims["exp"].(float64)) < time.Now().Unix() {
			return nil, status.Error(codes.PermissionDenied, "authorization denied")
		}

		authCtx := context.WithValue(ctx, "user_id", claims["sub"])

		return handler(authCtx, req)
	}
}

func getJWTString(md metadata.MD) (string, error) {
	values := md.Get("Authorization")
	if len(values[0]) < 1 {
		return "", fmt.Errorf("invalid auth header")
	}
	headerParts := strings.Split(values[0], " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", fmt.Errorf("invalid auth header")
	}
	return headerParts[1], nil
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
