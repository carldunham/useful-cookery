package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
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
	assert.Contains(t, output, "A tool for converting TROFF recipes to structured data and saving them to the database")
	assert.Contains(t, output, "test")
	assert.Contains(t, output, "save")
	assert.Contains(t, output, "batch")
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
			name:         "Test command help",
			subcommand:   "test",
			expectedText: "Parse a TROFF file and display the raw content and formatted recipe",
		},
		{
			name:         "Save command help",
			subcommand:   "save",
			expectedText: "Parse a TROFF file and save the recipe to the database",
		},
		{
			name:         "Batch command help",
			subcommand:   "batch",
			expectedText: "Parse all TROFF files in a directory and save them to the database",
		},
	}

	for _, testCase := range testCases {
		// Using Go 1.22+ loop variable capture semantics
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			// #nosec G204 - This is a test with controlled input
			cmd := exec.Command("go", "run", "main.go", testCase.subcommand, "--help")
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
	cmd := exec.Command("go", "run", "main.go", "test", sampleFile)
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
	cmd := exec.Command("go", "run", "main.go", "save", "--dry-run", sampleFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	require.NoError(t, err, "Command should execute without error")

	// Check stdout for log output
	stdoutOutput := stdout.String()
	assert.Contains(t, stdoutOutput, "Dry run mode")
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
	cmd := exec.Command("go", "run", "main.go", "batch", "--dry-run", "--verbose", sampleDir)
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
