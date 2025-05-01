//nolint:lll // Long function signatures are acceptable for GraphQL resolvers.
package resolvers

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/graphql"
	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	"github.com/carldunham/useful-cookery/internal/model"
)

// SearchResolver handles search-related resolvers.
type SearchResolver struct {
	DB        graphql.Database
	AIService *ai.Service
}

// NewSearchResolver creates a new search resolver.
func NewSearchResolver(db graphql.Database, aiService *ai.Service) *SearchResolver {
	return &SearchResolver{
		DB:        db,
		AIService: aiService,
	}
}

// SearchRecipes performs a natural language search for recipes.
//
//nolint:cyclop,funlen // This function is necessarily complex due to multiple search strategies and fallbacks
func (r *SearchResolver) SearchRecipes(ctx context.Context, query string, first *int, after *string) (*gqlmodel.RecipeConnection, error) {
	// Default values
	limitVal := 20
	offsetVal := 0

	if first != nil {
		limitVal = *first
	}

	// If after cursor is provided, decode it to get the offset
	if after != nil {
		var err error
		offsetVal, err = DecodeCursor(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		// Increment by 1 to get the next item after the cursor
		offsetVal++
	}

	if query == "" {
		return CreateRecipeConnection([]*model.Recipe{}, 0, limitVal, offsetVal), nil
	}

	// Start timing for the search
	start := time.Now()

	// Process the natural language query to extract structured parameters
	searchParams, err := r.AIService.ProcessNaturalLanguageQuery(ctx, query)
	if err != nil {
		log.Printf("Error processing natural language query: %v", err)
		// Fall back to basic text search if NLP fails
		recipes, searchErr := r.basicTextSearch(ctx, query, limitVal, offsetVal)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform basic text search after NLP error: %w", searchErr)
		}

		// Get total count for pagination info
		totalCount, err := r.countBasicTextSearch(ctx, query)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting basic text search results: %v", err)
			totalCount = len(recipes)
		}

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	// Generate embedding for vector search
	queryEmbedding, err := r.AIService.GenerateEmbedding(ctx, query)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		// Fall back to structured search if embedding fails
		recipes, searchErr := r.structuredSearch(ctx, searchParams, limitVal, offsetVal)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform structured search after embedding error: %w", searchErr)
		}

		// Get total count for pagination info
		totalCount, err := r.countStructuredSearch(ctx, searchParams)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting structured search results: %v", err)
			totalCount = len(recipes)
		}

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	// Hybrid search: combine vector search with structured filters
	results, err := r.hybridSearch(ctx, queryEmbedding, searchParams, limitVal, offsetVal)
	if err != nil {
		log.Printf("Error in hybrid search: %v", err)
		// Fall back to structured search
		recipes, searchErr := r.structuredSearch(ctx, searchParams, limitVal, offsetVal)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform structured search after hybrid search error: %w", searchErr)
		}

		// Get total count for pagination info
		totalCount, err := r.countStructuredSearch(ctx, searchParams)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting structured search results: %v", err)
			totalCount = len(recipes)
		}

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	// Log search performance
	elapsed := time.Since(start)
	log.Printf("Search completed in %s. Query: %s, Results: %d", elapsed, query, len(results))

	// Get total count for pagination info
	totalCount, err := r.countHybridSearch(ctx, queryEmbedding, searchParams)
	if err != nil {
		// If we can't get the total count, just use the length of the current page
		log.Printf("Error counting hybrid search results: %v", err)
		totalCount = len(results)
	}

	return CreateRecipeConnection(results, totalCount, limitVal, offsetVal), nil
}

// RecommendRecipes recommends recipes based on user preferences and available ingredients.
//
//nolint:cyclop,funlen // Complex function due to multiple recommendation strategies and error handling
func (r *SearchResolver) RecommendRecipes(
	ctx context.Context,
	userID *string,
	availableIngredients []string,
	first *int,
	after *string,
) (*gqlmodel.RecipeConnection, error) {
	// Default values
	limitVal := 20
	offsetVal := 0

	if first != nil {
		limitVal = *first
	}

	// If after cursor is provided, decode it to get the offset
	if after != nil {
		var err error
		offsetVal, err = DecodeCursor(*after)
		if err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		// Increment by 1 to get the next item after the cursor
		offsetVal++
	}

	// Start timing for the recommendation
	start := time.Now()

	// If no user ID provided and no ingredients, return popular recipes
	if userID == nil && len(availableIngredients) == 0 {
		recipes, err := r.getPopularRecipes(ctx, limitVal, offsetVal)
		if err != nil {
			return nil, fmt.Errorf("failed to get popular recipes: %w", err)
		}

		// Get total count for pagination info
		// Note: This is a simplification, as we don't have a direct way to count popular recipes
		// In a real implementation, you might want to add a CountPopularRecipes method
		totalCount, err := r.DB.CountRecipes(ctx, nil)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting recipes: %v", err)
			totalCount = len(recipes)
		}

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	var user *model.User
	var err error

	// Get user preferences if user ID is provided
	if userID != nil {
		user, err = r.DB.GetUser(ctx, *userID)
		if err != nil {
			log.Printf("Error getting user: %v", err)
		}
	}

	// If we have available ingredients, filter recipes by them
	if len(availableIngredients) > 0 {
		recipes, err := r.findRecipesByIngredients(ctx, availableIngredients, user, limitVal, offsetVal)
		if err != nil {
			log.Printf("Error finding recipes by ingredients: %v", err)
			return nil, fmt.Errorf("failed to find recipes by ingredients: %w", err)
		}

		// Get total count for pagination info
		totalCount, err := r.countRecipesByIngredients(ctx, availableIngredients, user)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting recipes by ingredients: %v", err)
			totalCount = len(recipes)
		}

		// Log recommendation performance
		elapsed := time.Since(start)
		log.Printf("Ingredient-based recommendation completed in %s. Results: %d", elapsed, len(recipes))

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	// If we have a user but no ingredients, recommend based on user preferences
	if user != nil {
		recipes, err := r.recommendBasedOnPreferences(ctx, user, limitVal, offsetVal)
		if err != nil {
			log.Printf("Error finding recipes by preferences: %v", err)
			return nil, fmt.Errorf("failed to find recipes by preferences: %w", err)
		}

		// Get total count for pagination info
		totalCount, err := r.countRecipesByPreferences(ctx, user)
		if err != nil {
			// If we can't get the total count, just use the length of the current page
			log.Printf("Error counting recipes by preferences: %v", err)
			totalCount = len(recipes)
		}

		// Log recommendation performance
		elapsed := time.Since(start)
		log.Printf("Preference-based recommendation completed in %s. Results: %d", elapsed, len(recipes))

		return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
	}

	// Fall back to popular recipes if nothing else works
	recipes, err := r.getPopularRecipes(ctx, limitVal, offsetVal)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular recipes: %w", err)
	}

	// Get total count for pagination info
	// Note: This is a simplification, as we don't have a direct way to count popular recipes
	// In a real implementation, you might want to add a CountPopularRecipes method
	totalCount, err := r.DB.CountRecipes(ctx, nil)
	if err != nil {
		// If we can't get the total count, just use the length of the current page
		log.Printf("Error counting recipes: %v", err)
		totalCount = len(recipes)
	}

	return CreateRecipeConnection(recipes, totalCount, limitVal, offsetVal), nil
}

// FindSubstitutes finds substitutes for an ingredient.
func (r *SearchResolver) FindSubstitutes(ctx context.Context, ingredientName string) ([]string, error) {
	substitutes, err := r.AIService.GenerateSubstitutes(ctx, ingredientName)
	if err != nil {
		return nil, fmt.Errorf("failed to generate substitutes: %w", err)
	}
	return substitutes, nil
}

// basicTextSearch performs a basic text search.
func (r *SearchResolver) basicTextSearch(ctx context.Context, searchText string, first int, offset int) ([]*model.Recipe, error) {
	// Use the database's GetRecipes method with a search filter
	filterMap := make(map[string]string)
	filterMap["search"] = searchText

	// Get recipes with the search filter
	recipes, err := r.DB.GetRecipes(ctx, filterMap, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute basic text search: %w", err)
	}

	return recipes, nil
}

// countBasicTextSearch counts recipes matching a basic text search.
func (r *SearchResolver) countBasicTextSearch(ctx context.Context, searchText string) (int, error) {
	// Use the database's CountRecipes method with a search filter
	filterMap := make(map[string]string)
	filterMap["search"] = searchText

	// Get count of recipes with the search filter
	count, err := r.DB.CountRecipes(ctx, filterMap)
	if err != nil {
		return 0, fmt.Errorf("failed to count basic text search results: %w", err)
	}

	return count, nil
}

// getPopularRecipes returns popular recipes.
func (r *SearchResolver) getPopularRecipes(ctx context.Context, first int, offset int) ([]*model.Recipe, error) {
	// Use the GetPopularRecipes method from the database implementation
	recipes, err := r.DB.GetPopularRecipes(ctx, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular recipes: %w", err)
	}
	return recipes, nil
}

// findRecipesByIngredients finds recipes that use the specified ingredients.
func (r *SearchResolver) findRecipesByIngredients(
	ctx context.Context,
	ingredients []string,
	user *model.User,
	first int,
	offset int,
) ([]*model.Recipe, error) {
	// Use the database's GetRecipes method with a filter map
	filterMap := make(map[string]string)

	// Join ingredients with commas for the filter
	if len(ingredients) > 0 {
		filterMap["ingredients"] = strings.Join(ingredients, ",")
	}

	// Add user preferences if available
	if user != nil && user.Preferences != nil {
		// Add dietary restrictions if any
		if len(user.Preferences.DietaryRestrictions) > 0 {
			filterMap["dietary_restrictions"] = strings.Join(user.Preferences.DietaryRestrictions, ",")
		}

		// Add disliked ingredients if any
		if len(user.Preferences.DislikedIngredients) > 0 {
			filterMap["excluded_ingredients"] = strings.Join(user.Preferences.DislikedIngredients, ",")
		}
	}

	// Get recipes with the filter
	recipes, err := r.DB.GetRecipes(ctx, filterMap, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by ingredients: %w", err)
	}

	return recipes, nil
}

// countRecipesByIngredients counts recipes that use the specified ingredients.
func (r *SearchResolver) countRecipesByIngredients(
	ctx context.Context,
	ingredients []string,
	user *model.User,
) (int, error) {
	// Use the database's CountRecipes method with a filter map
	filterMap := make(map[string]string)

	// Join ingredients with commas for the filter
	if len(ingredients) > 0 {
		filterMap["ingredients"] = strings.Join(ingredients, ",")
	}

	// Add user preferences if available
	if user != nil && user.Preferences != nil {
		// Add dietary restrictions if any
		if len(user.Preferences.DietaryRestrictions) > 0 {
			filterMap["dietary_restrictions"] = strings.Join(user.Preferences.DietaryRestrictions, ",")
		}

		// Add disliked ingredients if any
		if len(user.Preferences.DislikedIngredients) > 0 {
			filterMap["excluded_ingredients"] = strings.Join(user.Preferences.DislikedIngredients, ",")
		}
	}

	// Get count of recipes with the filter
	count, err := r.DB.CountRecipes(ctx, filterMap)
	if err != nil {
		return 0, fmt.Errorf("failed to count recipes by ingredients: %w", err)
	}

	return count, nil
}

// recommendBasedOnPreferences recommends recipes based on user preferences.
func (r *SearchResolver) recommendBasedOnPreferences(
	ctx context.Context,
	user *model.User,
	first int,
	offset int,
) ([]*model.Recipe, error) {
	if user.Preferences == nil {
		recipes, err := r.getPopularRecipes(ctx, first, offset)
		if err != nil {
			return nil, fmt.Errorf("failed to get popular recipes: %w", err)
		}
		return recipes, nil
	}

	// Use the database's GetRecipes method with a filter map
	filterMap := make(map[string]string)

	// Add favorite ingredients if any
	if len(user.Preferences.FavoriteIngredients) > 0 {
		filterMap["favorite_ingredients"] = strings.Join(user.Preferences.FavoriteIngredients, ",")
	}

	// Add cuisine preferences if any
	if len(user.Preferences.CuisinePreferences) > 0 {
		filterMap["cuisines"] = strings.Join(user.Preferences.CuisinePreferences, ",")
	}

	// Add skill level if set
	if user.Preferences.SkillLevel != "" {
		filterMap["difficulty"] = user.Preferences.SkillLevel
	}

	// Add disliked ingredients if any
	if len(user.Preferences.DislikedIngredients) > 0 {
		filterMap["excluded_ingredients"] = strings.Join(user.Preferences.DislikedIngredients, ",")
	}

	// Add dietary restrictions if any
	if len(user.Preferences.DietaryRestrictions) > 0 {
		filterMap["dietary_restrictions"] = strings.Join(user.Preferences.DietaryRestrictions, ",")
	}

	// Get recipes with the filter
	recipes, err := r.DB.GetRecipes(ctx, filterMap, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to find recipes by preferences: %w", err)
	}

	return recipes, nil
}

// countRecipesByPreferences counts recipes matching user preferences.
func (r *SearchResolver) countRecipesByPreferences(
	ctx context.Context,
	user *model.User,
) (int, error) {
	if user.Preferences == nil {
		// Fall back to counting popular recipes
		// Note: This is a simplification, as we don't have a direct way to count popular recipes
		// In a real implementation, you might want to add a CountPopularRecipes method
		count, err := r.DB.CountRecipes(ctx, nil)
		if err != nil {
			return 0, fmt.Errorf("failed to count recipes: %w", err)
		}
		return count, nil
	}

	// Use the database's CountRecipes method with a filter map
	filterMap := make(map[string]string)

	// Add favorite ingredients if any
	if len(user.Preferences.FavoriteIngredients) > 0 {
		filterMap["favorite_ingredients"] = strings.Join(user.Preferences.FavoriteIngredients, ",")
	}

	// Add cuisine preferences if any
	if len(user.Preferences.CuisinePreferences) > 0 {
		filterMap["cuisines"] = strings.Join(user.Preferences.CuisinePreferences, ",")
	}

	// Add skill level if set
	if user.Preferences.SkillLevel != "" {
		filterMap["difficulty"] = user.Preferences.SkillLevel
	}

	// Add disliked ingredients if any
	if len(user.Preferences.DislikedIngredients) > 0 {
		filterMap["excluded_ingredients"] = strings.Join(user.Preferences.DislikedIngredients, ",")
	}

	// Add dietary restrictions if any
	if len(user.Preferences.DietaryRestrictions) > 0 {
		filterMap["dietary_restrictions"] = strings.Join(user.Preferences.DietaryRestrictions, ",")
	}

	// Get count of recipes with the filter
	count, err := r.DB.CountRecipes(ctx, filterMap)
	if err != nil {
		return 0, fmt.Errorf("failed to count recipes by preferences: %w", err)
	}

	return count, nil
}

// structuredSearch performs a search using structured parameters.
func (r *SearchResolver) structuredSearch(ctx context.Context, params *model.SearchParams, first int, offset int) ([]*model.Recipe, error) {
	// Use the database's GetRecipes method with a filter map
	filterMap := make(map[string]string)

	// Add ingredients if any
	if len(params.Ingredients) > 0 {
		filterMap["ingredients"] = strings.Join(params.Ingredients, ",")
	}

	// Add excluded ingredients if any
	if len(params.ExcludedIngredients) > 0 {
		filterMap["excluded_ingredients"] = strings.Join(params.ExcludedIngredients, ",")
	}

	// Add categories if any
	if len(params.Categories) > 0 {
		filterMap["categories"] = strings.Join(params.Categories, ",")
	}

	// Add cuisine if set
	if params.Cuisine != "" {
		filterMap["cuisine"] = params.Cuisine
	}

	// Add dietary restrictions if any
	if len(params.DietaryRestrictions) > 0 {
		filterMap["dietary_restrictions"] = strings.Join(params.DietaryRestrictions, ",")
	}

	// Add max prep time if set
	if params.MaxPrepTime > 0 {
		filterMap["max_prep_time"] = strconv.Itoa(params.MaxPrepTime)
	}

	// Add difficulty if set
	if params.Difficulty != "" {
		filterMap["difficulty"] = params.Difficulty
	}

	// Get recipes with the filter
	recipes, err := r.DB.GetRecipes(ctx, filterMap, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to execute structured search: %w", err)
	}

	return recipes, nil
}

// countStructuredSearch counts recipes matching structured parameters.
func (r *SearchResolver) countStructuredSearch(ctx context.Context, params *model.SearchParams) (int, error) {
	// Use the database's CountRecipes method with a filter map
	filterMap := make(map[string]string)

	// Add ingredients if any
	if len(params.Ingredients) > 0 {
		filterMap["ingredients"] = strings.Join(params.Ingredients, ",")
	}

	// Add excluded ingredients if any
	if len(params.ExcludedIngredients) > 0 {
		filterMap["excluded_ingredients"] = strings.Join(params.ExcludedIngredients, ",")
	}

	// Add categories if any
	if len(params.Categories) > 0 {
		filterMap["categories"] = strings.Join(params.Categories, ",")
	}

	// Add cuisine if set
	if params.Cuisine != "" {
		filterMap["cuisine"] = params.Cuisine
	}

	// Add dietary restrictions if any
	if len(params.DietaryRestrictions) > 0 {
		filterMap["dietary_restrictions"] = strings.Join(params.DietaryRestrictions, ",")
	}

	// Add max prep time if set
	if params.MaxPrepTime > 0 {
		filterMap["max_prep_time"] = strconv.Itoa(params.MaxPrepTime)
	}

	// Add difficulty if set
	if params.Difficulty != "" {
		filterMap["difficulty"] = params.Difficulty
	}

	// Get count of recipes with the filter
	count, err := r.DB.CountRecipes(ctx, filterMap)
	if err != nil {
		return 0, fmt.Errorf("failed to count structured search results: %w", err)
	}

	return count, nil
}

// hybridSearch combines vector search with structured filters.
func (r *SearchResolver) hybridSearch(
	ctx context.Context,
	_ []float32, // TODO: use queryEmbedding or move to params.
	params *model.SearchParams,
	first int,
	offset int,
) ([]*model.Recipe, error) {
	// For PostgreSQL implementation, we'll use a simpler approach
	// that doesn't rely on vector search capabilities

	// Use the structured search as a fallback
	log.Printf("Vector search not implemented for PostgreSQL, falling back to structured search")
	recipes, err := r.structuredSearch(ctx, params, first, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to perform structured search: %w", err)
	}

	return recipes, nil
}

// countHybridSearch counts recipes matching hybrid search parameters.
func (r *SearchResolver) countHybridSearch(
	ctx context.Context,
	_ []float32, // TODO: use queryEmbedding or move to params.
	params *model.SearchParams,
) (int, error) {
	// For PostgreSQL implementation, we'll use a simpler approach
	// that doesn't rely on vector search capabilities

	// Use the structured search count as a fallback
	log.Printf("Vector search count not implemented for PostgreSQL, falling back to structured search count")
	count, err := r.countStructuredSearch(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("failed to count structured search results: %w", err)
	}

	return count, nil
}
