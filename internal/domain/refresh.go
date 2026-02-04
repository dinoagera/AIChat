package domain

import "time"

type RefreshToken struct {
	Token     string
	UserID    int64
	CreatedAt time.Time
}
