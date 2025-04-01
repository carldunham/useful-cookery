package resolvers

import (
	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver for GraphQL operations
type Resolver struct {
	DB             *database.DGraphClient
	AIService      *ai.AIService
	AuthService    *auth.Service
	RecipeResolver *RecipeResolver
	UserResolver   *UserResolver
	SearchResolver *SearchResolver
}

// NewRootResolver creates a new root resolver
func NewRootResolver(db *database.DGraphClient, aiService *ai.AIService, authService *auth.Service) *Resolver {
	// Create specialized resolvers
	recipeResolver := NewRecipeResolver(db)
	userResolver := NewUserResolver(db, authService)
	searchResolver := NewSearchResolver(db, aiService)

	return &Resolver{
		DB:             db,
		AIService:      aiService,
		AuthService:    authService,
		RecipeResolver: recipeResolver,
		UserResolver:   userResolver,
		SearchResolver: searchResolver,
	}
}

// RecipeResolver handles recipe-related resolvers
type RecipeResolver struct {
	DB *database.DGraphClient
}

// NewRecipeResolver creates a new recipe resolver
func NewRecipeResolver(db *database.DGraphClient) *RecipeResolver {
	return &RecipeResolver{
		DB: db,
	}
}

// UserResolver handles user-related resolvers
type UserResolver struct {
	DB          *database.DGraphClient
	AuthService *auth.Service
}

// NewUserResolver creates a new user resolver
func NewUserResolver(db *database.DGraphClient, authService *auth.Service) *UserResolver {
	return &UserResolver{
		DB:          db,
		AuthService: authService,
	}
}
