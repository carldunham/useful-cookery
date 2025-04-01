package resolvers

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/carldunham/useful-cookery/internal/ai"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/models"
)

// SearchResolver handles search-related resolvers
type SearchResolver struct {
	DB        *database.DGraphClient
	AIService *ai.AIService
}

// NewSearchResolver creates a new search resolver
func NewSearchResolver(db *database.DGraphClient, aiService *ai.AIService) *SearchResolver {
	return &SearchResolver{
		DB:        db,
		AIService: aiService,
	}
}

// SearchRecipes performs a natural language search for recipes
func (r *SearchResolver) SearchRecipes(ctx context.Context, query string) ([]*models.Recipe, error) {
	if query == "" {
		return []*models.Recipe{}, nil
	}

	// Start timing for the search
	start := time.Now()

	// Process the natural language query to extract structured parameters
	searchParams, err := r.AIService.ProcessNaturalLanguageQuery(ctx, query)
	if err != nil {
		log.Printf("Error processing natural language query: %v", err)
		// Fall back to basic text search if NLP fails
		return r.basicTextSearch(ctx, query)
	}

	// Generate embedding for vector search
	queryEmbedding, err := r.AIService.GenerateEmbedding(ctx, query)
	if err != nil {
		log.Printf("Error generating embedding: %v", err)
		// Fall back to structured search if embedding fails
		return r.structuredSearch(ctx, searchParams)
	}

	// Hybrid search: combine vector search with structured filters
	results, err := r.hybridSearch(ctx, queryEmbedding, searchParams)
	if err != nil {
		log.Printf("Error in hybrid search: %v", err)
		// Fall back to structured search
		return r.structuredSearch(ctx, searchParams)
	}

	// Log search performance
	elapsed := time.Since(start)
	log.Printf("Search completed in %s. Query: %s, Results: %d", elapsed, query, len(results))

	return results, nil
}

// RecommendRecipes recommends recipes based on user preferences and available ingredients
func (r *SearchResolver) RecommendRecipes(
	ctx context.Context, 
	userID *string, 
	availableIngredients []string,
) ([]*models.Recipe, error) {
	// Start timing for the recommendation
	start := time.Now()

	// If no user ID provided and no ingredients, return popular recipes
	if userID == nil && len(availableIngredients) == 0 {
		return r.getPopularRecipes(ctx)
	}

	var user *models.User
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
			return nil, err
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
			return nil, err
		}

		// Log recommendation performance
		elapsed := time.Since(start)
		log.Printf("Preference-based recommendation completed in %s. Results: %d", elapsed, len(recipes))

		return recipes, nil
	}

	// Fall back to popular recipes if nothing else works
	return r.getPopularRecipes(ctx)
}

// FindSubstitutes finds substitutes for an ingredient
func (r *SearchResolver) FindSubstitutes(ctx context.Context, ingredientName string) ([]string, error) {
	return r.AIService.GenerateSubstitutes(ctx, ingredientName)
}

// basicTextSearch performs a basic text search
func (r *SearchResolver) basicTextSearch(ctx context.Context, query string) ([]*models.Recipe, error) {
	// Basic text search implementation
	q := `
	query SearchRecipes($searchText: string) {
		recipes(func: alloftext(title, $searchText)) {
			uid
			expand(_all_)
		}
	}`

	variables := map[string]string{
		"searchText": query,
	}

	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, q, variables, &result)
	if err != nil {
		return nil, err
	}
	
	return result.Recipes, nil
}

// getPopularRecipes returns popular recipes
func (r *SearchResolver) getPopularRecipes(ctx context.Context) ([]*models.Recipe, error) {
	q := `
	query PopularRecipes() {
		recipes(func: has(title), orderasc: likes, first: 20) {
			uid
			expand(_all_)
		}
	}`
	
	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}
	
	err := r.DB.Query(ctx, q, nil, &result)
	if err != nil {
		return nil, err
	}
	
	return result.Recipes, nil
}

// findRecipesByIngredients finds recipes that use the specified ingredients
func (r *SearchResolver) findRecipesByIngredients(
	ctx context.Context,
	ingredients []string,
	user *models.User,
) ([]*models.Recipe, error) {
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
	
	q := fmt.Sprintf(`
	query RecipesByIngredients(%s) {
		recipes(func: has(title), @filter(%s), first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), conditions)
	
	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}
	
	err := r.DB.Query(ctx, q, variables, &result)
	if err != nil {
		return nil, err
	}
	
	return result.Recipes, nil
}

// recommendBasedOnPreferences recommends recipes based on user preferences
func (r *SearchResolver) recommendBasedOnPreferences(
	ctx context.Context,
	user *models.User,
) ([]*models.Recipe, error) {
	if user.Preferences == nil {
		return r.getPopularRecipes(ctx)
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
	
	q := fmt.Sprintf(`
	query RecommendRecipes(%s) {
		recipes(func: has(title), %s, first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), filterClause)
	
	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}
	
	err := r.DB.Query(ctx, q, variables, &result)
	if err != nil {
		return nil, err
	}
	
	return result.Recipes, nil
}

// buildVariableDeclarations creates the variable declaration string for DGraph queries
func buildVariableDeclarations(vars map[string]string) string {
	if len(vars) == 0 {
		return ""
	}
	
	declarations := []string{}
	for k := range vars {
		declarations = append(declarations, fmt.Sprintf("$%s: string", k))
	}
	
	return strings.Join(declarations, ", ")
}ctx, q, variables, &result)
	if err != nil {
		return nil, err
	}

	return result.Recipes, nil
}

// structuredSearch performs a search using structured parameters
func (r *SearchResolver) structuredSearch(ctx context.Context, params *models.SearchParams) ([]*models.Recipe, error) {
	// Build query conditions based on search parameters
	var conditions []string
	variables := make(map[string]string)

	// Add conditions for each parameter
	if len(params.Ingredients) > 0 {
		ingredientCondition := "has(ingredients)"
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
		categoryCondition := "has(categories)"
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
		variables["maxPrepTime"] = fmt.Sprintf("%d", params.MaxPrepTime)
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

	q := fmt.Sprintf(`
	query StructuredSearch(%s) {
		recipes(func: has(title), %s, first: 20) {
			uid
			expand(_all_)
		}
	}`, buildVariableDeclarations(variables), filterClause)

	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}

	err := r.DB.Query(ctx, q, variables, &result)
	if err != nil {
		return nil, err
	}

	return result.Recipes, nil
}

// hybridSearch combines vector search with structured filters
func (r *SearchResolver) hybridSearch(
	ctx context.Context,
	embedding []float32,
	params *models.SearchParams,
) ([]*models.Recipe, error) {
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
	
	// Add structured conditions...
	// (similar to structuredSearch method)
	
	// Construct DGraph query with vector search
	q := `
	query HybridSearch($embeddings: string, $distance: float) {
		recipes(func: vector(embeddings, $embeddings, $distance), first: 20) {
			uid
			expand(_all_)
		}
	}`
	
	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}
	
	err := r.DB.Query(