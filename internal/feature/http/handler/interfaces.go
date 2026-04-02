package handler

import "context"

type AuthService interface {
	SignUp(ctx context.Context, email, password string) error
	SignIn(ctx context.Context, email, password string) (string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
}
type BrigadeService interface {
	AddBrigade(ctx context.Context, name string, lat, lon float64, statis string) error
	UpdateStatus(ctx context.Context, id int64, status string) error
}
