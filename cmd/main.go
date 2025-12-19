package main

import (
	"log"

	"github.com/dinoagera/AIChat/config"
	"github.com/dinoagera/AIChat/internal/app"
	"github.com/dinoagera/AIChat/pkg/logger"
)

func main() {
	logger := logger.InitLogger()
	cfg, err := config.InitConfig(logger)
	if err != nil {
		log.Fatal("failed to init config", " err", err)
	}
	app.Run(cfg, logger)
}
