// Package main implements a command-line tool for managing users, categories, and recipes.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/cli"
	"github.com/carldunham/useful-cookery/internal/config"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/model"
)

// Error definitions.
var (
	// Interface errors.
	errDatabaseNotImplementsInterface = errors.New("database does not implement required interface")

	// Validation errors.
	errEmailPasswordRequired     = errors.New("email and password are required for creating admin user")
	errEmailPasswordNameRequired = errors.New("email, password, and name are required for creating user")
	errNameRequired              = errors.New("name is required for creating category")
	errTitleRequired             = errors.New("title is required for creating recipe")
	errRecipeIDRequired          = errors.New("recipe ID is required")
)

//nolint:gochecknoglobals // Using cobra command pattern which requires globals.
var (
	// Global flags.
	configFile string

	// User-related flags.
	email    string
	password string
	name     string
	role     string

	// Category and recipe flags.
	title       string
	desc        string
	difficulty  string
	cuisine     string
	prepTime    int
	cookTime    int
	servings    int
	ingredients []string
	steps       []string
	tags        []string

	// Root command.
	rootCmd = &cobra.Command{
		Use:   "useful-cookery-cli",
		Short: "Useful Cookery CLI",
		Long:  "Command-line tool for managing users, categories, and recipes in the Useful Cookery application",
	}

	// User commands.
	userCmd = &cobra.Command{
		Use:   "user",
		Short: "User management commands",
		Long:  "Commands for managing users in the Useful Cookery application",
	}

	// Create admin command.
	createAdminCmd = &cobra.Command{
		Use:   "create-admin",
		Short: "Create admin user",
		Long:  "Create an initial admin user with the specified email and password",
		RunE:  runCreateAdmin,
	}

	// Create user command.
	createUserCmd = &cobra.Command{
		Use:   "create",
		Short: "Create user",
		Long:  "Create a new user with the specified name, email, password, and role",
		RunE:  runCreateUser,
	}

	// List users command.
	listUsersCmd = &cobra.Command{
		Use:   "list",
		Short: "List users",
		Long:  "List all users in the system",
		RunE:  runListUsers,
	}

	// Category commands.
	categoryCmd = &cobra.Command{
		Use:   "category",
		Short: "Category management commands",
		Long:  "Commands for managing categories in the Useful Cookery application",
	}

	// Create category command.
	createCategoryCmd = &cobra.Command{
		Use:   "create",
		Short: "Create category",
		Long:  "Create a new category with the specified name and description",
		RunE:  runCreateCategory,
	}

	// List categories command.
	listCategoriesCmd = &cobra.Command{
		Use:   "list",
		Short: "List categories",
		Long:  "List all categories in the system",
		RunE:  runListCategories,
	}

	// Recipe commands.
	recipeCmd = &cobra.Command{
		Use:   "recipe",
		Short: "Recipe management commands",
		Long:  "Commands for managing recipes in the Useful Cookery application",
	}

	// Create recipe command.
	createRecipeCmd = &cobra.Command{
		Use:   "create",
		Short: "Create recipe",
		Long:  "Create a new recipe with the specified title, description, and other details",
		RunE:  runCreateRecipe,
	}

	// List recipes command.
	listRecipesCmd = &cobra.Command{
		Use:   "list",
		Short: "List recipes",
		Long:  "List all recipes in the system",
		RunE:  runListRecipes,
	}

	// Get recipe command.
	getRecipeCmd = &cobra.Command{
		Use:   "get",
		Short: "Get recipe",
		Long:  "Get a recipe by ID",
		RunE:  runGetRecipe,
	}

	// Recipe ID flag.
	recipeID string
)

// setupLogging initializes the logging system.
func setupLogging() {
	logHandler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	slog.SetDefault(slog.New(logHandler))
}

// setupUserCommands sets up the user-related commands and flags.
func setupUserCommands() {
	// Add user command flags.
	createAdminCmd.Flags().StringVarP(&email, "email", "e", "", "User email")
	createAdminCmd.Flags().StringVarP(&password, "password", "p", "", "User password")
	createAdminCmd.MarkFlagsRequiredTogether("email", "password")

	createUserCmd.Flags().StringVarP(&email, "email", "e", "", "User email")
	createUserCmd.Flags().StringVarP(&password, "password", "p", "", "User password")
	createUserCmd.Flags().StringVarP(&name, "name", "n", "", "User name")
	createUserCmd.Flags().StringVarP(&role, "role", "r", "USER", "User role (ADMIN or USER)")
	createUserCmd.MarkFlagsRequiredTogether("email", "password", "name")

	// Add commands to their parent command.
	userCmd.AddCommand(createAdminCmd)
	userCmd.AddCommand(createUserCmd)
	userCmd.AddCommand(listUsersCmd)
}

// setupCategoryCommands sets up the category-related commands and flags.
func setupCategoryCommands() {
	// Add category command flags.
	createCategoryCmd.Flags().StringVarP(&name, "name", "n", "", "Category name")
	createCategoryCmd.Flags().StringVarP(&desc, "description", "d", "", "Category description")
	if err := createCategoryCmd.MarkFlagRequired("name"); err != nil {
		slog.Error("Failed to mark flag as required", "flag", "name", "error", err)
	}

	// Add commands to their parent command.
	categoryCmd.AddCommand(createCategoryCmd)
	categoryCmd.AddCommand(listCategoriesCmd)
}

// setupRecipeCommands sets up the recipe-related commands and flags.
func setupRecipeCommands() {
	// Add recipe command flags.
	createRecipeCmd.Flags().StringVarP(&title, "title", "t", "", "Recipe title")
	createRecipeCmd.Flags().StringVarP(&desc, "description", "d", "", "Recipe description")
	createRecipeCmd.Flags().StringVar(&difficulty, "difficulty", "", "Recipe difficulty (Easy, Medium, Hard)")
	createRecipeCmd.Flags().StringVar(&cuisine, "cuisine", "", "Recipe cuisine type")
	createRecipeCmd.Flags().IntVar(&prepTime, "prep-time", 0, "Preparation time in minutes")
	createRecipeCmd.Flags().IntVar(&cookTime, "cook-time", 0, "Cooking time in minutes")
	createRecipeCmd.Flags().IntVar(&servings, "servings", 0, "Number of servings")
	createRecipeCmd.Flags().StringArrayVar(&ingredients, "ingredient", []string{},
		"Recipe ingredient (can be specified multiple times)")
	createRecipeCmd.Flags().StringArrayVar(&steps, "step", []string{}, "Recipe step (can be specified multiple times)")
	createRecipeCmd.Flags().StringArrayVar(&tags, "tag", []string{}, "Recipe tag (can be specified multiple times)")
	if err := createRecipeCmd.MarkFlagRequired("title"); err != nil {
		slog.Error("Failed to mark flag as required", "flag", "title", "error", err)
	}

	// Add get recipe command flags.
	getRecipeCmd.Flags().StringVarP(&recipeID, "id", "i", "", "Recipe ID")
	if err := getRecipeCmd.MarkFlagRequired("id"); err != nil {
		slog.Error("Failed to mark flag as required", "flag", "id", "error", err)
	}

	// Add commands to their parent command.
	recipeCmd.AddCommand(createRecipeCmd)
	recipeCmd.AddCommand(listRecipesCmd)
	recipeCmd.AddCommand(getRecipeCmd)
}

func main() {
	// Set up logging.
	setupLogging()

	// Add global flags.
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Path to config file")

	// Set up commands.
	setupUserCommands()
	setupCategoryCommands()
	setupRecipeCommands()

	// Add commands to root.
	rootCmd.AddCommand(userCmd)
	rootCmd.AddCommand(categoryCmd)
	rootCmd.AddCommand(recipeCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		slog.Error("Command execution failed", "error", err)
		os.Exit(1)
	}
}

// loadConfig loads the configuration from the specified file or environment variables.
func loadConfig(configFile string) (*config.Config, error) {
	if configFile != "" {
		viper.SetConfigFile(configFile)
	} else {
		// Use default config paths
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(".")
		viper.AddConfigPath("./config")
		viper.AddConfigPath("$HOME/.useful-cookery")
		viper.AddConfigPath("/etc/useful-cookery")
	}

	// Set default values
	// Database defaults - using Postgres as default
	viper.SetDefault("database.type", "postgres")
	viper.SetDefault(
		"database.connection_string", //nolint:ireturn,nolintlint // No idea why this has to be here.
		"postgres://postgres:postgres@localhost:5432/useful-cookery?sslmode=disable",
	)

	// Auth defaults
	viper.SetDefault("auth.jwt_secret", "change-me-in-production")
	viper.SetDefault("auth.token_expiry", "24h")

	// Read environment variables
	viper.SetEnvPrefix("UC")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Read configuration file
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &configFileNotFound) {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
		// Config file not found, using defaults and environment variables
		slog.Info("Config file not found, using defaults and environment variables")
	}

	// Parse configuration
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// Log database connection info
	slog.Info("Using database", "type", cfg.Database.Type, "connection", cfg.Database.ConnectionString)

	return &cfg, nil
}

// setupDatabase initializes the database connection.
//
//nolint:ireturn // Returning interfaces is a design choice for this function.
func setupDatabase(cfg *config.Config) (cli.Database, cli.AuthDatabase, error) {
	// Initialize database options
	dbOptions := database.Options{
		ConnectionString: cfg.Database.ConnectionString,
	}

	// Create database based on type
	dbInstance, err := database.CreateDatabase(dbtypes.DatabaseType(cfg.Database.Type), dbOptions)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Cast the database to the required interfaces
	var db cli.Database
	var authDB cli.AuthDatabase
	var ok bool

	if db, ok = dbInstance.(cli.Database); !ok {
		return nil, nil, fmt.Errorf("%w: cli.Database", errDatabaseNotImplementsInterface)
	}
	if authDB, ok = dbInstance.(cli.AuthDatabase); !ok {
		return nil, nil, fmt.Errorf("%w: cli.AuthDatabase", errDatabaseNotImplementsInterface)
	}

	slog.Info("Connected to database")
	return db, authDB, nil
}

// runCreateAdmin creates the initial admin user.
func runCreateAdmin(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	_, authDB, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize auth service
	authService := auth.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenExpiry,
		authDB,
	)

	// Create context
	ctx := context.Background()

	// Validate required flags
	if email == "" || password == "" {
		return errEmailPasswordRequired
	}

	// Create admin user
	if err := authService.CreateInitialAdminUser(ctx, email, password); err != nil {
		return fmt.Errorf("failed to create initial admin user: %w", err)
	}

	slog.Info("Admin user created successfully", "email", email)
	return nil
}

// runCreateUser creates a new user.
func runCreateUser(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	_, authDB, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize auth service
	authService := auth.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenExpiry,
		authDB,
	)

	// Create context
	ctx := context.Background()

	// Validate required flags
	if email == "" || password == "" || name == "" {
		return errEmailPasswordNameRequired
	}

	// Register user
	user, err := authService.RegisterUser(ctx, name, email, password)
	if err != nil {
		return fmt.Errorf("failed to register user: %w", err)
	}

	// Update role if needed
	if role == "ADMIN" && user.Role != model.AdminRole {
		user.Role = model.AdminRole
		if err := authService.UpdateUser(ctx, user); err != nil {
			return fmt.Errorf("failed to update user role: %w", err)
		}
	}

	slog.Info("User created successfully", "email", email, "role", user.Role)
	return nil
}

// runListUsers lists all users.
func runListUsers(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Get all users (no pagination)
	users, err := db.GetUsers(ctx, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to query users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found")
		return nil
	}

	fmt.Println("Users:")
	for _, user := range users {
		fmt.Printf("- ID: %s, Name: %s, Email: %s, Role: %s\n", user.ID, user.Name, user.Email, user.Role)
	}

	return nil
}

// runCreateCategory creates a new category.
func runCreateCategory(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Validate required flags
	if name == "" {
		return errNameRequired
	}

	// Create category
	category := &model.Category{
		Name:        name,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := db.CreateCategory(ctx, category); err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}

	slog.Info("Category created successfully", "name", name, "id", category.ID)
	return nil
}

// runListCategories lists all categories.
func runListCategories(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Get all categories (no pagination)
	categories, err := db.GetCategories(ctx, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to query categories: %w", err)
	}

	if len(categories) == 0 {
		fmt.Println("No categories found")
		return nil
	}

	fmt.Println("Categories:")
	for _, category := range categories {
		fmt.Printf("- ID: %s, Name: %s\n", category.ID, category.Name)
		if category.Description != "" {
			fmt.Printf("  Description: %s\n", category.Description)
		}
	}

	return nil
}

// createRecipeIngredient creates a new ingredient from a string in the format "unit: name" or just "name".
func createRecipeIngredient(ingredientStr string) model.DetailedIngredient {
	const ingredientParts = 2
	parts := strings.SplitN(ingredientStr, ":", ingredientParts)
	var unit, name string
	if len(parts) == ingredientParts {
		unit = strings.TrimSpace(parts[0])
		name = strings.TrimSpace(parts[1])
	} else {
		name = strings.TrimSpace(ingredientStr)
	}

	ingredient := model.DetailedIngredient{
		ID:   model.NewID(),
		Name: name,
		Unit: unit,
	}

	// Add a default unit if a unit was specified
	if unit != "" {
		unitID := model.NewID()
		ingredientUnit := model.IngredientUnit{
			ID:     unitID,
			System: "imperial", // Default system
			Value:  1.0,        // Default value
			Unit:   unit,
			IsMain: true,
		}
		ingredient.Units = append(ingredient.Units, ingredientUnit)
	}

	return ingredient
}

// createRecipeStep creates a new step with the given description and order index.
func createRecipeStep(index int, description string) model.Step {
	return model.Step{
		ID:          model.NewID(),
		OrderIndex:  index,
		Description: description,
	}
}

// createRecipeFromFlags creates a new recipe from the command-line flags.
func createRecipeFromFlags() *model.Recipe {
	now := time.Now()
	recipe := &model.Recipe{
		Title:       title,
		Description: desc,
		Difficulty:  difficulty,
		Cuisine:     cuisine,
		PrepTime:    prepTime,
		CookTime:    cookTime,
		Servings:    servings,
		Tags:        tags,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	// Add ingredients
	for _, ingredientStr := range ingredients {
		ingredient := createRecipeIngredient(ingredientStr)
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Add steps
	for i, stepDesc := range steps {
		step := createRecipeStep(i, stepDesc)
		recipe.Steps = append(recipe.Steps, step)
	}

	return recipe
}

// runCreateRecipe creates a new recipe.
func runCreateRecipe(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Validate required flags
	if title == "" {
		return errTitleRequired
	}

	// Create recipe from flags
	recipe := createRecipeFromFlags()

	// Save recipe to database
	if err := db.CreateRecipe(ctx, recipe); err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}

	slog.Info("Recipe created successfully", "title", title, "id", recipe.ID)
	return nil
}

// printRecipeBasicInfo prints the basic information of a recipe.
func printRecipeBasicInfo(recipe *model.Recipe) {
	fmt.Printf("Recipe: %s\n", recipe.Title)
	fmt.Printf("ID: %s\n", recipe.ID)
	if recipe.OriginalID != "" {
		fmt.Printf("Original ID: %s\n", recipe.OriginalID)
	}
	if recipe.Description != "" {
		fmt.Printf("Description: %s\n", recipe.Description)
	}
}

// printRecipeDetails2 prints the details of a recipe.
func printRecipeDetails2(recipe *model.Recipe) {
	if recipe.Difficulty != "" {
		fmt.Printf("Difficulty: %s\n", recipe.Difficulty)
	}
	if recipe.Cuisine != "" {
		fmt.Printf("Cuisine: %s\n", recipe.Cuisine)
	}
	if recipe.PrepTime > 0 {
		fmt.Printf("Prep Time: %d minutes\n", recipe.PrepTime)
	}
	if recipe.CookTime > 0 {
		fmt.Printf("Cook Time: %d minutes\n", recipe.CookTime)
	}
	if recipe.Servings > 0 {
		fmt.Printf("Servings: %d\n", recipe.Servings)
	}
}

// printRecipeIngredientsNoIndent prints the ingredients of a recipe without indentation.
func printRecipeIngredientsNoIndent(recipe *model.Recipe) {
	if len(recipe.Ingredients) == 0 {
		return
	}

	fmt.Println("\nIngredients:")
	for _, ingredient := range recipe.Ingredients {
		if ingredient.Unit != "" {
			fmt.Printf("- %s: %s\n", ingredient.Unit, ingredient.Name)
		} else {
			fmt.Printf("- %s\n", ingredient.Name)
		}
	}
}

// printRecipeStepsNoIndent prints the steps of a recipe without indentation.
func printRecipeStepsNoIndent(recipe *model.Recipe) {
	if len(recipe.Steps) == 0 {
		return
	}

	fmt.Println("\nSteps:")
	for i, step := range recipe.Steps {
		fmt.Printf("%d. %s\n", i+1, step.Description)
	}
}

// printRecipeTags prints the tags of a recipe.
func printRecipeTags(recipe *model.Recipe) {
	if len(recipe.Tags) > 0 {
		fmt.Printf("\nTags: %s\n", strings.Join(recipe.Tags, ", "))
	}
}

// printSingleRecipe prints a single recipe with a different format than the list view.
func printSingleRecipe(recipe *model.Recipe) {
	printRecipeBasicInfo(recipe)
	printRecipeDetails2(recipe)
	printRecipeIngredientsNoIndent(recipe)
	printRecipeStepsNoIndent(recipe)
	printRecipeTags(recipe)
}

// runGetRecipe gets a recipe by ID.
func runGetRecipe(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Validate required flags
	if recipeID == "" {
		return errRecipeIDRequired
	}

	// Get recipe
	recipe, err := db.GetRecipe(ctx, recipeID)
	if err != nil {
		return fmt.Errorf("failed to get recipe: %w", err)
	}

	// Print recipe details
	printSingleRecipe(recipe)

	return nil
}

// printRecipeDetails prints the details of a recipe in a formatted way.
func printRecipeDetails(recipe *model.Recipe, indent string) {
	fmt.Printf("%s- ID: %s, Title: %s\n", indent, recipe.ID, recipe.Title)
	if recipe.OriginalID != "" {
		fmt.Printf("%s  Original ID: %s\n", indent, recipe.OriginalID)
	}
	if recipe.Description != "" {
		fmt.Printf("%s  Description: %s\n", indent, recipe.Description)
	}
	if recipe.Difficulty != "" {
		fmt.Printf("%s  Difficulty: %s\n", indent, recipe.Difficulty)
	}
	if recipe.Cuisine != "" {
		fmt.Printf("%s  Cuisine: %s\n", indent, recipe.Cuisine)
	}
	if recipe.PrepTime > 0 {
		fmt.Printf("%s  Prep Time: %d minutes\n", indent, recipe.PrepTime)
	}
	if recipe.CookTime > 0 {
		fmt.Printf("%s  Cook Time: %d minutes\n", indent, recipe.CookTime)
	}
	if recipe.Servings > 0 {
		fmt.Printf("%s  Servings: %d\n", indent, recipe.Servings)
	}

	// Print ingredients
	printRecipeIngredients(recipe, indent)

	// Print steps
	printRecipeSteps(recipe, indent)

	// Print tags
	if len(recipe.Tags) > 0 {
		fmt.Printf("%s  Tags: %s\n", indent, strings.Join(recipe.Tags, ", "))
	}
}

// printRecipeIngredients prints the ingredients of a recipe.
func printRecipeIngredients(recipe *model.Recipe, indent string) {
	if len(recipe.Ingredients) == 0 {
		return
	}

	fmt.Printf("%s  Ingredients:\n", indent)
	for _, ingredient := range recipe.Ingredients {
		if ingredient.Unit != "" {
			fmt.Printf("%s    - %s: %s\n", indent, ingredient.Unit, ingredient.Name)
		} else {
			fmt.Printf("%s    - %s\n", indent, ingredient.Name)
		}
	}
}

// printRecipeSteps prints the steps of a recipe.
func printRecipeSteps(recipe *model.Recipe, indent string) {
	if len(recipe.Steps) == 0 {
		return
	}

	fmt.Printf("%s  Steps:\n", indent)
	for i, step := range recipe.Steps {
		fmt.Printf("%s    %d. %s\n", indent, i+1, step.Description)
	}
}

// runListRecipes lists all recipes.
func runListRecipes(_ *cobra.Command, _ []string) error {
	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Setup database
	db, _, err := setupDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Create context
	ctx := context.Background()

	// Get all recipes (no pagination)
	recipes, err := db.GetRecipes(ctx, nil, 0, 0)
	if err != nil {
		return fmt.Errorf("failed to query recipes: %w", err)
	}

	if len(recipes) == 0 {
		fmt.Println("No recipes found")
		return nil
	}

	fmt.Println("Recipes:")
	for _, recipe := range recipes {
		printRecipeDetails(recipe, "")
		fmt.Println() // Add a blank line between recipes
	}

	return nil
}
