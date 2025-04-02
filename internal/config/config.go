package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/openai/openai-go"
	"github.com/spf13/viper"
)

// Config holds the application configuration.
type Config struct {
	Server  ServerConfig  `mapstructure:"server"`
	DGraph  DGraphConfig  `mapstructure:"dgraph"`
	Redis   RedisConfig   `mapstructure:"redis"`
	Auth    AuthConfig    `mapstructure:"auth"`
	AI      AIConfig      `mapstructure:"ai"`
	Storage StorageConfig `mapstructure:"storage"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Port                string `mapstructure:"port"`
	CorsOrigin          string `mapstructure:"cors_origin"`
	EnablePlayground    bool   `mapstructure:"enable_playground"`
	ReadTimeoutSeconds  int    `mapstructure:"read_timeout_seconds"`
	WriteTimeoutSeconds int    `mapstructure:"write_timeout_seconds"`
	IdleTimeoutSeconds  int    `mapstructure:"idle_timeout_seconds"`
}

// DGraphConfig holds DGraph-related configuration.
type DGraphConfig struct {
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
	OpenAIKey       string                `mapstructure:"openai_key"`
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
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
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
	// Server defaults
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.cors_origin", "*")
	viper.SetDefault("server.enable_playground", true)
	viper.SetDefault("server.read_timeout_seconds", 30)
	viper.SetDefault("server.write_timeout_seconds", 30)
	viper.SetDefault("server.idle_timeout_seconds", 60)

	// DGraph defaults
	viper.SetDefault("dgraph.connection_string", "localhost:9080")

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
	viper.SetDefault("ai.cache_ttl_minutes", 1440) // 24 hours
	viper.SetDefault("ai.max_tokens", 4000)

	// Storage defaults
	viper.SetDefault("storage.type", "local")
	viper.SetDefault("storage.local_path", "./uploads")
}

// validateConfig validates the configuration.
func validateConfig(config *Config) error {
	// Validate server config
	if config.Server.Port == "" {
		return fmt.Errorf("server port is required")
	}

	// Validate DGraph config
	if config.DGraph.ConnectionString == "" {
		return fmt.Errorf("DGraph connection string is required")
	}

	// Validate Auth config
	if config.Auth.JWTSecret == "" {
		return fmt.Errorf("JWT secret is required")
	}

	// Validate AI config
	if config.AI.OpenAIKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}

	return nil
}
