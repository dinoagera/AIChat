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
	pool       *pgxpool.Pool
	ttlRefresh time.Duration
}

func NewAuthRepository(pool *pgxpool.Pool, ttlRefresh time.Duration) *AuthRepository {
	return &AuthRepository{pool: pool, ttlRefresh: ttlRefresh}
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
	err := ar.pool.QueryRow(ctx, `SELECT user_id, email, pass_hash FROM users WHERE email = $1`, email).Scan(&user.UserID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, domain.ErrUserNotFound
		}
		return domain.User{}, err
	}
	return user, nil
}
func (ar *AuthRepository) GetRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error) {
	var session domain.RefreshToken
	err := ar.pool.QueryRow(ctx, `SELECT token, user_id, expires_at, created_at FROM refresh_tokens WHERE token = $1`, refreshToken).Scan(&session.Token, &session.UserID, &session.ExpiresAt, &session.CreatedAt)
	if err != nil {
		return domain.RefreshToken{}, err
	}
	return session, nil
}
func (ar *AuthRepository) ReplaceRefreshToken(ctx context.Context, newToken, oldToken string) error {
	tx, err := ar.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	var (
		userID    string
		expiresAt time.Time
	)
	err = tx.QueryRow(ctx, `SELECT user_id, expires_at FROM refresh_tokens WHERE token = $1 AND expires_at > NOW()`, oldToken).Scan(&userID, &expiresAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		return err
	}
	_, err = tx.Exec(ctx, `DELETE FROM refresh_tokens WHERE token = $1`, oldToken)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO refresh_tokens (token, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`, newToken, userID, expiresAt, time.Now())
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
func (ar *AuthRepository) CreateSession(ctx context.Context, refreshToken string, userID int64) error {
	tx, err := ar.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)
	_, err = tx.Exec(ctx, `DELETE FROM refresh_tokens WHERE user_id = $1`, userID)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, `INSERT INTO refresh_tokens (token, user_id, expires_at, created_at) VALUES ($1, $2, $3, $4)`, refreshToken, userID, time.Now().Add(ar.ttlRefresh), time.Now())
	if err != nil {
		return err
	}
	return tx.Commit(ctx)
}
