// Package config provides configuration functionality for the application.
// DEPRECATED: This package is deprecated and will be removed in a future version.
// Use the configuration packages in cmd/api/config, cmd/cli/config, and cmd/migration/config instead.
package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/spf13/viper"
)

// Default configuration values.
const (
	DefaultReadTimeoutSeconds  = 30
	DefaultWriteTimeoutSeconds = 30
	DefaultIdleTimeoutSeconds  = 60
	DefaultCacheTTLMinutes     = 1440 // 24 hours.
	DefaultMaxTokens           = 4000
)

// Configuration errors.
var (
	ErrServerPortRequired      = errors.New("server port is required")
	ErrDatabaseConnStrRequired = errors.New("database connection string is required")
	ErrDatabaseTypeRequired    = errors.New("database type is required")
	ErrJWTSecretRequired       = errors.New("JWT secret is required")
	ErrOpenAIAPIKeyRequired    = errors.New("OpenAI API key is required")
)

// Config holds the application configuration.
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	AI       AIConfig       `mapstructure:"ai"`
	Storage  StorageConfig  `mapstructure:"storage"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port                string `mapstructure:"port"`
	CORSOrigin          string `mapstructure:"cors_origin"`
	EnablePlayground    bool   `mapstructure:"enable_playground"`
	ReadTimeoutSeconds  int    `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `mapstructure:"write_timeout_seconds"`
	IdleTimeoutSeconds  int    `mapstructure:"idle_timeout_seconds"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Type             string `mapstructure:"type"`
	ConnectionString string `mapstructure:"connection_string"`
}

// RedisConfig holds Redis-related configuration.
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
}

// AuthConfig holds authentication-related configuration.
type AuthConfig struct {
	JWTSecret   string        `mapstructure:"jwt_secret"`
	TokenExpiry time.Duration `mapstructure:"token_expiry"`
}

// AIConfig holds AI-related configuration.
type AIConfig struct {
	OpenAIAPIKey    string                `mapstructure:"openai_api_key"`
	EmbeddingModel  openai.EmbeddingModel `mapstructure:"embedding_model"`
	CompletionModel string                `mapstructure:"completion_model"`
	EnableCache     bool                  `mapstructure:"enable_cache"`
	CacheTTLMinutes int                   `mapstructure:"cache_ttl_minutes"`
	MaxTokens       int                   `mapstructure:"max_tokens"`
}

// StorageConfig holds storage-related configuration.
type StorageConfig struct {
	Type      string `mapstructure:"type"`
	S3Bucket  string `mapstructure:"s3_bucket"`
	S3Region  string `mapstructure:"s3_region"`
	LocalPath string `mapstructure:"local_path"`
}

// Load loads the configuration from environment variables and config files.
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath(filepath.Join("$HOME", ".useful-cookery"))
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

	// Process environment variable placeholders in strings
	processEnvPlaceholders(&config)

	// Validate configuration
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// processEnvPlaceholders replaces ${ENV_VAR} placeholders with actual environment variable values.
func processEnvPlaceholders(config *Config) {
	// Process database connection string
	config.Database.ConnectionString = replaceEnvVars(config.Database.ConnectionString)

	// Process JWT secret
	config.Auth.JWTSecret = replaceEnvVars(config.Auth.JWTSecret)

	// Process OpenAI key
	config.AI.OpenAIAPIKey = replaceEnvVars(config.AI.OpenAIAPIKey)
}

// replaceEnvVars replaces ${ENV_VAR} placeholders with actual environment variable values.
func replaceEnvVars(input string) string {
	re := regexp.MustCompile(`\${([^}]+)}`)
	return re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract environment variable name (remove ${ and })
		envVar := match[2 : len(match)-1]

		// Get environment variable value
		value := os.Getenv(envVar)
		if value == "" {
			fmt.Fprintf(os.Stderr, "Warning: Environment variable %s not set\n", envVar)
			return match // Return original placeholder if env var not set
		}

		return value
	})
}

// setDefaults sets default values for configuration.
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.cors_origin", "*")
	viper.SetDefault("server.enable_playground", true)
	viper.SetDefault("server.read_timeout_seconds", DefaultReadTimeoutSeconds)
	viper.SetDefault("server.write_timeout_seconds", DefaultWriteTimeoutSeconds)
	viper.SetDefault("server.idle_timeout_seconds", DefaultIdleTimeoutSeconds)

	// Database defaults
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault(
		"database.connection_string",
		"postgres://postgres:postgres@localhost:5432/useful-cookery?sslmode=disable",
	)

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.password", "")
	viper.SetDefault("redis.db", 0)

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "change-me-in-production")
	viper.SetDefault("auth.token_expiry", "24h")

	// AI defaults
	viper.SetDefault("ai.embedding_model", "text-embedding-ada-002")
	viper.SetDefault("ai.completion_model", "gpt-3.5-turbo")
	viper.SetDefault("ai.enable_cache", true)
	viper.SetDefault("ai.cache_ttl_minutes", DefaultCacheTTLMinutes) // 24 hours
	viper.SetDefault("ai.max_tokens", DefaultMaxTokens)

	// Storage defaults
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local_path", "./uploads")
}

// validateConfig validates the configuration.
func validateConfig(config *Config) error {
	// Validate server config
	if config.Server.Port == "" {
		return ErrServerPortRequired
	}

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

	// Validate AI config
	if config.AI.OpenAIAPIKey == "" {
		return ErrOpenAIAPIKeyRequired
	}

	return nil
}
