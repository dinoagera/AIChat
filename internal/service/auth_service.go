package service

import (
	"context"
	"log/slog"

	domain "github.com/dinoagera/AIChat/internal/domain/errors"
	"github.com/dinoagera/AIChat/pkg/auth"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	log            *slog.Logger
	authRepository AuthRepository
	tokenManager   auth.TokenManager
}

func NewAuthService(log *slog.Logger, authRepository AuthRepository, tokenManager auth.TokenManager) *AuthService {
	return &AuthService{
		log:            log,
		authRepository: authRepository,
		tokenManager:   tokenManager,
	}
}
func (as *AuthService) SignUp(ctx context.Context, email, password string) error {
	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		as.log.Info("failed to generate passhash", "err", err)
		return err
	}
	err = as.authRepository.CreateUser(ctx, email, string(passHash))
	if err != nil {
		as.log.Info("failed to create user", "err", err)
		return err
	}
	return nil
}
func (as *AuthService) SignIn(ctx context.Context, email, password string) (string, string, error) {
	user, err := as.authRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return "", "", domain.ErrPasswordWrong
	}
	//todo:add create jwt token and create session
	return "", "", nil
}
