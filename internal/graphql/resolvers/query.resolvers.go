package resolvers

import (
	"context"
	"errors"
	"fmt"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// Recipe returns the recipe with the given ID.
func (r *Resolver) recipe(ctx context.Context, id string) (*domainmodel.Recipe, error) {
	return r.RecipeResolver.getRecipe(ctx, id)
}

// Recipes returns a list of recipes based on the given filters.
func (r *Resolver) recipes(
	ctx context.Context,
	filter *gqlmodel.RecipeFilter,
	order *gqlmodel.RecipeOrder,
	first *int,
	offset *int,
) ([]*domainmodel.Recipe, error) {
	return r.RecipeResolver.getRecipes(ctx, filter, order, first, offset)
}

// GetRecipe returns a recipe by ID.
func (r *RecipeResolver) getRecipe(ctx context.Context, id string) (*domainmodel.Recipe, error) {
	recipe, err := r.DB.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, ErrRecipeNotFound
		}
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}
	return recipe, nil
}

// GetRecipes returns recipes based on filters.
func (r *RecipeResolver) getRecipes(
	ctx context.Context,
	filter *gqlmodel.RecipeFilter,
	_ *gqlmodel.RecipeOrder,
	first *int,
	offset *int,
) ([]*domainmodel.Recipe, error) {
	// Default values
	limitVal := 20
	offsetVal := 0

	if first != nil {
		limitVal = *first
	}
	if offset != nil {
		offsetVal = *offset
	}

	// Convert filter to map for database query
	filterMap := make(map[string]string)
	if filter != nil {
		if filter.Search != nil {
			filterMap["title"] = *filter.Search
		}
		if filter.Cuisine != nil {
			filterMap["cuisine"] = *filter.Cuisine
		}
		if filter.AuthorID != nil {
			filterMap["authorID"] = *filter.AuthorID
		}
		// Other filters would be handled similarly
	}

	recipes, err := r.DB.GetRecipes(ctx, filterMap, limitVal, offsetVal)
	if err != nil {
		return nil, fmt.Errorf("failed to get recipes: %w", err)
	}

	return recipes, nil
}

// Me returns the currently authenticated user.
func (r *Resolver) me(ctx context.Context) (*domainmodel.User, error) {
	return r.UserResolver.getCurrentUser(ctx)
}

// User returns a user by ID.
func (r *Resolver) user(ctx context.Context, id string) (*domainmodel.User, error) {
	return r.UserResolver.getUser(ctx, id)
}

// GetCurrentUser returns the authenticated user from context.
func (r *UserResolver) getCurrentUser(ctx context.Context) (*domainmodel.User, error) {
	user := auth.GetUserFromContext(ctx)
	if user == nil {
		return nil, ErrNotAuthenticated
	}
	return user, nil
}

// GetUser returns a user by ID.
func (r *UserResolver) getUser(ctx context.Context, userID string) (*domainmodel.User, error) {
	// Check if the requester has permission to view the user
	currentUser := auth.GetUserFromContext(ctx)
	currentRole := auth.GetRoleFromContext(ctx)

	// Allow admins to view any user, otherwise users can only view themselves
	if currentUser == nil || (currentRole != string(domainmodel.AdminRole) && currentUser.ID != userID) {
		return nil, ErrPermissionDenied
	}

	user, err := r.DB.GetUser(ctx, userID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

// Categories returns all recipe categories.
func (r *Resolver) categories(_ context.Context) ([]*domainmodel.Category, error) {
	// Implementation would query the database for all categories
	return []*domainmodel.Category{}, nil
}

// Category returns a category by ID.
func (r *Resolver) category(_ context.Context, _ string) (*domainmodel.Category, error) {
	// Implementation would query the database for the category
	return &domainmodel.Category{}, nil
}

// SearchRecipes performs a natural language search for recipes.
func (r *Resolver) searchRecipes(ctx context.Context, query string) ([]*domainmodel.Recipe, error) {
	return r.SearchResolver.SearchRecipes(ctx, query)
}

// RecommendRecipes recommends recipes based on user preferences and ingredients.
func (r *Resolver) recommendRecipes(
	ctx context.Context,
	userID *string,
	availableIngredients []string,
) ([]*domainmodel.Recipe, error) {
	return r.SearchResolver.RecommendRecipes(ctx, userID, availableIngredients)
}

// FindSubstitutes finds substitutes for an ingredient.
func (r *Resolver) findSubstitutes(ctx context.Context, ingredientName string) ([]string, error) {
	return r.SearchResolver.FindSubstitutes(ctx, ingredientName)
}
