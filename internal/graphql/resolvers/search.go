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
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// Constants for search queries.
const (
	typeRecipe = "type(Recipe)"
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
func (r *SearchResolver) SearchRecipes(ctx context.Context, query string) ([]*domainmodel.Recipe, error) {
	if query == "" {
		return []*domainmodel.Recipe{}, nil
	}

	// Start timing for the search
	start := time.Now()

	// Process the natural language query to extract structured parameters
	searchParams, err := r.AIService.ProcessNaturalLanguageQuery(ctx, query)
	if err != nil {
		log.Printf("Error processing natural language query: %v", err)
		// Fall back to basic text search if NLP fails
		recipes, searchErr := r.basicTextSearch(ctx, query)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform basic text search after NLP error: %w", searchErr)
		}
		return recipes, nil
	}

	// Generate embedding for vector search
	queryEmbedding, err := r.AIService.GenerateEmbedding(ctx, query)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		// Fall back to structured search if embedding fails
		recipes, searchErr := r.structuredSearch(ctx, searchParams)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform structured search after embedding error: %w", searchErr)
		}
		return recipes, nil
	}

	// Hybrid search: combine vector search with structured filters
	results, err := r.hybridSearch(ctx, queryEmbedding, searchParams)
	if err != nil {
		log.Printf("Error in hybrid search: %v", err)
		// Fall back to structured search
		recipes, searchErr := r.structuredSearch(ctx, searchParams)
		if searchErr != nil {
			return nil, fmt.Errorf("failed to perform structured search after hybrid search error: %w", searchErr)
		}
		return recipes, nil
	}

	// Log search performance
	elapsed := time.Since(start)
	log.Printf("Search completed in %s. Query: %s, Results: %d", elapsed, query, len(results))

	return results, nil
}

// RecommendRecipes recommends recipes based on user preferences and available ingredients.
//
//nolint:cyclop // Complex function due to multiple recommendation strategies.
func (r *SearchResolver) RecommendRecipes(
	ctx context.Context,
	userID *string,
	availableIngredients []string,
) ([]*domainmodel.Recipe, error) {
	// Start timing for the recommendation
	start := time.Now()

	// If no user ID provided and no ingredients, return popular recipes
	if userID == nil && len(availableIngredients) == 0 {
		recipes, err := r.getPopularRecipes(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get popular recipes: %w", err)
		}
		return recipes, nil
	}

	var user *domainmodel.User
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
		recipes, err := r.findRecipesByIngredients(ctx, availableIngredients, user)
		if err != nil {
			log.Printf("Error finding recipes by ingredients: %v", err)
			return nil, fmt.Errorf("failed to find recipes by ingredients: %w", err)
		}

		// Log recommendation performance
		elapsed := time.Since(start)
		log.Printf("Ingredient-based recommendation completed in %s. Results: %d", elapsed, len(recipes))

		return recipes, nil
	}

	// If we have a user but no ingredients, recommend based on user preferences
	if user != nil {
		recipes, err := r.recommendBasedOnPreferences(ctx, user)
		if err != nil {
			log.Printf("Error finding recipes by preferences: %v", err)
			return nil, fmt.Errorf("failed to find recipes by preferences: %w", err)
		}

		// Log recommendation performance
		elapsed := time.Since(start)
		log.Printf("Preference-based recommendation completed in %s. Results: %d", elapsed, len(recipes))

		return recipes, nil
	}

	// Fall back to popular recipes if nothing else works
	recipes, err := r.getPopularRecipes(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular recipes: %w", err)
	}
	return recipes, nil
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
func (r *SearchResolver) basicTextSearch(ctx context.Context, searchText string) ([]*domainmodel.Recipe, error) {
	// Basic text search implementation
	queryStr := `
	query SearchRecipes($searchText: string) {
		recipes(func: alloftext(title, $searchText)) {
			uid
			expand(_all_)
		}
	}`

	variables := map[string]string{
		"searchText": searchText,
	}

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute basic text search query: %w", err)
	}

	return result.Recipes, nil
}

// getPopularRecipes returns popular recipes.
func (r *SearchResolver) getPopularRecipes(ctx context.Context) ([]*domainmodel.Recipe, error) {
	queryStr := `
	query PopularRecipes() {
		recipes(func: type(Recipe), orderasc: likes, first: 20) {
			uid
			expand(_all_)
		}
	}`

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute popular recipes query: %w", err)
	}

	return result.Recipes, nil
}

// findRecipesByIngredients finds recipes that use the specified ingredients.
func (r *SearchResolver) findRecipesByIngredients(
	ctx context.Context,
	ingredients []string,
	user *domainmodel.User,
) ([]*domainmodel.Recipe, error) {
	variables := make(map[string]string)

	// Build ingredient match conditions
	ingConditions := []string{}
	for i, ing := range ingredients {
		varName := fmt.Sprintf("ing%d", i)
		ingConditions = append(ingConditions, fmt.Sprintf("anyoftext(ingredients, $%s)", varName))
		variables[varName] = ing
	}

	// Add user preferences if available
	userFilters := []string{}
	if user != nil && user.Preferences != nil {
		// Exclude disliked ingredients
		for i, ing := range user.Preferences.DislikedIngredients {
			varName := fmt.Sprintf("disliked%d", i)
			userFilters = append(userFilters, fmt.Sprintf("NOT anyoftext(ingredients, $%s)", varName))
			variables[varName] = ing
		}

		// Match dietary restrictions
		for i, diet := range user.Preferences.DietaryRestrictions {
			varName := fmt.Sprintf("diet%d", i)
			userFilters = append(userFilters, fmt.Sprintf("anyoftext(tags, $%s)", varName))
			variables[varName] = diet
		}
	}

	// Combine all conditions
	conditions := fmt.Sprintf("(%s)", strings.Join(ingConditions, " OR "))
	if len(userFilters) > 0 {
		for _, filter := range userFilters {
			conditions = fmt.Sprintf("%s AND %s", conditions, filter)
		}
	}

	queryStr := fmt.Sprintf(`
	query RecipesByIngredients(%s) {
		recipes(func: type(Recipe), @filter(%s), first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), conditions)

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute recipes by ingredients query: %w", err)
	}

	return result.Recipes, nil
}

// recommendBasedOnPreferences recommends recipes based on user preferences.
//
//nolint:cyclop,funlen // Complex function due to handling multiple user preference fields.
func (r *SearchResolver) recommendBasedOnPreferences(
	ctx context.Context,
	user *domainmodel.User,
) ([]*domainmodel.Recipe, error) {
	if user.Preferences == nil {
		recipes, err := r.getPopularRecipes(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get popular recipes: %w", err)
		}
		return recipes, nil
	}

	variables := make(map[string]string)
	conditions := []string{}

	// Match favorite ingredients
	favIngredients := []string{}
	for i, ing := range user.Preferences.FavoriteIngredients {
		varName := fmt.Sprintf("fav%d", i)
		favIngredients = append(favIngredients, fmt.Sprintf("anyoftext(ingredients, $%s)", varName))
		variables[varName] = ing
	}

	if len(favIngredients) > 0 {
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(favIngredients, " OR ")))
	}

	// Match favorite cuisines
	favCuisines := []string{}
	for i, cuisine := range user.Preferences.CuisinePreferences {
		varName := fmt.Sprintf("cuisine%d", i)
		favCuisines = append(favCuisines, fmt.Sprintf("eq(cuisine, $%s)", varName))
		variables[varName] = cuisine
	}

	if len(favCuisines) > 0 {
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(favCuisines, " OR ")))
	}

	// Match skill level
	if user.Preferences.SkillLevel != "" {
		conditions = append(conditions, "eq(difficulty, $skillLevel)")
		variables["skillLevel"] = user.Preferences.SkillLevel
	}

	// Exclude disliked ingredients
	for i, ing := range user.Preferences.DislikedIngredients {
		varName := fmt.Sprintf("disliked%d", i)
		conditions = append(conditions, fmt.Sprintf("NOT anyoftext(ingredients, $%s)", varName))
		variables[varName] = ing
	}

	// Match dietary restrictions
	for i, diet := range user.Preferences.DietaryRestrictions {
		varName := fmt.Sprintf("diet%d", i)
		conditions = append(conditions, fmt.Sprintf("anyoftext(tags, $%s)", varName))
		variables[varName] = diet
	}

	// Add conditions to query
	var filterClause string
	if len(conditions) > 0 {
		filterClause = fmt.Sprintf("@filter(%s)", strings.Join(conditions, " AND "))
	}

	queryStr := fmt.Sprintf(`
	query RecommendRecipes(%s) {
		recipes(func: type(Recipe), %s, first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), filterClause)

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute recommendation by preferences query: %w", err)
	}

	return result.Recipes, nil
}

// buildVariableDeclarations creates the variable declaration string for DGraph queries.
func buildVariableDeclarations(vars map[string]string) string {
	if len(vars) == 0 {
		return ""
	}

	declarations := []string{}
	for k := range vars {
		declarations = append(declarations, fmt.Sprintf("$%s: string", k))
	}

	return strings.Join(declarations, ", ")
}

// structuredSearch performs a search using structured parameters.
//
//nolint:cyclop,funlen // Complex function due to handling multiple search parameters.
func (r *SearchResolver) structuredSearch(ctx context.Context, params *domainmodel.SearchParams) ([]*domainmodel.Recipe, error) {
	// Build query conditions based on search parameters
	var conditions []string
	variables := make(map[string]string)

	// Add conditions for each parameter
	if len(params.Ingredients) > 0 {
		ingredientCondition := typeRecipe
		for i, ing := range params.Ingredients {
			varName := fmt.Sprintf("ing%d", i)
			ingredientCondition = fmt.Sprintf("%s AND anyoftext(ingredients, $%s)", ingredientCondition, varName)
			variables[varName] = ing
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", ingredientCondition))
	}

	// Add excluded ingredients
	if len(params.ExcludedIngredients) > 0 {
		for i, ing := range params.ExcludedIngredients {
			varName := fmt.Sprintf("excing%d", i)
			conditions = append(conditions, fmt.Sprintf("NOT anyoftext(ingredients, $%s)", varName))
			variables[varName] = ing
		}
	}

	// Add categories
	if len(params.Categories) > 0 {
		categoryCondition := typeRecipe
		for i, cat := range params.Categories {
			varName := fmt.Sprintf("cat%d", i)
			categoryCondition = fmt.Sprintf("%s AND anyoftext(categories, $%s)", categoryCondition, varName)
			variables[varName] = cat
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", categoryCondition))
	}

	// Add cuisine
	if params.Cuisine != "" {
		conditions = append(conditions, "anyoftext(cuisine, $cuisine)")
		variables["cuisine"] = params.Cuisine
	}

	// Add dietary restrictions
	if len(params.DietaryRestrictions) > 0 {
		for i, diet := range params.DietaryRestrictions {
			varName := fmt.Sprintf("diet%d", i)
			conditions = append(conditions, fmt.Sprintf("anyoftext(tags, $%s)", varName))
			variables[varName] = diet
		}
	}

	// Add max prep time
	if params.MaxPrepTime > 0 {
		conditions = append(conditions, "le(prepTime, $maxPrepTime)")
		variables["maxPrepTime"] = strconv.Itoa(params.MaxPrepTime)
	}

	// Add difficulty
	if params.Difficulty != "" {
		conditions = append(conditions, "eq(difficulty, $difficulty)")
		variables["difficulty"] = params.Difficulty
	}

	// Construct query
	var filterClause string
	if len(conditions) > 0 {
		filterClause = fmt.Sprintf("@filter(%s)", conditions[0])
		for i := 1; i < len(conditions); i++ {
			filterClause = fmt.Sprintf("%s AND %s", filterClause, conditions[i])
		}
	}

	queryStr := fmt.Sprintf(`
	query StructuredSearch(%s) {
		recipes(func: type(Recipe), %s, first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), filterClause)

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute structured search query: %w", err)
	}

	return result.Recipes, nil
}

// hybridSearch combines vector search with structured filters.
//
//nolint:cyclop,funlen // Complex function due to handling multiple search parameters and vector search.
func (r *SearchResolver) hybridSearch(
	ctx context.Context,
	embedding []float32,
	params *domainmodel.SearchParams,
) ([]*domainmodel.Recipe, error) {
	// Implement vector search with filters
	// This is a simplified example - actual implementation would depend on DGraph's vector search capabilities

	// Convert embedding to string for query
	embeddingStr := "["
	for i, val := range embedding {
		if i > 0 {
			embeddingStr += ", "
		}
		embeddingStr += fmt.Sprintf("%f", val)
	}
	embeddingStr += "]"

	// Build conditions from search params
	var conditions []string
	variables := map[string]string{
		"embeddings": embeddingStr,
		"distance":   "0.8", // Threshold for vector similarity
	}

	// Add conditions for each parameter
	if len(params.Ingredients) > 0 {
		ingredientCondition := typeRecipe
		for i, ing := range params.Ingredients {
			varName := fmt.Sprintf("ing%d", i)
			ingredientCondition = fmt.Sprintf("%s AND anyoftext(ingredients, $%s)", ingredientCondition, varName)
			variables[varName] = ing
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", ingredientCondition))
	}

	// Add excluded ingredients
	if len(params.ExcludedIngredients) > 0 {
		for i, ing := range params.ExcludedIngredients {
			varName := fmt.Sprintf("excing%d", i)
			conditions = append(conditions, fmt.Sprintf("NOT anyoftext(ingredients, $%s)", varName))
			variables[varName] = ing
		}
	}

	// Add categories
	if len(params.Categories) > 0 {
		categoryCondition := typeRecipe
		for i, cat := range params.Categories {
			varName := fmt.Sprintf("cat%d", i)
			categoryCondition = fmt.Sprintf("%s AND anyoftext(categories, $%s)", categoryCondition, varName)
			variables[varName] = cat
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", categoryCondition))
	}

	// Add cuisine
	if params.Cuisine != "" {
		conditions = append(conditions, "anyoftext(cuisine, $cuisine)")
		variables["cuisine"] = params.Cuisine
	}

	// Add dietary restrictions
	if len(params.DietaryRestrictions) > 0 {
		for i, diet := range params.DietaryRestrictions {
			varName := fmt.Sprintf("diet%d", i)
			conditions = append(conditions, fmt.Sprintf("anyoftext(tags, $%s)", varName))
			variables[varName] = diet
		}
	}

	// Add max prep time
	if params.MaxPrepTime > 0 {
		conditions = append(conditions, "le(prepTime, $maxPrepTime)")
		variables["maxPrepTime"] = strconv.Itoa(params.MaxPrepTime)
	}

	// Add difficulty
	if params.Difficulty != "" {
		conditions = append(conditions, "eq(difficulty, $difficulty)")
		variables["difficulty"] = params.Difficulty
	}

	// Construct filter clause if we have conditions
	var filterClause string
	if len(conditions) > 0 {
		filterClause = fmt.Sprintf("@filter(%s)", strings.Join(conditions, " AND "))
	}

	// Construct DGraph query with vector search and filters
	queryStr := fmt.Sprintf(`
	query HybridSearch(%s) {
		recipes(func: vector(embeddings, $embeddings, $distance), %s, first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), filterClause)

	var result struct {
		Recipes []*domainmodel.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, queryStr, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to execute hybrid search query: %w", err)
	}

	// If we got no results from vector search, fall back to structured search
	if len(result.Recipes) == 0 {
		log.Printf("No results from vector search, falling back to structured search")
		recipes, err := r.structuredSearch(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("failed to perform structured search after empty vector search: %w", err)
		}
		return recipes, nil
	}

	return result.Recipes, nil
}
