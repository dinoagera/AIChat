package tests

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/dinoagera/AIChat/internal/domain"
	"github.com/dinoagera/AIChat/internal/service"
	authMocks "github.com/dinoagera/AIChat/internal/service/mocks"
	tokenMocks "github.com/dinoagera/AIChat/pkg/auth/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestAuthService_SignIn(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		setupMocks func(*authMocks.AuthRepository, *tokenMocks.TokenManager)
		wantErr    error
		wantToken  bool
	}{
		{
			name:     "success",
			email:    "test@mail.ru",
			password: "pass123",
			setupMocks: func(ar *authMocks.AuthRepository, tm *tokenMocks.TokenManager) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
				ar.On("GetUserByEmail", mock.Anything, "test@mail.ru").Return(domain.User{UserID: 1, Email: "test@mail.ru", PassHash: hash}, nil)
				tm.On("NewJWT", int64(1)).Return("access_token", nil)
				tm.On("NewRefreshToken").Return("refresh_token", nil)
				ar.On("CreateSession", mock.Anything, "refresh_token", int64(1)).Return(nil)
			},
			wantErr:   nil,
			wantToken: true,
		},
		{
			name:     "user not found",
			email:    "test@mail.ru",
			password: "pass123",
			setupMocks: func(ar *authMocks.AuthRepository, tm *tokenMocks.TokenManager) {
				ar.On("GetUserByEmail", mock.Anything, "test@mail.ru").Return(domain.User{}, domain.ErrUserNotFound)
			},
			wantErr:   domain.ErrUserNotFound,
			wantToken: false,
		},
		{
			name:     "wrong password",
			email:    "test@mail.ru",
			password: "wrongpass",
			setupMocks: func(ar *authMocks.AuthRepository, tm *tokenMocks.TokenManager) {
				hash, _ := bcrypt.GenerateFromPassword([]byte("correctpass"), bcrypt.DefaultCost)
				ar.On("GetUserByEmail", mock.Anything, "test@mail.ru").Return(domain.User{UserID: int64(1), Email: "test@mail.ru", PassHash: hash}, nil)
			},
			wantErr:   domain.ErrPasswordWrong,
			wantToken: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := authMocks.NewAuthRepository(t)
			mockTokenMgr := tokenMocks.NewTokenManager(t)
			svc := service.NewAuthService(slog.Default(), mockRepo, mockTokenMgr)
			tt.setupMocks(mockRepo, mockTokenMgr)
			access, refresh, err := svc.SignIn(context.Background(), tt.email, tt.password)
			if tt.wantErr != nil {
				if errors.Is(tt.wantErr, domain.ErrPasswordWrong) {
					assert.ErrorIs(t, err, domain.ErrPasswordWrong)
				} else if errors.Is(tt.wantErr, domain.ErrUserNotFound) {
					assert.ErrorIs(t, err, domain.ErrUserNotFound)
				}
			} else {
				assert.NoError(t, err)
				if tt.wantToken {
					assert.NotEmpty(t, access)
					assert.NotEmpty(t, refresh)
				}
			}
		})
	}
}
func TestAuthService_SignUp(t *testing.T) {
	tests := []struct {
		name       string
		email      string
		password   string
		setupMocks func(*authMocks.AuthRepository)
		wantErr    error
	}{
		{
			name:     "success",
			email:    "test@mail.ru",
			password: "pass123",
			setupMocks: func(ar *authMocks.AuthRepository) {
				// hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
				ar.On("CreateUser", mock.Anything, "test@mail.ru", mock.Anything).Return(nil)
			},
			wantErr: nil,
		},
		{
			name:     "user is already exists",
			email:    "test@mail.ru",
			password: "pass123",
			setupMocks: func(ar *authMocks.AuthRepository) {
				// hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.DefaultCost)
				ar.On("CreateUser", mock.Anything, "test@mail.ru", mock.Anything).Return(domain.ErrUserAlreadyExists)
			},
			wantErr: domain.ErrUserAlreadyExists,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := authMocks.NewAuthRepository(t)
			mockTokenMgr := tokenMocks.NewTokenManager(t)
			srv := service.NewAuthService(slog.Default(), mockRepo, mockTokenMgr)
			tt.setupMocks(mockRepo)
			err := srv.SignUp(context.Background(), tt.email, tt.password)
			if tt.wantErr != nil {
				if errors.Is(err, domain.ErrUserAlreadyExists) {
					assert.ErrorIs(t, err, domain.ErrUserAlreadyExists)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
