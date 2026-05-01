package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ajianaz/mikrosaas/internal/config"
	mikrosaas "github.com/ajianaz/mikrosaas/pkg"
)

func main() {
	// Load config from environment
	loader := config.NewLoader()

	cfg := mikrosaas.Config{
		Host:                 loader.String("SERVER_HOST", "0.0.0.0"),
		Port:                 loader.Int("SERVER_PORT", 8080),
		ReadTimeout:          loader.Duration("SERVER_READ_TIMEOUT", 0),
		WriteTimeout:         loader.Duration("SERVER_WRITE_TIMEOUT", 0),
		IdleTimeout:          loader.Duration("SERVER_IDLE_TIMEOUT", 0),
		DatabaseURL:          loader.MustString("DATABASE_URL"),
		MaxOpenConns:         loader.Int("DB_MAX_OPEN_CONNS", 0),
		MaxIdleConns:         loader.Int("DB_MAX_IDLE_CONNS", 0),
		ConnMaxLifetime:      loader.Duration("DB_CONN_MAX_LIFETIME", 0),
		JWTSecret:            loader.MustString("JWT_SECRET"),
		JWTAccessTokenTTL:    loader.Duration("JWT_ACCESS_TOKEN_TTL", 0),
		JWTRefreshTokenTTL:   loader.Duration("JWT_REFRESH_TOKEN_TTL", 0),
		BcryptCost:           loader.Int("BCRYPT_COST", 0),
		AppName:              loader.String("APP_NAME", "mikrosaas"),
		Env:                  loader.String("APP_ENV", "development"),
	}

	// Create app
	app, err := mikrosaas.NewApp(cfg)
	if err != nil {
		slog.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	// Graceful shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go func() {
		<-ctx.Done()
		slog.Info("received shutdown signal")
		if err := app.Shutdown(); err != nil {
			slog.Error("shutdown error", "error", err)
		}
	}()

	// Start server
	if err := app.Start(); err != nil {
		slog.Error("server error", "error", err)
		os.Exit(1)
	}

	fmt.Println("server stopped")
}
