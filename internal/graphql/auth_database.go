package graphql

import (
	"context"

	"github.com/carldunham/useful-cookery/internal/model"
)

// AuthDatabase defines the database methods needed by the auth service in the GraphQL application.
type AuthDatabase interface {
	// User operations
	GetUser(ctx context.Context, userID string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error)
	CreateUser(ctx context.Context, user *model.User) error
	UpdateUser(ctx context.Context, user *model.User) error

	// Recipe operations (needed for permission checks)
	GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error)
}
