package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Log      LogConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port string
	Host string
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret    string
	Algorithm string
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	viper.AutomaticEnv()

	// Set defaults
	viper.SetDefault("PORT", "8005")
	viper.SetDefault("HOST", "0.0.0.0")
	viper.SetDefault("JWT_ALGORITHM", "HS256")
	viper.SetDefault("LOG_LEVEL", "info")

	config := &Config{
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Host: viper.GetString("HOST"),
		},
		Database: DatabaseConfig{
			URL: viper.GetString("DATABASE_URL"),
		},
		JWT: JWTConfig{
			Secret:    viper.GetString("JWT_SECRET"),
			Algorithm: viper.GetString("JWT_ALGORITHM"),
		},
		Log: LogConfig{
			Level: viper.GetString("LOG_LEVEL"),
		},
	}

	// Validate required fields
	if config.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if config.JWT.Secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}

	return config, nil
}

// GetEnv gets environment variable with fallback
func GetEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
