package resolvers

import (
	"context"
	"errors"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/graphql/models"
	domainModels "github.com/carldunham/useful-cookery/internal/model"
)

// Recipe returns the recipe with the given ID
func (r *Resolver) Recipe(ctx context.Context, id string) (*domainModels.Recipe, error) {
	return r.RecipeResolver.GetRecipe(ctx, id)
}

// Recipes returns a list of recipes based on the given filters
func (r *Resolver) Recipes(
	ctx context.Context,
	filter *models.RecipeFilter,
	order *models.RecipeOrder,
	first *int,
	offset *int,
) ([]*domainModels.Recipe, error) {
	return r.RecipeResolver.GetRecipes(ctx, filter, order, first, offset)
}

// GetRecipe returns a recipe by ID
func (r *RecipeResolver) GetRecipe(ctx context.Context, id string) (*domainModels.Recipe, error) {
	recipe, err := r.DB.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("recipe not found")
		}
		return nil, err
	}
	return recipe, nil
}

// GetRecipes returns recipes based on filters
func (r *RecipeResolver) GetRecipes(
	ctx context.Context,
	filter *models.RecipeFilter,
	order *models.RecipeOrder,
	first *int,
	offset *int,
) ([]*domainModels.Recipe, error) {
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
			filterMap["authorId"] = *filter.AuthorID
		}
		// Other filters would be handled similarly
	}

	recipes, err := r.DB.GetRecipes(ctx, filterMap, limitVal, offsetVal)
	if err != nil {
		return nil, err
	}

	return recipes, nil
}

// Me returns the currently authenticated user
func (r *Resolver) Me(ctx context.Context) (*domainModels.User, error) {
	return r.UserResolver.GetCurrentUser(ctx)
}

// User returns a user by ID
func (r *Resolver) User(ctx context.Context, id string) (*domainModels.User, error) {
	return r.UserResolver.GetUser(ctx, id)
}

// GetCurrentUser returns the authenticated user from context
func (r *UserResolver) GetCurrentUser(ctx context.Context) (*domainModels.User, error) {
	user := auth.GetUserFromContext(ctx)
	if user == nil {
		return nil, errors.New("not authenticated")
	}
	return user, nil
}

// GetUser returns a user by ID
func (r *UserResolver) GetUser(ctx context.Context, id string) (*domainModels.User, error) {
	// Check if the requester has permission to view the user
	currentUser := auth.GetUserFromContext(ctx)
	currentRole := auth.GetRoleFromContext(ctx)

	// Allow admins to view any user, otherwise users can only view themselves
	if currentUser == nil || (currentRole != "ADMIN" && currentUser.ID != id) {
		return nil, errors.New("permission denied")
	}

	user, err := r.DB.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// Categories returns all recipe categories
func (r *Resolver) Categories(ctx context.Context) ([]*domainModels.Category, error) {
	// Implementation would query the database for all categories
	return []*domainModels.Category{}, nil
}

// Category returns a category by ID
func (r *Resolver) Category(ctx context.Context, id string) (*domainModels.Category, error) {
	// Implementation would query the database for the category
	return &domainModels.Category{}, nil
}

// SearchRecipes performs a natural language search for recipes
func (r *Resolver) SearchRecipes(ctx context.Context, query string) ([]*domainModels.Recipe, error) {
	return r.SearchResolver.SearchRecipes(ctx, query)
}

// RecommendRecipes recommends recipes based on user preferences and ingredients
func (r *Resolver) RecommendRecipes(
	ctx context.Context,
	userID *string,
	availableIngredients []string,
) ([]*domainModels.Recipe, error) {
	return r.SearchResolver.RecommendRecipes(ctx, userID, availableIngredients)
}

// FindSubstitutes finds substitutes for an ingredient
func (r *Resolver) FindSubstitutes(ctx context.Context, ingredientName string) ([]string, error) {
	return r.SearchResolver.FindSubstitutes(ctx, ingredientName)
}
