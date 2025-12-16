package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewTeamRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}
func (ar *AuthRepository) CreateUser(ctx context.Context, email, passHash string) error {
	return err
}
