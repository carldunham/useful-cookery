package resolvers

import (
	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver for GraphQL operations.
type Resolver struct {
	DB                   *database.DGraphClient
	AIService            *ai.AIService
	AuthService          *auth.Service
	RecipeResolver       *RecipeResolver
	UserResolver         *UserResolver
	SearchResolver       *SearchResolver
	SubscriptionResolver *SubscriptionResolver
}

// NewRootResolver creates a new root resolver.
func NewRootResolver(db *database.DGraphClient, aiService *ai.AIService, authService *auth.Service) *Resolver {
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
	DB *database.DGraphClient
}

// NewRecipeResolver creates a new recipe resolver.
func NewRecipeResolver(db *database.DGraphClient) *RecipeResolver {
	return &RecipeResolver{
		DB: db,
	}
}

// UserResolver handles user-related resolvers.
type UserResolver struct {
	DB          *database.DGraphClient
	AuthService *auth.Service
}

// NewUserResolver creates a new user resolver.
func NewUserResolver(db *database.DGraphClient, authService *auth.Service) *UserResolver {
	return &UserResolver{
		DB:          db,
		AuthService: authService,
	}
}
