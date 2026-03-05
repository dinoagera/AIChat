package service

import (
	"context"
	"errors"
	"log/slog"

	domain "github.com/dinoagera/AIChat/internal/domain"
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
		if errors.Is(err, domain.ErrUserNotFound) {
			return "", "", domain.ErrUserNotFound
		}
		return "", "", err
	}
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(password))
	if err != nil {
		return "", "", domain.ErrPasswordWrong
	}
	accessToken, err := as.tokenManager.NewJWT(user.UserID)
	if err != nil {
		return "", "", domain.ErrPasswordWrong
	}
	refreshToken, err := as.tokenManager.NewRefreshToken()
	if err != nil {
		return "", "", domain.ErrPasswordWrong
	}
	err = as.authRepository.CreateSession(ctx, refreshToken, user.UserID)
	if err != nil {
		return "", "", domain.ErrService
	}
	return accessToken, refreshToken, nil
}
func (as *AuthService) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	session, err := as.authRepository.GetRefreshToken(ctx, refreshToken)
	if err != nil {
		as.log.Info("err", err)
		return "", "", err
	}
	err = as.tokenManager.ParseRefreshToken(session, refreshToken)
	if err != nil {
		as.log.Info("err", err)
		return "", "", err
	}
	accessToken, err := as.tokenManager.NewJWT(session.UserID)
	if err != nil {
		as.log.Info("wrong", "err", err)
		return "", "", err
	}
	newRefreshToken, err := as.tokenManager.NewRefreshToken()
	if err != nil {
		as.log.Info("wrong", "err", err)
		return "", "", err
	}
	err = as.authRepository.ReplaceRefreshToken(ctx, newRefreshToken, refreshToken)
	if err != nil {
		as.log.Info("err", err)
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}
