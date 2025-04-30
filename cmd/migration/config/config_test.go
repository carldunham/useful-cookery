//nolint:testpackage // Need to access unexported validateConfig function
package config

import (
	"errors"
	"testing"
)

// Helper function to create a valid config for testing.
func createValidConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Type:             "postgres",
			ConnectionString: "postgres://user:pass@localhost:5432/db",
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

//nolint:paralleltest // Using t.Setenv() which cannot be used with t.Parallel()
func TestValidateConfig(t *testing.T) {
	// Test valid config
	t.Run("valid config", func(t *testing.T) {
		cfg := createValidConfig()
		testConfigValidation(t, cfg, nil)
	})

	// Test missing database type
	t.Run("missing database type", func(t *testing.T) {
		cfg := createValidConfig()
		cfg.Database.Type = ""
		testConfigValidation(t, cfg, ErrDatabaseTypeRequired)
	})

	// Test missing database connection string
	t.Run("missing database connection string", func(t *testing.T) {
		cfg := createValidConfig()
		cfg.Database.ConnectionString = ""
		testConfigValidation(t, cfg, ErrDatabaseConnStrRequired)
	})
}

func TestLoad(t *testing.T) {
	// Set environment variables for testing
	t.Setenv("UC_DATABASE_TYPE", "postgres")
	t.Setenv("UC_DATABASE_CONNECTION_STRING", "postgres://test:test@localhost:5432/testdb")

	// Test loading configuration from environment variables
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Verify that the configuration was loaded correctly
	if cfg.Database.Type != "postgres" {
		t.Errorf("Expected Database.Type = postgres, got %s", cfg.Database.Type)
	}
	expectedConnStr := "postgres://test:test@localhost:5432/testdb"
	if cfg.Database.ConnectionString != expectedConnStr {
		t.Errorf("Expected Database.ConnectionString = %s, got %s",
			expectedConnStr, cfg.Database.ConnectionString)
	}
}

func TestSetDefaults(t *testing.T) {
	// Clear environment variables for this test
	t.Setenv("UC_DATABASE_TYPE", "")
	t.Setenv("UC_DATABASE_CONNECTION_STRING", "")

	// Load configuration with defaults
	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	// Test that defaults are set correctly
	if cfg.Database.Type != "postgres" {
		t.Errorf("Expected Database.Type = postgres, got %s", cfg.Database.Type)
	}
	expectedConnStr := "postgres://postgres:postgres@localhost:5432/useful-cookery?sslmode=disable"
	if cfg.Database.ConnectionString != expectedConnStr {
		t.Errorf("Expected Database.ConnectionString = %s, got %s",
			expectedConnStr, cfg.Database.ConnectionString)
	}
}
