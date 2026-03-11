package main

import (
	"context"
	"os"
	"os/signal"
	"settings-service-go/internal/application/usecase"
	"settings-service-go/internal/config"
	"settings-service-go/internal/infrastructure/http"
	"settings-service-go/internal/infrastructure/persistence"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	// Setup logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load configuration")
	}

	// Set log level
	level, err := zerolog.ParseLevel(cfg.Log.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	log.Info().Msg("Starting Settings Service")

	// Initialize database
	ctx := context.Background()
	db, err := persistence.NewPostgresPool(ctx, cfg.Database.URL)
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer db.Close()

	// Initialize repositories
	userSettingsRepo := persistence.NewUserSettingsRepository(db)
	systemSettingsRepo := persistence.NewSystemSettingsRepository(db)

	// Initialize use cases
	userSettingsUseCase := usecase.NewUserSettingsUseCase(userSettingsRepo)
	systemSettingsUseCase := usecase.NewSystemSettingsUseCase(systemSettingsRepo)

	// Initialize HTTP server
	server := http.NewServer(cfg, userSettingsUseCase, systemSettingsUseCase)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatal().Err(err).Msg("Failed to start server")
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Received shutdown signal")

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error().Err(err).Msg("Server shutdown error")
	}

	log.Info().Msg("Server stopped")
}
