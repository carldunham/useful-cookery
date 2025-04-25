package resolvers

import (
	"errors"

	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/graphql"
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
	DB                   graphql.Database
	AIService            *ai.Service
	AuthService          *auth.Service
	RecipeResolver       *RecipeResolver
	UserResolver         *UserResolver
	SearchResolver       *SearchResolver
	SubscriptionResolver *SubscriptionResolver
}

// NewRootResolver creates a new root resolver.
func NewRootResolver(db graphql.Database, aiService *ai.Service, authService *auth.Service) *Resolver {
	// Create the root resolver first (needed for circular dependencies)
	resolver := &Resolver{
		DB:          db,
		AIService:   aiService,
		AuthService: authService,
	}

	// Create specialized resolvers
	resolver.RecipeResolver = NewRecipeResolver(db)
	resolver.UserResolver = NewUserResolver(db, authService)
	resolver.SearchResolver = NewSearchResolver(db, aiService)
	resolver.SubscriptionResolver = NewSubscriptionResolver(resolver)

	return resolver
}

// RecipeResolver handles recipe-related resolvers.
type RecipeResolver struct {
	DB graphql.Database
}

// NewRecipeResolver creates a new recipe resolver.
func NewRecipeResolver(db graphql.Database) *RecipeResolver {
	return &RecipeResolver{
		DB: db,
	}
}

// UserResolver handles user-related resolvers.
type UserResolver struct {
	DB          graphql.Database
	AuthService *auth.Service
}

// NewUserResolver creates a new user resolver.
func NewUserResolver(db graphql.Database, authService *auth.Service) *UserResolver {
	return &UserResolver{
		DB:          db,
		AuthService: authService,
	}
}
