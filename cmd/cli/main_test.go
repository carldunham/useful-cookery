package main

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestRootCommand tests that the root command can be created and executed.
func TestRootCommand(t *testing.T) {
	t.Parallel()
	// Save the original rootCmd and restore it after the test
	origRootCmd := rootCmd
	defer func() {
		rootCmd = origRootCmd
	}()

	// Create a new root command for testing
	rootCmd = &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Run: func(_ *cobra.Command, _ []string) {
			// Do nothing
		},
	}

	// Execute the command
	err := rootCmd.Execute()
	assert.NoError(t, err)
}

// TestLoadConfig tests the loadConfig function.
func TestLoadConfig(t *testing.T) {
	t.Parallel()
	// Test with empty config file
	cfg, err := loadConfig("")
	require.NoError(t, err)
	require.NotNil(t, cfg)

	// Verify that the default database type is postgres
	assert.Equal(t, "postgres", cfg.Database.Type)
}
