package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	StoragePath   string `env:"DB_URL" env-required:"true"`
	ServerAddress string `env:"SERVER_ADDRESS" env-default:":8080"`
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
