package config

import (
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
}

func InitConfig(log *slog.Logger) (*Config, error) {
	pathCfg := ".env"
	var cfg *Config
	err := cleanenv.ReadConfig(pathCfg, cfg)
	if err != nil {
		log.Debug("failed to init config", "err", err)
		return nil, err
	}
	return cfg, nil
}
