package resolvers

import (
	"errors"

	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Common sentinel errors for resolvers.
var (
	ErrNotAuthenticated      = errors.New("not authenticated")
	ErrPermissionDenied      = errors.New("permission denied")
	ErrRecipeNotFound        = errors.New("recipe not found")
	ErrRecipeNotFoundInSaved = errors.New("recipe not found in saved recipes")
	ErrReviewNotFound        = errors.New("review not found")
	ErrCategoryNotFound      = errors.New("category not found")
	ErrUserNotFound          = errors.New("user not found")
)

// Resolver is the root resolver for GraphQL operations.
type Resolver struct {
	DB                   *database.DGraphClient
	AIService            *ai.Service
	AuthService          *auth.Service
	RecipeResolver       *RecipeResolver
	UserResolver         *UserResolver
	SearchResolver       *SearchResolver
	SubscriptionResolver *SubscriptionResolver
}

// NewRootResolver creates a new root resolver.
func NewRootResolver(dbClient *database.DGraphClient, aiService *ai.Service, authService *auth.Service) *Resolver {
	// Create the root resolver first (needed for circular dependencies)
	resolver := &Resolver{
		DB:          dbClient,
		AIService:   aiService,
		AuthService: authService,
	}

	// Create specialized resolvers
	resolver.RecipeResolver = NewRecipeResolver(dbClient)
	resolver.UserResolver = NewUserResolver(dbClient, authService)
	resolver.SearchResolver = NewSearchResolver(dbClient, aiService)
	resolver.SubscriptionResolver = NewSubscriptionResolver(resolver)

	return resolver
}

// RecipeResolver handles recipe-related resolvers.
type RecipeResolver struct {
	DB *database.DGraphClient
}

// NewRecipeResolver creates a new recipe resolver.
func NewRecipeResolver(dbClient *database.DGraphClient) *RecipeResolver {
	return &RecipeResolver{
		DB: dbClient,
	}
}

// UserResolver handles user-related resolvers.
type UserResolver struct {
	DB          *database.DGraphClient
	AuthService *auth.Service
}

// NewUserResolver creates a new user resolver.
func NewUserResolver(dbClient *database.DGraphClient, authService *auth.Service) *UserResolver {
	return &UserResolver{
		DB:          dbClient,
		AuthService: authService,
	}
}
