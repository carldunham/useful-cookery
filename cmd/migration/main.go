// Package main implements a command-line tool for migrating TROFF recipes to the database.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"

	"github.com/carldunham/useful-cookery/internal/config"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
	"github.com/carldunham/useful-cookery/internal/troff"
)

// Database defines the database operations needed by the migration tool.
type Database interface {
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
	CreateRecipe(ctx context.Context, recipe *model.Recipe) error
}

// Constants.
const (
	filePermission = 0600 // Read/write permissions for owner only
)

//nolint:gochecknoglobals // Using cobra command pattern which requires globals.
var (
	// Global flags.
	verbose      bool
	dryRun       bool
	skipExisting bool
	maxErrors    int
	outputPath   string

	// Error definitions.
	errMaxErrorsReached  = errors.New("maximum error count reached, terminating")
	errUnsupportedDBType = errors.New("unsupported database type for migration")

	// Root command.
	rootCmd = &cobra.Command{
		Use:   "migration",
		Short: "TROFF recipe migration tool",
		Long:  "A tool for converting TROFF recipes to structured data and saving them to the database",
	}

	// Test command.
	testCmd = &cobra.Command{
		Use:   "test [file]",
		Short: "Test TROFF parsing",
		Long:  "Parse a TROFF file and display the raw content and formatted recipe",
		Args:  cobra.ExactArgs(1),
		RunE:  runTest,
	}

	// Save command.
	saveCmd = &cobra.Command{
		Use:   "save [file]",
		Short: "Save a recipe to the database",
		Long:  "Parse a TROFF file and save the recipe to the database",
		Args:  cobra.ExactArgs(1),
		RunE:  runSave,
	}

	// Batch command.
	batchCmd = &cobra.Command{
		Use:   "batch [directory]",
		Short: "Batch process TROFF files",
		Long:  "Parse all TROFF files in a directory and save them to the database",
		Args:  cobra.ExactArgs(1),
		RunE:  runBatch,
	}
)

func main() {
	// Set up logging.
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logHandler))

	// Add global flags.
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false,
		"Don't actually save to database, just parse and validate")
	rootCmd.PersistentFlags().BoolVar(&skipExisting, "skip-existing", false,
		"Skip recipes that already exist in the database (by ID)")

	// Add command-specific flags.
	testCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output file for formatted recipe")

	batchCmd.Flags().IntVar(&maxErrors, "max-errors", 1,
		"Maximum number of errors to tolerate before terminating (0 = don't terminate on errors)")

	// Add commands to root.
	rootCmd.AddCommand(testCmd)
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(batchCmd)

	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

// runTest parses a single TROFF file and displays raw and formatted output.
func runTest(_ *cobra.Command, args []string) error {
	filePath := args[0]

	// Read the raw TROFF content.
	rawContent, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	// Parse the TROFF file.
	recipe, err := troff.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse TROFF file: %w", err)
	}

	// Display the raw TROFF content and formatted recipe.
	fmt.Println("=== Raw TROFF Content ===")
	fmt.Println(string(rawContent))
	fmt.Println("=== Formatted Recipe ===")
	printRecipe(recipe)

	// If output path is provided, write the formatted recipe to file.
	if outputPath != "" {
		// Create a simple text representation of the recipe.
		output := formatRecipeAsText(recipe)
		if err := os.WriteFile(outputPath, []byte(output), filePermission); err != nil {
			return fmt.Errorf("failed to write output file: %w", err)
		}
		fmt.Println("Formatted recipe written to file", "path", outputPath)
	}

	return nil
}

// runSave parses a single TROFF file and saves it to the database.
func runSave(_ *cobra.Command, args []string) error {
	filePath := args[0]

	// Parse the TROFF file.
	recipe, err := troff.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to parse TROFF file: %w", err)
	}

	// If dry run, just print the recipe and return.
	if dryRun {
		fmt.Println("Dry run mode - recipe would be saved to database, originalID:", recipe.OriginalID)
		printRecipe(recipe)
		return nil
	}

	// Connect to the database.
	db, err := connectToDatabase()
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Check if recipe already exists.
	if skipExisting {
		existingRecipe, err := db.GetRecipe(context.Background(), recipe.ID)
		if err == nil && existingRecipe != nil {
			fmt.Println("Recipe already exists, skipping", "id", recipe.ID)
			return nil
		}
	}

	// Set creation and update timestamps.
	now := time.Now()
	recipe.CreatedAt = now
	recipe.UpdatedAt = now

	// Save the recipe to the database.
	if err := db.CreateRecipe(context.Background(), &recipe); err != nil {
		return fmt.Errorf("failed to save recipe to database: %w", err)
	}

	fmt.Println("Recipe saved to database. ", "id:", recipe.ID, "originalID:", recipe.OriginalID, "title:", recipe.Title)
	return nil
}

// runBatch parses all TROFF files in a directory and saves them to the database.
//
//nolint:gocognit,cyclop,funlen // Complex function that would require significant refactoring to simplify
func runBatch(_ *cobra.Command, args []string) error {
	dirPath := args[0]

	// Connect to the database.
	var db Database
	var err error
	if !dryRun {
		db, err = connectToDatabase()
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
	}

	// Get all files in the directory.
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	// Process each file.
	errorCount := 0
	successCount := 0
	skipCount := 0

	for _, file := range files {
		// Skip directories.
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(dirPath, file.Name())
		if verbose {
			fmt.Println("Processing file", "path", filePath)
		}

		// Parse the TROFF file.
		recipe, err := troff.ParseFile(filePath)
		if err != nil {
			slog.Error("Failed to parse file", "path", filePath, "error", err)
			errorCount++
			if maxErrors > 0 && errorCount >= maxErrors {
				return fmt.Errorf("%w: %d", errMaxErrorsReached, maxErrors)
			}
			continue
		}

		// If dry run, just print the recipe and continue.
		if dryRun {
			if verbose {
				fmt.Println("Successfully parsed file (dry run)", "path", filePath, "originalID", recipe.OriginalID)
			}
			successCount++
			continue
		}

		// Check if recipe already exists.
		if skipExisting {
			existingRecipe, err := db.GetRecipe(context.Background(), recipe.ID)
			if err == nil && existingRecipe != nil {
				if verbose {
					fmt.Println("Recipe already exists, skipping", "originalID", recipe.OriginalID, "path", filePath)
				}
				skipCount++
				continue
			}
		}

		// Set creation and update timestamps.
		now := time.Now()
		recipe.CreatedAt = now
		recipe.UpdatedAt = now

		// Save the recipe to the database.
		if err := db.CreateRecipe(context.Background(), &recipe); err != nil {
			slog.Error("Failed to save recipe to database", "path", filePath, "error", err)
			errorCount++
			if maxErrors > 0 && errorCount >= maxErrors {
				return fmt.Errorf("%w: %d", errMaxErrorsReached, maxErrors)
			}
			continue
		}

		if verbose {
			fmt.Println("Successfully saved recipe", "path", filePath, "id", recipe.ID, "originalID", recipe.OriginalID)
		}
		successCount++
	}

	// Print summary.
	fmt.Println("Batch processing complete",
		"successful", successCount,
		"skipped", skipCount,
		"errors", errorCount)
	return nil
}

// connectToDatabase connects to the database.
//
//nolint:ireturn // Intentionally returning interface for database abstraction
func connectToDatabase() (Database, error) {
	// Load configuration.
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize database options
	dbOptions := database.Options{
		ConnectionString: cfg.Database.ConnectionString,
	}

	// If Database connection string is not set, fall back to DGraph connection string
	// for backward compatibility
	if dbOptions.ConnectionString == "" {
		dbOptions.ConnectionString = cfg.DGraph.ConnectionString
	}

	// Create database based on type
	var db Database
	switch cfg.Database.Type {
	case "dgraph", "":
		dgraphDB, err := database.NewDGraphDatabase(dbOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		db = dgraphDB
	case "memory":
		memoryDB, err := database.NewInMemoryDatabase(dbOptions)
		if err != nil {
			return nil, fmt.Errorf("failed to connect to database: %w", err)
		}
		db = memoryDB
	default:
		return nil, fmt.Errorf("%w: %s", errUnsupportedDBType, cfg.Database.Type)
	}

	return db, nil
}

// printRecipe logs a recipe in a readable format.
func printRecipe(recipe model.Recipe) {
	fmt.Println("Recipe details", "id:", recipe.ID, "originalID:", recipe.OriginalID, "title:", recipe.Title)

	if recipe.Description != "" {
		fmt.Println("Description", "text", recipe.Description)
	}

	if len(recipe.Categories) > 0 {
		categories := make([]string, 0, len(recipe.Categories))
		for _, category := range recipe.Categories {
			categories = append(categories, category.Name)
		}
		fmt.Println("Categories", "list", categories)
	}

	if recipe.Author != nil {
		fmt.Println("Author", "name", recipe.Author.Name)
	}

	if len(recipe.Ingredients) > 0 {
		fmt.Println("Ingredients")
		for _, ingredient := range recipe.Ingredients {
			fmt.Println("- Ingredient", "unit", ingredient.Unit, "name", ingredient.Name)
		}
	}

	if len(recipe.Steps) > 0 {
		fmt.Println("Steps")
		for _, step := range recipe.Steps {
			fmt.Println("- Step", "number", step.OrderIndex+1, "description", step.Description)
		}
	}

	if recipe.Notes != "" {
		fmt.Println("Notes", "text", recipe.Notes)
	}
}

// formatRecipeAsText formats a recipe as plain text.
func formatRecipeAsText(recipe model.Recipe) string {
	var result string

	result += "# " + recipe.Title + "\n"
	result += "Original ID: " + recipe.OriginalID + "\n\n"

	if recipe.Description != "" {
		result += recipe.Description + "\n\n"
	}

	if len(recipe.Categories) > 0 {
		result += "## Categories\n"
		for _, category := range recipe.Categories {
			result += "- " + category.Name + "\n"
		}
		result += "\n"
	}

	if recipe.Author != nil {
		result += "## Author\n" + recipe.Author.Name + "\n\n"
	}

	if len(recipe.Ingredients) > 0 {
		result += "## Ingredients\n"
		for _, ingredient := range recipe.Ingredients {
			result += "- " + ingredient.Unit + ": " + ingredient.Name + "\n"
		}
		result += "\n"
	}

	if len(recipe.Steps) > 0 {
		result += "## Steps\n"
		for _, step := range recipe.Steps {
			result += fmt.Sprintf("%d. ", step.OrderIndex+1) + step.Description + "\n"
		}
		result += "\n"
	}

	if recipe.Notes != "" {
		result += "## Notes\n" + recipe.Notes + "\n"
	}

	return result
}
