package config

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"
)

// Configuration errors.
var (
	ErrDatabaseConnStrRequired = errors.New("database connection string is required")
	ErrDatabaseTypeRequired    = errors.New("database type is required")
)

// Config holds the migration application configuration.
type Config struct {
	Database DatabaseConfig `mapstructure:"database"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Type             string `mapstructure:"type"`
	ConnectionString string `mapstructure:"connection_string"`
}

// Load loads the configuration from environment variables and config files.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("./cmd/migration")
	viper.AddConfigPath("./cmd/migration/config")
	viper.AddConfigPath(os.ExpandEnv("$HOME/.useful-cookery"))
	viper.AddConfigPath("/etc/useful-cookery")

	// Set default values
	setDefaults()

	// Read environment variables
	viper.SetEnvPrefix("UC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundErr viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFoundErr) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, using defaults and environment variables
		fmt.Fprintln(os.Stderr, "Config file not found, using defaults and environment variables")
	}

	// Parse configuration
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// setDefaults sets default values for configuration.
func setDefaults() {
	// Database defaults
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault(
		"database.connection_string",
		"postgres://postgres:postgres@localhost:5432/useful-cookery?sslmode=disable",
	)
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

	return nil
}
