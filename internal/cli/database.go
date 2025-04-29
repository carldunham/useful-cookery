package cli

import (
	"context"

	"github.com/carldunham/useful-cookery/internal/model"
)

// Database defines the database methods needed by the CLI application.
type Database interface {
	// Category operations
	CreateCategory(ctx context.Context, category *model.Category) error
	GetCategories(ctx context.Context, limit, offset int) ([]*model.Category, error)

	// Recipe operations
	CreateRecipe(ctx context.Context, recipe *model.Recipe) error
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
	GetRecipes(ctx context.Context, filter map[string]string, limit, offset int) ([]*model.Recipe, error)

	// User operations
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)

	// Query operations
	Query(ctx context.Context, query string, vars map[string]string, result any) error
}
