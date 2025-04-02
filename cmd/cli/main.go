package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/config"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/model"
)

const (
	cmdCreateAdmin    = "create-admin"
	cmdCreateUser     = "create-user"
	cmdCreateCategory = "create-category"
	cmdCreateRecipe   = "create-recipe"
	cmdListUsers      = "list-users"
	cmdListCategories = "list-categories"
	cmdListRecipes    = "list-recipes"
)

func main() {
	// Define command line flags
	var (
		configFile string
		command    string
		email      string
		password   string
		name       string
		role       string
		title      string
		desc       string
	)

	// Global flags
	pflag.StringVarP(&configFile, "config", "c", "", "Path to config file")
	pflag.StringVarP(&command, "command", "m", "", "Command to execute")

	// User-related flags
	pflag.StringVarP(&email, "email", "e", "", "User email")
	pflag.StringVarP(&password, "password", "p", "", "User password")
	pflag.StringVarP(&name, "name", "n", "", "User or category name")
	pflag.StringVarP(&role, "role", "r", "USER", "User role (ADMIN or USER)")

	// Category and recipe flags
	pflag.StringVarP(&title, "title", "t", "", "Recipe title")
	pflag.StringVarP(&desc, "description", "d", "", "Description")

	// Parse flags
	pflag.Parse()

	// Print usage if no command is provided
	if command == "" {
		printUsage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := loadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to DGraph
	dgraphClient, err := database.NewDGraphClient(cfg.DGraph.ConnectionString)
	if err != nil {
		log.Fatalf("Failed to connect to DGraph: %v", err)
	}
	log.Println("Connected to DGraph")

	// Initialize auth service
	authService := auth.NewService(
		cfg.Auth.JWTSecret,
		cfg.Auth.TokenExpiry,
		dgraphClient,
	)

	// Create context
	ctx := context.Background()

	// Execute command
	switch command {
	case cmdCreateAdmin:
		if email == "" || password == "" {
			log.Fatalf("Email and password are required for creating admin user")
		}
		if err := createAdminUser(ctx, authService, email, password); err != nil {
			log.Fatalf("Failed to create admin user: %v", err)
		}
		log.Printf("Admin user created successfully: %s", email)

	case cmdCreateUser:
		if email == "" || password == "" || name == "" {
			log.Fatalf("Email, password, and name are required for creating user")
		}
		if err := createUser(ctx, authService, name, email, password, role); err != nil {
			log.Fatalf("Failed to create user: %v", err)
		}
		log.Printf("User created successfully: %s", email)

	case cmdCreateCategory:
		if name == "" {
			log.Fatalf("Name is required for creating category")
		}
		if err := createCategory(ctx, dgraphClient, name, desc); err != nil {
			log.Fatalf("Failed to create category: %v", err)
		}
		log.Printf("Category created successfully: %s", name)

	case cmdCreateRecipe:
		if title == "" {
			log.Fatalf("Title is required for creating recipe")
		}
		if err := createRecipe(ctx, dgraphClient, title, desc); err != nil {
			log.Fatalf("Failed to create recipe: %v", err)
		}
		log.Printf("Recipe created successfully: %s", title)

	case cmdListUsers:
		if err := listUsers(ctx, dgraphClient); err != nil {
			log.Fatalf("Failed to list users: %v", err)
		}

	case cmdListCategories:
		if err := listCategories(ctx, dgraphClient); err != nil {
			log.Fatalf("Failed to list categories: %v", err)
		}

	case cmdListRecipes:
		if err := listRecipes(ctx, dgraphClient); err != nil {
			log.Fatalf("Failed to list recipes: %v", err)
		}

	default:
		log.Fatalf("Unknown command: %s", command)
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
	viper.SetDefault("dgraph.connection_string", "localhost:9080")
	viper.SetDefault("auth.jwt_secret", "change-me-in-production")
	viper.SetDefault("auth.token_expiry", "24h")

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
	var cfg config.Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	// We don't need to modify the DGraph connection string
	// The dgo.Open function expects a URL with the "dgraph://" scheme
	fmt.Fprintf(os.Stderr, "Using DGraph connection string: %s\n", cfg.DGraph.ConnectionString)

	return &cfg, nil
}

// createAdminUser creates the initial admin user.
func createAdminUser(ctx context.Context, authService *auth.Service, email, password string) error {
	return authService.CreateInitialAdminUser(ctx, email, password)
}

// createUser creates a new user.
func createUser(ctx context.Context, authService *auth.Service, name, email, password, role string) error {
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

	return nil
}

// createCategory creates a new category.
func createCategory(ctx context.Context, dbClient *database.DGraphClient, name, description string) error {
	category := &model.Category{
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return dbClient.CreateCategory(ctx, category)
}

// createRecipe creates a new recipe.
func createRecipe(ctx context.Context, dbClient *database.DGraphClient, title, description string) error {
	recipe := &model.Recipe{
		Title:       title,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return dbClient.CreateRecipe(ctx, recipe)
}

// listUsers lists all users
func listUsers(ctx context.Context, dbClient *database.DGraphClient) error {
	// Get all users (no pagination)
	users, err := dbClient.GetUsers(ctx, 0, 0)
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

// listCategories lists all categories
func listCategories(ctx context.Context, dbClient *database.DGraphClient) error {
	// Query for categories - use a more specific query to avoid getting users
	query := `
	{
		categories(func: has(name)) @filter(NOT has(email)) {
			uid
			name
			description
			createdAt
		}
	}`

	var result struct {
		Categories []*model.Category `json:"categories"`
	}

	err := dbClient.Query(ctx, query, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to query categories: %w", err)
	}

	if len(result.Categories) == 0 {
		fmt.Println("No categories found")
		return nil
	}

	fmt.Println("Categories:")
	for _, category := range result.Categories {
		fmt.Printf("- ID: %s, Name: %s\n", category.ID, category.Name)
		if category.Description != "" {
			fmt.Printf("  Description: %s\n", category.Description)
		}
	}

	return nil
}

// listRecipes lists all recipes
func listRecipes(ctx context.Context, dbClient *database.DGraphClient) error {
	// Query for recipes
	q := `
	{
		recipes(func: has(title)) {
			uid
			title
			description
			createdAt
			updatedAt
		}
	}`

	var result struct {
		Recipes []*model.Recipe `json:"recipes"`
	}

	err := dbClient.Query(ctx, q, nil, &result)
	if err != nil {
		return fmt.Errorf("failed to query recipes: %w", err)
	}

	if len(result.Recipes) == 0 {
		fmt.Println("No recipes found")
		return nil
	}

	fmt.Println("Recipes:")
	for _, recipe := range result.Recipes {
		fmt.Printf("- ID: %s, Title: %s\n", recipe.ID, recipe.Title)
		if recipe.Description != "" {
			fmt.Printf("  Description: %s\n", recipe.Description)
		}
	}

	return nil
}

// printUsage prints the usage information
func printUsage() {
	fmt.Println("Useful Cookery CLI")
	fmt.Println("\nUsage:")
	fmt.Println("  useful-cookery-cli --command <command> [options]")
	fmt.Println("\nCommands:")
	fmt.Printf("  %s\tCreate initial admin user\n", cmdCreateAdmin)
	fmt.Printf("  %s\tCreate a new user\n", cmdCreateUser)
	fmt.Printf("  %s\tCreate a new category\n", cmdCreateCategory)
	fmt.Printf("  %s\tCreate a new recipe\n", cmdCreateRecipe)
	fmt.Printf("  %s\tList all users\n", cmdListUsers)
	fmt.Printf("  %s\tList all categories\n", cmdListCategories)
	fmt.Printf("  %s\tList all recipes\n", cmdListRecipes)
	fmt.Println("\nOptions:")
	fmt.Println("  -c, --config string       Path to config file")
	fmt.Println("  -m, --command string      Command to execute")
	fmt.Println("  -e, --email string        User email")
	fmt.Println("  -p, --password string     User password")
	fmt.Println("  -n, --name string         User or category name")
	fmt.Println("  -r, --role string         User role (ADMIN or USER)")
	fmt.Println("  -t, --title string        Recipe title")
	fmt.Println("  -d, --description string  Description")
	fmt.Println("\nExamples:")
	fmt.Println("  Create admin user:")
	fmt.Println("    useful-cookery-cli --command create-admin --email admin@example.com --password securepassword")
	fmt.Println("  Create regular user:")
	fmt.Println("    useful-cookery-cli --command create-user --name \"John Doe\" --email john@example.com --password userpassword")
	fmt.Println("  Create category:")
	fmt.Println("    useful-cookery-cli --command create-category --name \"Desserts\" --description \"Sweet treats\"")
	fmt.Println("  Create recipe:")
	fmt.Println("    useful-cookery-cli --command create-recipe --title \"Chocolate Cake\" --description \"Delicious chocolate cake recipe\"")
	fmt.Println("  List users:")
	fmt.Println("    useful-cookery-cli --command list-users")
	fmt.Println("  List categories:")
	fmt.Println("    useful-cookery-cli --command list-categories")
	fmt.Println("  List recipes:")
	fmt.Println("    useful-cookery-cli --command list-recipes")
}
