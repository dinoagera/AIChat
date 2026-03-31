package service

import (
	"context"

	"github.com/dinoagera/AIChat/internal/domain"
)

//go:generate
type AuthRepository interface {
	CreateUser(ctx context.Context, email, passHash string) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
	ReplaceRefreshToken(ctx context.Context, oldToken, newToken string) error
	CreateSession(ctx context.Context, newToken string, userID int64) error
}
type BrigadeRepository interface {
	AddBrigade(ctx context.Context, req *domain.Brigade) error
	CheckName(ctx context.Context, name string) bool
}
