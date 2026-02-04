package auth

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/dinoagera/AIChat/internal/domain"
)

type TokenManager interface {
	NewJWT(userId int64) (string, error)
	Parse(accessToken string) (string, error)
	NewRefreshToken() (string, error)
	ParseRefreshToken(session domain.RefreshToken, oldRefreshToken string) error
}

var (
	ErrEmptySigningKey     = errors.New("empty signing key")
	ErrRefreshTTLUp        = errors.New("ttl refresh token is up")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
	ErrInvalidAccessToken  = errors.New("invalid token")
)

type Manager struct {
	signingKey string
	accessTTL  time.Duration
	refreshTTL time.Duration
}

func NewManager(signingKey string, accessTTL time.Duration, refreshTTL time.Duration) (*Manager, error) {
	if signingKey == "" {
		return nil, ErrEmptySigningKey
	}
	return &Manager{signingKey: signingKey, accessTTL: accessTTL, refreshTTL: refreshTTL}, nil
}

func (m *Manager) NewJWT(userId int64) (string, error) {
	claims := jwt.MapClaims{
		"uid": strconv.FormatInt(userId, 10),
		"exp": time.Now().Add(m.accessTTL).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.signingKey))
}

func (m *Manager) Parse(accessToken string) (int64, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(m.signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	if !token.Valid {
		return 0, ErrInvalidAccessToken
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, ErrInvalidAccessToken
	}
	uidStr, ok := claims["uid"].(string)
	if !ok || uidStr == "" {
		return 0, ErrInvalidAccessToken
	}
	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		return 0, ErrInvalidAccessToken
	}
	return uid, nil
}

func (m *Manager) NewRefreshToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", b), nil
}
func (m *Manager) ParseRefreshToken(session domain.RefreshToken, oldRefreshToken string) error {
	if time.Since(session.CreatedAt) > m.refreshTTL {
		return ErrRefreshTTLUp
	}
	if session.Token != oldRefreshToken {
		return ErrInvalidRefreshToken
	}
	return nil
}
