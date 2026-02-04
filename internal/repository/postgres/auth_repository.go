package postgres

import (
	"context"
	"errors"
	"time"

	domain "github.com/dinoagera/AIChat/internal/domain"
	"github.com/jackc/pgx/v5"
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
func (ar *AuthRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var user domain.User
	err := ar.pool.QueryRow(ctx, `SELECT user_id, email, pass_hash FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}
func (ar *AuthRepository) GetRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {

}
func (ar *AuthRepository) PutRefreshToken(ctx context.Context, refreshToken string) error {

}
