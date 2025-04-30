package config

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Default configuration values.
const (
	DefaultTokenExpiry = "24h"
)

// Configuration errors.
var (
	ErrDatabaseConnStrRequired = errors.New("database connection string is required")
	ErrDatabaseTypeRequired    = errors.New("database type is required")
	ErrJWTSecretRequired       = errors.New("JWT secret is required")
)

// Config holds the CLI application configuration.
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
	Auth     AuthConfig     `mapstructure:"auth"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Type             string `mapstructure:"type"`
	ConnectionString string `mapstructure:"connection_string"`
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	JWTSecret   string        `mapstructure:"jwt_secret"`
	TokenExpiry time.Duration `mapstructure:"token_expiry"`
}

// Load loads the configuration from the specified file or environment variables.
func Load(configFile string) (*Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Use default config paths
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("./cmd/cli/config")
		viper.AddConfigPath("$HOME/.useful-cookery")
		viper.AddConfigPath("/etc/useful-cookery")
	}

	// Set default values
	// Database defaults - using Postgres as default
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault(
		"database.connection_string",
		"postgres://postgres:postgres@localhost:5432/useful-cookery?sslmode=disable",
	)

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "change-me-in-production")
	viper.SetDefault("auth.token_expiry", DefaultTokenExpiry)

	// Read environment variables
	viper.SetEnvPrefix("UC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, using defaults and environment variables
		slog.Info("Config file not found, using defaults and environment variables")
	}

	// Parse configuration
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	// Log database connection info
	slog.Info("Using database", "type", cfg.Database.Type, "connection", cfg.Database.ConnectionString)

	return &cfg, nil
}

// validateConfig validates the configuration.
func validateConfig(config *Config) error {
	// Validate database config
	if config.Database.Type == "" {
		return ErrDatabaseTypeRequired
	}
	if config.Database.ConnectionString == "" {
		return ErrDatabaseConnStrRequired
	}

	// Validate Auth config
	if config.Auth.JWTSecret == "" {
		return ErrJWTSecretRequired
	}

	return nil
}
