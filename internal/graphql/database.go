package graphql

import (
	"context"

	"github.com/carldunham/useful-cookery/internal/model"
)

// Database defines the database methods needed by the GraphQL resolvers.
//
//nolint:interfacebloat // This interface is intentionally large to cover all database operations
type Database interface {
	// User operations
	GetUser(ctx context.Context, userID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error

	// Category operations
	GetCategory(ctx context.Context, categoryID string) (*model.Category, error)
	CreateCategory(ctx context.Context, category *model.Category) error
	UpdateCategory(ctx context.Context, category *model.Category) error
	DeleteCategory(ctx context.Context, categoryID string) error

	// Review operations
	CreateReview(ctx context.Context, review *model.Review) error
	GetReview(ctx context.Context, reviewID string) (*model.Review, error)
	UpdateReview(ctx context.Context, review *model.Review) error
	DeleteReview(ctx context.Context, reviewID string) error

	// Recipe operations
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
	GetRecipes(ctx context.Context, filter map[string]string, first, offset int) ([]*model.Recipe, error)
	CountRecipes(ctx context.Context, filter map[string]string) (int, error)
	GetPopularRecipes(ctx context.Context, limit, offset int) ([]*model.Recipe, error)
	CreateRecipe(ctx context.Context, recipe *model.Recipe) error
	UpdateRecipe(ctx context.Context, recipe *model.Recipe) error
	DeleteRecipe(ctx context.Context, recipeID string) error

	// Raw query operations
	Query(ctx context.Context, query string, vars map[string]string, result any) error
}
