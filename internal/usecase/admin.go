package usecase

import (
	"context"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"go_store/config"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ AdminUseCase = (*adminUseCase)(nil)

type adminUseCase struct {
	logger *zap.Logger
	cfg    *config.Admin
}

func NewAdminUseCase(logger *zap.Logger, cfg *config.Admin) AdminUseCase {
	return &adminUseCase{logger: logger, cfg: cfg}
}

func (a *adminUseCase) Login(_ context.Context, username string, password string) (string, error) {
	if username != a.cfg.Username {
		a.logger.Warn("invalid username", zap.String("username", username))
		return "", status.Error(codes.Unauthenticated, "wrong credentials")
	}
	err := bcrypt.CompareHashAndPassword([]byte(a.cfg.PasswordHash), []byte(password))
	if err != nil {
		a.logger.Warn("wrong password", zap.String("password", password), zap.String("password hash", a.cfg.PasswordHash), zap.Error(err))
		return "", status.Error(codes.Unauthenticated, "wrong credentials")
	}

	token := jwt.New(jwt.SigningMethodHS256)

	tokenString, err := token.SignedString([]byte(a.cfg.JWTSecret))
	if err != nil {
		return "", status.Error(codes.Unauthenticated, "invalid token")
	}
	return tokenString, nil
}
