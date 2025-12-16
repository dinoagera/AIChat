package service

import "context"

type AuthRepository interface {
	CreateUser(ctx context.Context, email, passHash string) error
}
