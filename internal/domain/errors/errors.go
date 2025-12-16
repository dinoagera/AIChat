package errors

const (
	ErrInvalidRequest    = "invalid request"
	ErrEmailRequired     = "email is required"
	ErrPasswordTooWeak   = "password must be at least 8 characters"
	ErrUserAlreadyExists = "user already exists"
	ErrInternal          = "internal server error"
	ErrUnauthorized      = "unauthorized"
)
