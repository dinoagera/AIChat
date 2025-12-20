package handler

import "context"

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (string, string, error)
}
