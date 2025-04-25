package database

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/dgraph-io/dgo/v240"
	"github.com/dgraph-io/dgo/v240/protos/api"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/model"
)

// DGraphDatabase implements database operations using DGraph.
type DGraphDatabase struct {
	client *dgo.Dgraph
	cache  dbtypes.Cache
}

// NewDGraphDatabase creates a new DGraph database instance.
func NewDGraphDatabase(options dbtypes.DatabaseOptions) (*DGraphDatabase, error) {
	dgraphClient, err := dgo.Open(options.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("creating client: %w", err)
	}

	db := &DGraphDatabase{
		client: dgraphClient,
	}

	// Set up cache if enabled
	if options.CacheEnabled && options.Cache != nil {
		db.cache = options.Cache
	}

	return db, nil
}

// Query executes a GraphQL+ query against DGraph.
func (db *DGraphDatabase) Query(ctx context.Context, query string, vars map[string]string, result any) error {
	txn := db.client.NewTxn()
	defer func() {
		if err := txn.Discard(ctx); err != nil {
			log.Printf("error discarding Query transaction: %v", err)
		}
	}()

	// Create request
	req := &api.Request{
		Query: query,
		Vars:  vars,
	}

	// Execute query
	resp, err := txn.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	// Unmarshal result
	err = json.Unmarshal(resp.GetJson(), result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Mutate executes a mutation against DGraph.
func (db *DGraphDatabase) Mutate(ctx context.Context, data any) (*api.Response, error) {
	txn := db.client.NewTxn()
	defer func() {
		if err := txn.Discard(ctx); err != nil {
			log.Printf("error discarding Mutation transaction: %v", err)
		}
	}()

	// Marshal data
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal data: %w", err)
	}

	// Execute mutation
	mu := &api.Mutation{
		SetJson:   jsonData,
		CommitNow: true,
	}

	resp, err := txn.Mutate(ctx, mu)
	if err != nil {
		return nil, fmt.Errorf("failed to execute mutation: %w", err)
	}

	return resp, nil
}

// GetUser fetches a user by ID.
func (db *DGraphDatabase) GetUser(ctx context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	userQuery := `
	query GetUser($id: string) {
		user(func: uid($id)) {
			uid
			name
			email
			role
			preferences {
				dietaryRestrictions
				favoriteIngredients
				dislikedIngredients
				skillLevel
				cuisinePreferences
			}
			createdAt
			updatedAt
		}
	}`

	vars := map[string]string{
		"$id": userID,
	}

	var result struct {
		Users []*model.User `json:"user"`
	}

	err := db.Query(ctx, userQuery, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query user by ID: %w", err)
	}

	if len(result.Users) == 0 {
		return nil, dbtypes.ErrNotFound
	}

	return result.Users[0], nil
}

// GetUserByEmail fetches a user by email.
func (db *DGraphDatabase) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, dbtypes.ErrEmailRequired
	}

	// TODO: use type(User) instead of has(email).
	emailQuery := `
	{
		user(func: has(email)) {
			uid
			name
			email
			password
			role
			preferences {
				dietaryRestrictions
				favoriteIngredients
				dislikedIngredients
				skillLevel
				cuisinePreferences
			}
			createdAt
			updatedAt
		}
	}`

	var result struct {
		Users []*model.User `json:"user"`
	}

	err := db.Query(ctx, emailQuery, nil, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	// Filter users by email in code
	for _, user := range result.Users {
		if user.Email == email {
			return user, nil
		}
	}

	return nil, dbtypes.ErrNotFound
}

// CreateUser creates a new user.
func (db *DGraphDatabase) CreateUser(ctx context.Context, user *model.User) error {
	// Check if user already exists
	existingUser, err := db.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return dbtypes.ErrAlreadyExists
	}

	// Create user
	_, err = db.Mutate(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// GetUsers fetches a list of users with pagination.
func (db *DGraphDatabase) GetUsers(ctx context.Context, limit, offset int) ([]*model.User, error) {
	query := `
	query GetUsers($limit: int, $offset: int) {
		users(func: has(email), first: $limit, offset: $offset) {
			uid
			name
			email
			role
			createdAt
			updatedAt
		}
	}`

	vars := map[string]string{
		"$limit":  strconv.Itoa(limit),
		"$offset": strconv.Itoa(offset),
	}

	var result struct {
		Users []*model.User `json:"users"`
	}

	err := db.Query(ctx, query, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}

	return result.Users, nil
}

// UpdateUser updates a user.
func (db *DGraphDatabase) UpdateUser(ctx context.Context, user *model.User) error {
	_, err := db.Mutate(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// GetCategory fetches a category by ID.
func (db *DGraphDatabase) GetCategory(ctx context.Context, categoryID string) (*model.Category, error) {
	if categoryID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	categoryQuery := `
	query GetCategory($id: string) {
		category(func: uid($id)) {
			uid
			name
			description
			createdAt
			updatedAt
		}
	}`

	vars := map[string]string{
		"$id": categoryID,
	}

	var result struct {
		Categories []*model.Category `json:"category"`
	}

	err := db.Query(ctx, categoryQuery, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query category by ID: %w", err)
	}

	if len(result.Categories) == 0 {
		return nil, dbtypes.ErrNotFound
	}

	return result.Categories[0], nil
}

// CreateCategory creates a new category.
func (db *DGraphDatabase) CreateCategory(ctx context.Context, category *model.Category) error {
	_, err := db.Mutate(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to create category: %w", err)
	}
	return nil
}

// UpdateCategory updates an existing category.
func (db *DGraphDatabase) UpdateCategory(ctx context.Context, category *model.Category) error {
	_, err := db.Mutate(ctx, category)
	if err != nil {
		return fmt.Errorf("failed to update category: %w", err)
	}
	return nil
}

// DeleteCategory deletes a category.
func (db *DGraphDatabase) DeleteCategory(ctx context.Context, categoryID string) error {
	if categoryID == "" {
		return dbtypes.ErrInvalidID
	}

	txn := db.client.NewTxn()
	defer func() {
		if err := txn.Discard(ctx); err != nil {
			log.Printf("error discarding DeleteCategory transaction: %v", err)
		}
	}()

	// Delete category
	d := map[string]string{"uid": categoryID}
	deleteJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal category deletion data: %w", err)
	}

	mu := &api.Mutation{
		DeleteJson: deleteJSON,
		CommitNow:  true,
	}

	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// CreateReview creates a new review.
func (db *DGraphDatabase) CreateReview(ctx context.Context, review *model.Review) error {
	_, err := db.Mutate(ctx, review)
	if err != nil {
		return fmt.Errorf("failed to create review: %w", err)
	}
	return nil
}

// GetReview fetches a review by ID.
func (db *DGraphDatabase) GetReview(ctx context.Context, reviewID string) (*model.Review, error) {
	if reviewID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	reviewQuery := `
	query GetReview($id: string) {
		review(func: uid($id)) {
			uid
			rating
			comment
			author {
				uid
				name
			}
			recipe {
				uid
				title
			}
			createdAt
			updatedAt
		}
	}`

	vars := map[string]string{
		"$id": reviewID,
	}

	var result struct {
		Reviews []*model.Review `json:"review"`
	}

	err := db.Query(ctx, reviewQuery, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query review by ID: %w", err)
	}

	if len(result.Reviews) == 0 {
		return nil, dbtypes.ErrNotFound
	}

	return result.Reviews[0], nil
}

// UpdateReview updates an existing review.
func (db *DGraphDatabase) UpdateReview(ctx context.Context, review *model.Review) error {
	_, err := db.Mutate(ctx, review)
	if err != nil {
		return fmt.Errorf("failed to update review: %w", err)
	}
	return nil
}

// DeleteReview deletes a review.
func (db *DGraphDatabase) DeleteReview(ctx context.Context, reviewID string) error {
	if reviewID == "" {
		return dbtypes.ErrInvalidID
	}

	txn := db.client.NewTxn()
	defer func() {
		if err := txn.Discard(ctx); err != nil {
			log.Printf("error discarding DeleteReview transaction: %v", err)
		}
	}()

	// Delete review
	d := map[string]string{"uid": reviewID}
	deleteJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal review deletion data: %w", err)
	}

	mu := &api.Mutation{
		DeleteJson: deleteJSON,
		CommitNow:  true,
	}

	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to delete review: %w", err)
	}

	return nil
}

// GetRecipe fetches a recipe by ID.
//
//nolint:funlen // This function is necessarily long due to the complex query structure
func (db *DGraphDatabase) GetRecipe(ctx context.Context, recipeID string) (*model.Recipe, error) {
	if recipeID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	recipeQuery := `
	query GetRecipe($id: string) {
		recipe(func: uid($id)) {
			uid
			title
			description
			author {
				uid
				name
			}
			categories {
				uid
				name
			}
			cuisine
			prepTime
			cookTime
			servings
			difficulty
			ingredients {
				uid
				name
				quantity
				unit
				preparation
				substitutes
				isOptional
			}
			steps {
				uid
				orderIndex
				description
				timeEstimate
			}
			nutritionInfo {
				calories
				protein
				carbs
				fat
				fiber
				sugar
				sodium
			}
			images {
				uid
				url
				alt
				width
				height
			}
			tags
			likes
			reviews {
				uid
				author {
					uid
					name
				}
				rating
				comment
				createdAt
				updatedAt
			}
			averageRating
			createdAt
			updatedAt
		}
	}`

	vars := map[string]string{
		"$id": recipeID,
	}

	var result struct {
		Recipes []*model.Recipe `json:"recipe"`
	}

	err := db.Query(ctx, recipeQuery, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipe by ID: %w", err)
	}

	if len(result.Recipes) == 0 {
		return nil, dbtypes.ErrNotFound
	}

	return result.Recipes[0], nil
}

// GetRecipes fetches recipes based on filters.
//
//nolint:funlen // This function is necessarily long due to the complex filter handling
func (db *DGraphDatabase) GetRecipes(
	ctx context.Context, filter map[string]string, first, offset int,
) ([]*model.Recipe, error) {
	// Construct filter conditions
	conditions := ""
	vars := make(map[string]string)

	for key, value := range filter {
		if value == "" {
			continue
		}

		switch key {
		case "title":
			conditions += ", anyoftext(title, $title)"
			vars["title"] = value
		case "category":
			conditions += ", type(Recipe) AND anyoftext(categories, $category)"
			vars["category"] = value
		case "cuisine":
			conditions += ", eq(cuisine, $cuisine)"
			vars["cuisine"] = value
		case "authorID":
			conditions += ", uid_in(author, $authorID)"
			vars["authorID"] = value
		}
	}

	if conditions != "" {
		conditions = ", @filter(" + conditions[2:] + ")"
	}

	recipesQuery := fmt.Sprintf(`
	query GetRecipes(%s) {
		recipes(func: type(Recipe) %s, first: %d, offset: %d) {
			uid
			title
			description
			author {
				uid
				name
			}
			categories {
				uid
				name
			}
			cuisine
			prepTime
			cookTime
			difficulty
			images {
				uid
				url
				alt
			}
			likes
			averageRating
			createdAt
		}
	}`, BuildVarDeclarations(vars), conditions, first, offset)

	var result struct {
		Recipes []*model.Recipe `json:"recipes"`
	}

	err := db.Query(ctx, recipesQuery, vars, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query recipes: %w", err)
	}

	return result.Recipes, nil
}

// CreateRecipe creates a new recipe.
func (db *DGraphDatabase) CreateRecipe(ctx context.Context, recipe *model.Recipe) error {
	_, err := db.Mutate(ctx, recipe)
	if err != nil {
		return fmt.Errorf("failed to create recipe: %w", err)
	}

	return nil
}

// UpdateRecipe updates an existing recipe.
func (db *DGraphDatabase) UpdateRecipe(ctx context.Context, recipe *model.Recipe) error {
	_, err := db.Mutate(ctx, recipe)
	if err != nil {
		return fmt.Errorf("failed to update recipe: %w", err)
	}

	return nil
}

// DeleteRecipe deletes a recipe.
func (db *DGraphDatabase) DeleteRecipe(ctx context.Context, recipeID string) error {
	if recipeID == "" {
		return dbtypes.ErrInvalidID
	}

	txn := db.client.NewTxn()
	defer func() {
		if err := txn.Discard(ctx); err != nil {
			log.Printf("error discarding DeleteRecipe transaction: %v", err)
		}
	}()

	// Delete recipe
	d := map[string]string{"uid": recipeID}
	deleteJSON, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal recipe deletion data: %w", err)
	}

	mu := &api.Mutation{
		DeleteJson: deleteJSON,
		CommitNow:  true,
	}

	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		return fmt.Errorf("failed to delete recipe: %w", err)
	}

	return nil
}

// BuildVarDeclarations creates a string of variable declarations for DGraph queries.
// This function is exported for testing purposes.
func BuildVarDeclarations(vars map[string]string) string {
	if len(vars) == 0 {
		return ""
	}

	result := ""
	for k := range vars {
		if result != "" {
			result += ", "
		}
		result += "$" + k + ": string"
	}

	return result
}
