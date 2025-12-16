package service

import (
	"context"
	"log/slog"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log            *slog.Logger
	authRepository AuthRepository
}

func NewAuthService(log *slog.Logger, authRepository AuthRepository) *AuthService {
	return &AuthService{
		log:            log,
		authRepository: authRepository,
	}
}
func (as *AuthService) SignUp(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		as.log.Info("failed to generate passhash", "err", err)
		return err
	}
	err = as.authRepository.CreateUser(ctx, email, string(passHash))
	return err
}
