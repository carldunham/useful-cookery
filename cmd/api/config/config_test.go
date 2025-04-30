//nolint:testpackage // Need to access unexported validateConfig function
package config

import (
	"errors"
	"testing"
	"time"
)

// Helper function to create a valid config for testing.
func createValidConfig() *Config {
	return &Config{
		Server: ServerConfig{
			Port:                "8080",
			CORSOrigin:          "*",
			EnablePlayground:    true,
			ReadTimeoutSeconds:  30,
			WriteTimeoutSeconds: 30,
			IdleTimeoutSeconds:  60,
		},
		Database: DatabaseConfig{
			Type:             "postgres",
			ConnectionString: "postgres://user:pass@localhost:5432/db",
		},
		Redis: RedisConfig{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		},
		Auth: AuthConfig{
			JWTSecret:   "secret",
			TokenExpiry: 24 * time.Hour,
		},
		AI: AIConfig{
			OpenAIKey:       "test-key",
			EmbeddingModel:  "text-embedding-ada-002",
			CompletionModel: "gpt-3.5-turbo",
			EnableCache:     true,
			CacheTTLMinutes: 1440,
			MaxTokens:       4000,
		},
		Storage: StorageConfig{
			Type:      "local",
			LocalPath: "./uploads",
		},
	}
}

// Helper function to test config validation.
func testConfigValidation(t *testing.T, cfg *Config, wantErr error) {
	t.Helper()

	// Call validateConfig directly
	err := validateConfig(cfg)
	if !errors.Is(err, wantErr) {
		t.Errorf("validateConfig() error = %v, wantErr %v", err, wantErr)
	}
}

func TestValidateConfig(t *testing.T) {
	t.Parallel()

	// Test valid config
	t.Run("valid config", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		testConfigValidation(t, cfg, nil)
	})

	// Test missing server port
	t.Run("missing server port", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		cfg.Server.Port = ""
		testConfigValidation(t, cfg, ErrServerPortRequired)
	})

	// Test missing database type
	t.Run("missing database type", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		cfg.Database.Type = ""
		testConfigValidation(t, cfg, ErrDatabaseTypeRequired)
	})

	// Test missing database connection string
	t.Run("missing database connection string", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		cfg.Database.ConnectionString = ""
		testConfigValidation(t, cfg, ErrDatabaseConnStrRequired)
	})

	// Test missing JWT secret
	t.Run("missing JWT secret", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		cfg.Auth.JWTSecret = ""
		testConfigValidation(t, cfg, ErrJWTSecretRequired)
	})

	// Test missing OpenAI API key
	t.Run("missing OpenAI API key", func(t *testing.T) {
		t.Parallel()
		cfg := createValidConfig()
		cfg.AI.OpenAIKey = ""
		testConfigValidation(t, cfg, ErrOpenAIAPIKeyRequired)
	})
}

// Helper function to verify server config defaults.
func verifyServerDefaults(t *testing.T, cfg *Config) {
	t.Helper()

	// Test that environment variables override defaults
	if cfg.Server.Port != "9090" {
		t.Errorf("Expected Server.Port = 9090, got %s", cfg.Server.Port)
	}

	// Test that defaults are set for values not in environment
	expectedCORSOrigin := "http://localhost:3000"
	if cfg.Server.CORSOrigin != expectedCORSOrigin {
		t.Errorf("Expected Server.CORSOrigin = %s, got %s", expectedCORSOrigin, cfg.Server.CORSOrigin)
	}
	if !cfg.Server.EnablePlayground {
		t.Errorf("Expected Server.EnablePlayground = true, got %v", cfg.Server.EnablePlayground)
	}
	if cfg.Server.ReadTimeoutSeconds != DefaultReadTimeoutSeconds {
		t.Errorf("Expected Server.ReadTimeoutSeconds = %d, got %d",
			DefaultReadTimeoutSeconds, cfg.Server.ReadTimeoutSeconds)
	}
	if cfg.Server.WriteTimeoutSeconds != DefaultWriteTimeoutSeconds {
		t.Errorf("Expected Server.WriteTimeoutSeconds = %d, got %d",
			DefaultWriteTimeoutSeconds, cfg.Server.WriteTimeoutSeconds)
	}
	if cfg.Server.IdleTimeoutSeconds != DefaultIdleTimeoutSeconds {
		t.Errorf("Expected Server.IdleTimeoutSeconds = %d, got %d",
			DefaultIdleTimeoutSeconds, cfg.Server.IdleTimeoutSeconds)
	}
}

// Helper function to verify AI config defaults.
func verifyAIDefaults(t *testing.T, cfg *Config) {
	t.Helper()

	expectedEmbeddingModel := "text-embedding-ada-002"
	if cfg.AI.EmbeddingModel != expectedEmbeddingModel {
		t.Errorf("Expected AI.EmbeddingModel = %s, got %s", expectedEmbeddingModel, cfg.AI.EmbeddingModel)
	}
	if cfg.AI.CompletionModel != "gpt-3.5-turbo" {
		t.Errorf("Expected AI.CompletionModel = gpt-3.5-turbo, got %s", cfg.AI.CompletionModel)
	}
	if !cfg.AI.EnableCache {
		t.Errorf("Expected AI.EnableCache = true, got %v", cfg.AI.EnableCache)
	}
	if cfg.AI.CacheTTLMinutes != DefaultCacheTTLMinutes {
		t.Errorf("Expected AI.CacheTTLMinutes = %d, got %d",
			DefaultCacheTTLMinutes, cfg.AI.CacheTTLMinutes)
	}
	if cfg.AI.MaxTokens != DefaultMaxTokens {
		t.Errorf("Expected AI.MaxTokens = %d, got %d",
			DefaultMaxTokens, cfg.AI.MaxTokens)
	}
}

// Helper function to verify storage config defaults.
func verifyStorageDefaults(t *testing.T, cfg *Config) {
	t.Helper()

	if cfg.Storage.Type != "local" {
		t.Errorf("Expected Storage.Type = local, got %s", cfg.Storage.Type)
	}
	if cfg.Storage.LocalPath != "./uploads" {
		t.Errorf("Expected Storage.LocalPath = ./uploads, got %s", cfg.Storage.LocalPath)
	}
}

func TestSetDefaults(t *testing.T) {
	// Set up environment for testing
	t.Setenv("UC_SERVER_PORT", "9090")
	t.Setenv("UC_DATABASE_TYPE", "postgres")
	t.Setenv("UC_DATABASE_CONNECTION_STRING", "postgres://test:test@localhost:5432/testdb")
	t.Setenv("UC_AUTH_JWT_SECRET", "test-secret")
	t.Setenv("UC_AI_OPENAI_KEY", "test-openai-key")

	// Load configuration
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify defaults using helper functions
	verifyServerDefaults(t, cfg)
	verifyAIDefaults(t, cfg)
	verifyStorageDefaults(t, cfg)
}
