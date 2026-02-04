package service

import (
	"context"

	"github.com/dinoagera/AIChat/internal/domain"
)

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passHash string) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
	PutRefreshToken(ctx context.Context, refreshToken string) error
}
