package interceptor

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"strings"
)

func AuthInterceptor(secret string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if strings.HasPrefix(info.FullMethod, "/store.admin.AdminService/") &&
			info.FullMethod != "/store.admin.AdminService/Login" {

			md, ok := metadata.FromIncomingContext(ctx)
			if !ok {
				return nil, status.Error(codes.Unauthenticated, "missing metadata")
			}

			authHeader := md["authorization"]
			if len(authHeader) == 0 || !strings.HasPrefix(authHeader[0], "Bearer ") {
				return nil, status.Error(codes.Unauthenticated, "invalid or missing authorization header")
			}

			tokenString := strings.TrimPrefix(authHeader[0], "Bearer ")

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(secret), nil
			})

			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
			}

			if !token.Valid {
				return nil, status.Error(codes.Unauthenticated, "invalid token")
			}
		}

		return handler(ctx, req)
	}
}
