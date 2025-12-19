package postgres

import (
	"context"
	"errors"
	"time"

	domain "github.com/dinoagera/AIChat/internal/domain/errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	pool *pgxpool.Pool
}

func NewAuthRepository(pool *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{pool: pool}
}
func (ar *AuthRepository) CreateUser(ctx context.Context, email, passHash string) error {
	_, err := ar.pool.Exec(ctx, `INSERT INTO users (email, pass_hash, created_at) VALUES ($1, $2, $3)`, email, passHash, time.Now())
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return domain.ErrUserAlreadyExists
		}
		return err
	}
	return nil
}
