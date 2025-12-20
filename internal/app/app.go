package app

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dinoagera/AIChat/config"
	"github.com/dinoagera/AIChat/internal/http/handler"
	"github.com/dinoagera/AIChat/internal/repository/postgres"
	"github.com/dinoagera/AIChat/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func Run(cfg *config.Config, l *slog.Logger) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	pool, err := pgxpool.New(ctx, cfg.StoragePath)
	if err != nil {
		log.Fatalf("can't connect to postgresql: %v", err)
	}
	defer pool.Close()
	authRepository := postgres.NewAuthRepository(pool)
	authService := service.NewAuthService(l, authRepository)
	authHandler := handler.NewAuthHandler(l, authService)
	r := gin.New()
	authHandler.SetupRoutes(r)
	server := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: r,
	}
	go func() {
		l.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server failed: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	l.Info("shutting down server...")
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server shutdown failed: %v", err)
	}
	l.Info("server exited")
}
