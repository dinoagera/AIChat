package domain

type User struct {
	UserID   int64
	Email    string
	PassHash []byte
}
