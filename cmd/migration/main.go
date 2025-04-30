// Package main implements a command-line tool for migrating TROFF recipes to the database
// and managing database schema migrations.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/spf13/cobra"

	"github.com/carldunham/useful-cookery/cmd/migration/config"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
	"github.com/carldunham/useful-cookery/internal/troff"
)

// Database defines the database operations needed by the migration tool.
type Database interface {
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
	GetRecipes(ctx context.Context, filter map[string]string, first, offset int) ([]*model.Recipe, error)
	CreateRecipe(ctx context.Context, recipe *model.Recipe) error
	UpdateRecipe(ctx context.Context, recipe *model.Recipe) error
}

// Constants.
const (
	filePermission = 0600 // Read/write permissions for owner only
	migrationsPath = "file://db/migrations"
)

//nolint:gochecknoglobals // Using cobra command pattern which requires globals.
var (
	// Global flags.
	verbose      bool
	dryRun       bool
	skipExisting bool
	maxErrors    int
	outputPath   string
	steps        int

	// Error definitions.
	errMaxErrorsReached              = errors.New("maximum error count reached, terminating")
	errUnsupportedDBType             = errors.New("unsupported database type for migration")
	errNoMigrationFound              = errors.New("no migration file found for version")
	errUnsupportedDBTypeForMigration = errors.New("migrations only supported for postgres database type")

	// Root command.
	rootCmd = &cobra.Command{
		Use:   "migration",
		Short: "Migration tool",
		Long:  "A tool for converting TROFF recipes to structured data and managing database schema migrations",
	}

	// TROFF recipe migration commands.
	troffCmd = &cobra.Command{
		Use:   "troff",
		Short: "TROFF recipe migration commands",
		Long:  "Commands for converting TROFF recipes to structured data and saving them to the database",
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

	// Database schema migration commands.
	dbCmd = &cobra.Command{
		Use:   "db",
		Short: "Database schema migration commands",
		Long:  "Commands for managing database schema migrations",
	}

	// Create migration command.
	createCmd = &cobra.Command{
		Use:   "create [name]",
		Short: "Create a new migration",
		Long:  "Create a new empty migration with up and down files",
		Args:  cobra.ExactArgs(1),
		RunE:  runCreate,
	}

	// Up command.
	upCmd = &cobra.Command{
		Use:   "up",
		Short: "Run migrations up",
		Long:  "Apply all or a limited number of up migrations",
		RunE:  runUp,
	}

	// Down command.
	downCmd = &cobra.Command{
		Use:   "down",
		Short: "Run migrations down",
		Long:  "Apply all or a limited number of down migrations",
		RunE:  runDown,
	}

	// Status command.
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "Show migration status",
		Long:  "Show the current migration version and dirty state",
		RunE:  runStatus,
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

	// Add TROFF command flags
	troffCmd.PersistentFlags().BoolVar(&dryRun, "dry-run", false,
		"Don't actually save to database, just parse and validate")
	troffCmd.PersistentFlags().BoolVar(&skipExisting, "skip-existing", false,
		"Skip recipes that already exist in the database (by ID)")

	// Add command-specific flags.
	testCmd.Flags().StringVarP(&outputPath, "output", "o", "", "Path to output file for formatted recipe")

	batchCmd.Flags().IntVar(&maxErrors, "max-errors", 1,
		"Maximum number of errors to tolerate before terminating (0 = don't terminate on errors)")

	// Add database migration command flags
	upCmd.Flags().IntVarP(&steps, "steps", "n", 0, "Number of migrations to apply (0 = all)")
	downCmd.Flags().IntVarP(&steps, "steps", "n", 0, "Number of migrations to apply (0 = all)")

	// Add TROFF commands to troff command
	troffCmd.AddCommand(testCmd)
	troffCmd.AddCommand(saveCmd)
	troffCmd.AddCommand(batchCmd)

	// Add database migration commands to db command
	dbCmd.AddCommand(createCmd)
	dbCmd.AddCommand(upCmd)
	dbCmd.AddCommand(downCmd)
	dbCmd.AddCommand(statusCmd)

	// Add commands to root
	rootCmd.AddCommand(troffCmd)
	rootCmd.AddCommand(dbCmd)

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

// checkExistingRecipe checks if a recipe already exists in the database.
func checkExistingRecipe(db Database, recipe *model.Recipe) (*model.Recipe, bool) {
	// First try to find by ID
	existingRecipe, err := db.GetRecipe(context.Background(), recipe.ID)
	if err == nil && existingRecipe != nil {
		return existingRecipe, true
	}

	// If not found by ID, try to find by OriginalID
	recipes, err := db.GetRecipes(context.Background(), map[string]string{"originalID": recipe.OriginalID}, 1, 0)
	if err == nil && len(recipes) > 0 {
		existingRecipe := recipes[0]
		// Use the existing ID
		recipe.ID = existingRecipe.ID
		return existingRecipe, true
	}

	return nil, false
}

// saveRecipeToDatabase saves or updates a recipe in the database.
func saveRecipeToDatabase(db Database, recipe *model.Recipe, shouldUpdate bool, existingRecipe *model.Recipe) error {
	// Set creation and update timestamps.
	now := time.Now()
	if shouldUpdate {
		recipe.CreatedAt = existingRecipe.CreatedAt
		recipe.UpdatedAt = now
	} else {
		recipe.CreatedAt = now
		recipe.UpdatedAt = now
	}

	// Save the recipe to the database.
	if shouldUpdate {
		if err := db.UpdateRecipe(context.Background(), recipe); err != nil {
			return fmt.Errorf("failed to update recipe in database: %w", err)
		}
		fmt.Println(
			"Recipe updated in database. ",
			"id:", recipe.ID,
			"originalID:", recipe.OriginalID,
			"title:", recipe.Title,
		)
	} else {
		if err := db.CreateRecipe(context.Background(), recipe); err != nil {
			return fmt.Errorf("failed to save recipe to database: %w", err)
		}
		fmt.Println(
			"Recipe saved to database. ",
			"id:", recipe.ID,
			"originalID:", recipe.OriginalID,
			"title:", recipe.Title,
		)
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

	// Check if recipe already exists
	existingRecipe, shouldUpdate := checkExistingRecipe(db, &recipe)

	// Skip if requested and recipe exists
	if skipExisting && shouldUpdate {
		fmt.Println("Recipe already exists, skipping", "id", recipe.ID, "originalID", recipe.OriginalID)
		return nil
	}

	// Save or update the recipe
	return saveRecipeToDatabase(db, &recipe, shouldUpdate, existingRecipe)
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
		"errors", errorCount,
	)
	return nil
}

// connectToDatabase connects to the database.
//
//nolint:ireturn,nolintlint // Intentionally returning the Database abstraction.
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

	// Create database based on type
	dbInstance, err := database.CreateDatabase(database.Type(cfg.Database.Type), dbOptions)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Cast the database to the required interface
	db, ok := dbInstance.(Database)
	if !ok {
		return nil, fmt.Errorf("%w: %s", errUnsupportedDBType, cfg.Database.Type)
	}

	return db, nil
}

// printRecipeBasicInfo prints the basic information of a recipe.
func printRecipeBasicInfo(recipe model.Recipe) {
	fmt.Println("Recipe details", "id:", recipe.ID, "originalID:", recipe.OriginalID, "title:", recipe.Title)

	if recipe.Description != "" {
		fmt.Println("Description", "text", recipe.Description)
	}

	if recipe.Difficulty != "" {
		fmt.Println("Difficulty", "level", recipe.Difficulty)
	}

	if recipe.Cuisine != "" {
		fmt.Println("Cuisine", "type", recipe.Cuisine)
	}
}

// printRecipeTimingInfo prints the timing information of a recipe.
func printRecipeTimingInfo(recipe model.Recipe) {
	if recipe.PrepTime > 0 {
		fmt.Println("Prep Time", "minutes", recipe.PrepTime)
	}

	if recipe.CookTime > 0 {
		fmt.Println("Cook Time", "minutes", recipe.CookTime)
	}

	if recipe.Servings > 0 {
		fmt.Println("Servings", "count", recipe.Servings)
	}
}

// printRecipeCategories prints the categories of a recipe.
func printRecipeCategories(recipe model.Recipe) {
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
}

// printRecipeIngredients prints the ingredients of a recipe.
func printRecipeIngredients(recipe model.Recipe) {
	if len(recipe.Ingredients) > 0 {
		fmt.Println("Ingredients")
		for _, ingredient := range recipe.Ingredients {
			fmt.Println("- Ingredient", "unit", ingredient.Unit, "name", ingredient.Name)

			if len(ingredient.Units) > 0 {
				fmt.Println("  Units:")
				for _, unit := range ingredient.Units {
					fmt.Printf("  - %s: %g %s (system: %s, main: %t)\n",
						unit.ID, unit.Value, unit.Unit, unit.System, unit.IsMain)
				}
			}
		}
	}
}

// printRecipeSteps prints the steps of a recipe.
func printRecipeSteps(recipe model.Recipe) {
	if len(recipe.Steps) > 0 {
		fmt.Println("Steps")
		for _, step := range recipe.Steps {
			fmt.Println("- Step", "number", step.OrderIndex+1, "description", step.Description)
		}
	}
}

// printRecipeNotesAndTags prints the notes and tags of a recipe.
func printRecipeNotesAndTags(recipe model.Recipe) {
	if recipe.Notes != "" {
		fmt.Println("Notes", "text", recipe.Notes)
	}

	if len(recipe.Tags) > 0 {
		fmt.Println("Tags", "list", recipe.Tags)
	}
}

// printRecipe logs a recipe in a readable format.
func printRecipe(recipe model.Recipe) {
	printRecipeBasicInfo(recipe)
	printRecipeTimingInfo(recipe)
	printRecipeCategories(recipe)
	printRecipeIngredients(recipe)
	printRecipeSteps(recipe)
	printRecipeNotesAndTags(recipe)
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

// runCreate creates a new migration file.
func runCreate(_ *cobra.Command, args []string) error {
	name := args[0]

	// Create migrations directory if it doesn't exist
	const dirPermission = 0700 // Owner can read, write, and execute
	if err := os.MkdirAll("db/migrations", dirPermission); err != nil {
		return fmt.Errorf("failed to create migrations directory: %w", err)
	}

	// Get current timestamp for versioning
	timestamp := time.Now().Unix()

	// Create migration file names
	upFile := fmt.Sprintf("db/migrations/%d_%s.up.sql", timestamp, name)
	downFile := fmt.Sprintf("db/migrations/%d_%s.down.sql", timestamp, name)

	// Create empty migration files
	if err := os.WriteFile(upFile, []byte("-- Migration Up\n\n"), filePermission); err != nil {
		return fmt.Errorf("failed to create up migration file: %w", err)
	}

	if err := os.WriteFile(downFile, []byte("-- Migration Down\n\n"), filePermission); err != nil {
		return fmt.Errorf("failed to create down migration file: %w", err)
	}

	fmt.Printf("Created migration files:\n  %s\n  %s\n", upFile, downFile)
	return nil
}

// getDatabaseURL gets the database URL from the config.
func getDatabaseURL() (string, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return "", fmt.Errorf("failed to load configuration: %w", err)
	}

	// Check database type
	if cfg.Database.Type != "postgres" {
		return "", fmt.Errorf("%w: got %s", errUnsupportedDBTypeForMigration, cfg.Database.Type)
	}

	return cfg.Database.ConnectionString, nil
}

// createMigrate creates a new migrate instance.
func createMigrate() (*migrate.Migrate, error) {
	dbURL, err := getDatabaseURL()
	if err != nil {
		return nil, err
	}

	// Create migrate instance
	m, err := migrate.New(migrationsPath, dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m, nil
}

// runUp runs migrations up.
func runUp(_ *cobra.Command, _ []string) error {
	mig, err := createMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	if steps > 0 {
		if err := mig.Steps(steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		fmt.Printf("Applied %d migrations\n", steps)
	} else {
		if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to apply migrations: %w", err)
		}
		fmt.Println("Applied all migrations")
	}

	return nil
}

// runDown runs migrations down.
func runDown(_ *cobra.Command, _ []string) error {
	mig, err := createMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	if steps > 0 {
		if err := mig.Steps(-steps); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		fmt.Printf("Reverted %d migrations\n", steps)
	} else {
		if err := mig.Down(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
			return fmt.Errorf("failed to revert migrations: %w", err)
		}
		fmt.Println("Reverted all migrations")
	}

	return nil
}

// runStatus shows the current migration version.
func runStatus(_ *cobra.Command, _ []string) error {
	mig, err := createMigrate()
	if err != nil {
		return err
	}
	defer mig.Close()

	version, dirty, err := mig.Version()
	if err != nil {
		if errors.Is(err, migrate.ErrNilVersion) {
			fmt.Println("No migrations applied")
			return nil
		}
		return fmt.Errorf("failed to get migration version: %w", err)
	}

	// Get migration level (name) from the migration files
	migrationName, err := getMigrationName(version)
	if err != nil {
		slog.Warn("Could not determine migration name", "error", err)
	}

	fmt.Printf("Current migration version: %d\n", version)
	if migrationName != "" {
		fmt.Printf("Migration level: %s\n", migrationName)
	}
	fmt.Printf("Dirty: %t\n", dirty)

	return nil
}

// getMigrationName extracts the migration name from the migration files for a given version.
func getMigrationName(version uint) (string, error) {
	// List migration files
	files, err := os.ReadDir("db/migrations")
	if err != nil {
		return "", fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Look for a file with the matching version number
	for _, file := range files {
		var fileVersion uint
		parts := strings.Split(file.Name(), "_")
		if len(parts) > 0 {
			if _, err := fmt.Sscanf(parts[0], "%d", &fileVersion); err != nil {
				continue
			}
			if fileVersion == version && strings.HasSuffix(file.Name(), ".up.sql") {
				// Extract the name part from the filename
				namePart := strings.Join(parts[1:], "_")
				namePart = strings.TrimSuffix(namePart, ".up.sql")
				return namePart, nil
			}
		}
	}
	return "", fmt.Errorf("%w: %d", errNoMigrationFound, version)
}
