package config

import (
	"log/slog"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath    string        `env:"DB_URL" env-required:"true"`
	ServerAddress  string        `env:"SERVER_ADDRESS" env-default:":8080"`
	SecretKey      string        `env:"JWT_SECRET" env-required:"true"`
	TokenTTL       time.Duration `env:"TOKEN_TTL" env-default:"15m"`
	RefreshTTL     time.Duration `env:"REFRESH_TTL" env-default:"720h"`
	APIKey         string        `env:"API_KEY" env-required:"true"`
	BaseURLWeather string        `env:"BaseURLWeather"`
}

func InitConfig(log *slog.Logger) (*Config, error) {
	pathCfg := ".env"
	cfg := &Config{}
	err := cleanenv.ReadConfig(pathCfg, cfg)
	if err != nil {
		log.Debug("failed to init config", " err", err)
		return nil, err
	}
	return cfg, nil
}
