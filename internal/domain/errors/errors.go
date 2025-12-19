package domain

import "errors"

var (
	ErrInvalidRequest    = errors.New("invalid request")
	ErrEmailRequired     = errors.New("email is required")
	ErrPasswordTooWeak   = errors.New("password must be at least 8 characters")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInternal          = errors.New("internal server error")
	ErrUnauthorized      = errors.New("unauthorized")
)
