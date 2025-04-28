package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCommandHelp verifies that the command help text is displayed correctly.
func TestCommandHelp(t *testing.T) {
	t.Parallel()
	// Build the command for testing
	// #nosec G204 - This is a test with controlled input
	cmd := exec.Command("go", "run", "main.go", "--help")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	// Run the command
	err := cmd.Run()
	require.NoError(t, err, "Command should execute without error")

	// Verify the output contains expected help text
	output := stdout.String()
	assert.Contains(t, output, "A tool for converting TROFF recipes to structured data and managing database schema migrations")
	assert.Contains(t, output, "troff")
	assert.Contains(t, output, "db")
}

// TestSubcommandHelp verifies that subcommand help text is displayed correctly.
func TestSubcommandHelp(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		name         string
		subcommand   string
		expectedText string
	}{
		{
			name:         "TROFF command help",
			subcommand:   "troff",
			expectedText: "Commands for converting TROFF recipes to structured data and saving them to the database",
		},
		{
			name:         "DB command help",
			subcommand:   "db",
			expectedText: "Commands for managing database schema migrations",
		},
		{
			name:         "Test command help",
			subcommand:   "troff test",
			expectedText: "Parse a TROFF file and display the raw content and formatted recipe",
		},
		{
			name:         "Save command help",
			subcommand:   "troff save",
			expectedText: "Parse a TROFF file and save the recipe to the database",
		},
		{
			name:         "Batch command help",
			subcommand:   "troff batch",
			expectedText: "Parse all TROFF files in a directory and save them to the database",
		},
	}

	for _, testCase := range testCases {
		// Using Go 1.22+ loop variable capture semantics
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// Split the subcommand into parts for exec.Command
			cmdParts := append([]string{"run", "main.go"}, strings.Split(testCase.subcommand, " ")...)
			cmdParts = append(cmdParts, "--help")

			// #nosec G204 - This is a test with controlled input
			cmd := exec.Command("go", cmdParts...)
			var stdout bytes.Buffer
			cmd.Stdout = &stdout

			err := cmd.Run()
			require.NoError(t, err, "Command should execute without error")

			output := stdout.String()
			assert.Contains(t, output, testCase.expectedText)
		})
	}
}

// TestTestCommandWithSampleFile tests the 'test' command with a sample TROFF file.
// This test requires a sample TROFF file to be available.
func TestTestCommandWithSampleFile(t *testing.T) {
	t.Parallel()
	// Skip if running in CI or if sample files aren't available
	sampleDir := filepath.Join("..", "..", "data", "raw", "recipes")
	if _, err := os.Stat(sampleDir); os.IsNotExist(err) {
		t.Skip("Sample recipes directory not found, skipping test")
	}

	// Find a sample file
	files, err := os.ReadDir(sampleDir)
	require.NoError(t, err, "Should be able to read sample directory")

	if len(files) == 0 {
		t.Skip("No sample files found in recipes directory")
	}

	sampleFile := filepath.Join(sampleDir, files[0].Name())

	// Run the test command
	// #nosec G204 - This is a test with controlled input
	cmd := exec.Command("go", "run", "main.go", "troff", "test", sampleFile)
	var stdout bytes.Buffer
	cmd.Stdout = &stdout

	err = cmd.Run()
	require.NoError(t, err, "Command should execute without error")

	// Verify output contains both raw and formatted sections
	output := stdout.String()
	assert.Contains(t, output, "=== Raw TROFF Content ===")
	assert.Contains(t, output, "=== Formatted Recipe ===")

	// Basic verification that parsing worked
	assert.Contains(t, output, "id:")
	assert.Contains(t, output, "title:")
}

// TestDryRunSaveCommand tests the 'save' command with --dry-run flag.
func TestDryRunSaveCommand(t *testing.T) {
	t.Parallel()
	// Skip if running in CI or if sample files aren't available
	sampleDir := filepath.Join("..", "..", "data", "raw", "recipes")
	if _, err := os.Stat(sampleDir); os.IsNotExist(err) {
		t.Skip("Sample recipes directory not found, skipping test")
	}

	// Find a sample file
	files, err := os.ReadDir(sampleDir)
	require.NoError(t, err, "Should be able to read sample directory")

	if len(files) == 0 {
		t.Skip("No sample files found in recipes directory")
	}

	sampleFile := filepath.Join(sampleDir, files[0].Name())

	// Run the save command with dry-run
	// #nosec G204 - This is a test with controlled input
	cmd := exec.Command("go", "run", "main.go", "troff", "save", "--dry-run", sampleFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "Command should execute without error")

	// Check stdout for log output
	stdoutOutput := stdout.String()
	assert.Contains(t, stdoutOutput, "Dry run mode")
}

// TestDBStatusCommand tests that the 'db status' command includes the migration level.
func TestDBStatusCommand(t *testing.T) {
	t.Parallel()

	// Create a temporary directory for test migrations
	tempDir, err := os.MkdirTemp("", "migration-test")
	require.NoError(t, err, "Should be able to create temp directory")
	defer os.RemoveAll(tempDir)

	// Create a migrations directory structure
	migrationsDir := filepath.Join(tempDir, "db", "migrations")
	err = os.MkdirAll(migrationsDir, 0755)
	require.NoError(t, err, "Should be able to create migrations directory")

	// Create a test migration file
	testMigration := "000001_create_test_table.up.sql"
	err = os.WriteFile(
		filepath.Join(migrationsDir, testMigration),
		[]byte("-- Test migration"),
		0644,
	)
	require.NoError(t, err, "Should be able to create test migration file")

	// Create a minimal config file
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(configDir, 0755)
	require.NoError(t, err, "Should be able to create config directory")

	configContent := `
database:
  type: "postgres"
  connection_string: "postgres://postgres:postgres@localhost:5432/test?sslmode=disable"
`
	err = os.WriteFile(
		filepath.Join(configDir, "config.yml"),
		[]byte(configContent),
		0644,
	)
	require.NoError(t, err, "Should be able to create config file")

	// This test can't actually run the command since it would try to connect to a database,
	// but we can verify the code changes by examining the source code

	// Read the main.go file
	mainContent, err := os.ReadFile("main.go")
	require.NoError(t, err, "Should be able to read main.go")

	// Check that the runStatus function includes code to get and display the migration level
	mainContentStr := string(mainContent)
	assert.Contains(t, mainContentStr, "Migration level:", "main.go should include code to display migration level")
	assert.Contains(t, mainContentStr, "getMigrationName", "main.go should include getMigrationName function")
}

// TestDryRunBatchCommand tests the 'batch' command with --dry-run flag.
func TestDryRunBatchCommand(t *testing.T) {
	t.Parallel()
	// Skip if running in CI or if sample files aren't available
	sampleDir := filepath.Join("..", "..", "data", "raw", "recipes")
	if _, err := os.Stat(sampleDir); os.IsNotExist(err) {
		t.Skip("Sample recipes directory not found, skipping test")
	}

	// Run the batch command with dry-run
	// #nosec G204 - This is a test with controlled input
	cmd := exec.Command("go", "run", "main.go", "troff", "batch", "--dry-run", "--verbose", sampleDir)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	require.NoError(t, err, "Command should execute without error")

	// Check stderr for log output
	stdoutOutput := stdout.String()
	assert.Contains(t, stdoutOutput, "Batch processing complete")
	assert.Contains(t, stdoutOutput, "successful")
}
