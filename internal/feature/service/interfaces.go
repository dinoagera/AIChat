package service

import (
	"context"

	client "github.com/dinoagera/AIChat/internal/clients/http"
	"github.com/dinoagera/AIChat/internal/domain"
)

//go:generate
type AuthRepository interface {
	CreateUser(ctx context.Context, email, passHash string) error
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (domain.RefreshToken, error)
	ReplaceRefreshToken(ctx context.Context, oldToken, newToken string) error
	CreateSession(ctx context.Context, newToken string, userID int64) error
}
type BrigadeRepository interface {
	AddBrigade(ctx context.Context, req *domain.Brigade) error
	CheckName(ctx context.Context, name string) (bool, error)
	CheckBrigadeByID(ctx context.Context, id int64) (bool, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}
type LLMService interface {
	ParseEmergencyText(ctx context.Context, text string) (client.ParsedEmergency, error)
}
type GeocoderService interface {
	Geocode(ctx context.Context, address string) (lat, lon float64, err error)
}
type RoutingAPI interface {
	GetETA(ctx context.Context, fromLat, fromLon, toLat, toLon float64) (etaMin int, distanceKm float64, err error)
}
