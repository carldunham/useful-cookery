package database

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/dgraph-io/dgo/v2"
	"github.com/dgraph-io/dgo/v2/protos/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Common errors
var (
	ErrNotFound      = errors.New("entity not found")
	ErrInvalidID     = errors.New("invalid ID")
	ErrAlreadyExists = errors.New("entity already exists")
)

// DGraphClient is a client for DGraph operations
type DGraphClient struct {
	client *dgo.Dgraph
}

// NewDGraphClient creates a new DGraph client
func NewDGraphClient(hosts []string) (*DGraphClient, error) {
	// Connect to all hosts
	conns := make([]*grpc.ClientConn, len(hosts))
	clients := make([]api.DgraphClient, len(hosts))

	for i, host := range hosts {
		conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			return nil, fmt.Errorf("failed to connect to DGraph at %s: %w", host, err)
		}
		conns[i] = conn
		clients[i] = api.NewDgraphClient(conn)
	}

	// Create DGraph client
	dgraphClient := dgo.NewDgraphClient(clients...)

	return &DGraphClient{
		client: dgraphClient,
	}, nil
}

// Query executes a GraphQL+ query against DGraph
func (c *DGraphClient) Query(ctx context.Context, q string, vars map[string]string, result interface{}) error {
	txn := c.client.NewTxn()
	defer txn.Discard(ctx)

	// Create request
	req := &api.Request{
		Query: q,
		Vars:  vars,
	}

	// Execute query
	resp, err := txn.Do(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to execute query: %w", err)
	}

	// Unmarshal result
	err = json.Unmarshal(resp.Json, result)
	if err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return nil
}

// Mutate executes a mutation against DGraph
func (c *DGraphClient) Mutate(ctx context.Context, data interface{}) (*api.Response, error) {
	txn := c.client.NewTxn()
	defer txn.Discard(ctx)

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

// GetUser fetches a user by ID
func (c *DGraphClient) GetUser(ctx context.Context, userID string) (*models.User, error) {
	if userID == "" {
		return nil, ErrInvalidID
	}

	q := `
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
		"id": userID,
	}

	var result struct {
		Users []*models.User `json:"user"`
	}

	err := c.Query(ctx, q, vars, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Users) == 0 {
		return nil, ErrNotFound
	}

	return result.Users[0], nil
}

// GetUserByEmail fetches a user by email
func (c *DGraphClient) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	if email == "" {
		return nil, errors.New("email is required")
	}

	q := `
	query GetUserByEmail($email: string) {
		user(func: eq(email, $email)) {
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

	vars := map[string]string{
		"email": email,
	}

	var result struct {
		Users []*models.User `json:"user"`
	}

	err := c.Query(ctx, q, vars, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Users) == 0 {
		return nil, ErrNotFound
	}

	return result.Users[0], nil
}

// CreateUser creates a new user
func (c *DGraphClient) CreateUser(ctx context.Context, user *models.User) error {
	// Check if user already exists
	existingUser, err := c.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return ErrAlreadyExists
	}

	// Create user
	_, err = c.Mutate(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// GetUsers fetches a list of users with pagination
func (c *DGraphClient) GetUsers(ctx context.Context, limit, offset int) ([]*models.User, error) {
	q := `
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
		"limit":  fmt.Sprintf("%d", limit),
		"offset": fmt.Sprintf("%d", offset),
	}

	var result struct {
		Users []*models.User `json:"users"`
	}

	err := c.Query(ctx, q, vars, &result)
	if err != nil {
		return nil, err
	}

	return result.Users, nil
}

// UpdateUser updates a user
func (c *DGraphClient) UpdateUser(ctx context.Context, user *models.User) error {
	_, err := c.Mutate(ctx, user)
	return err
}

// GetCategory fetches a category by ID
func (c *DGraphClient) GetCategory(ctx context.Context, categoryID string) (*models.Category, error) {
	if categoryID == "" {
		return nil, ErrInvalidID
	}

	q := `
	query GetCategory($id: string) {
		category(func: uid($id)) {
			uid
			name
			description
		}
	}`

	vars := map[string]string{
		"id": categoryID,
	}

	var result struct {
		Categories []*models.Category `json:"category"`
	}

	err := c.Query(ctx, q, vars, &result)
	if err != nil {
		return nil, err
	}

	if len(result.Categories) == 0 {
		return nil, ErrNotFound
	}

	return result.Categories[0], nil
}

// CreateReview creates a new review
func (c *DGraphClient) CreateReview(ctx context.Context, review *models.Review) error {
	_, err := c.Mutate(ctx, review)
	return err
}

// GetRecipes fetches recipes based on filters
func (c *DGraphClient) GetRecipes(ctx context.Context, filter map[string]string, first, offset int) ([]*models.Recipe, error) {
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
			conditions += ", has(categories) AND anyoftext(categories, $category)"
			vars["category"] = value
		case "cuisine":
			conditions += ", eq(cuisine, $cuisine)"
			vars["cuisine"] = value
		case "authorId":
			conditions += ", uid_in(author, $authorId)"
			vars["authorId"] = value
		}
	}

	if conditions != "" {
		conditions = "@filter(" + conditions[2:] + ")"
	}

	q := fmt.Sprintf(`
	query GetRecipes(%s) {
		recipes(func: has(title), %s, first: %d, offset: %d) {
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
	}`, buildVarDeclarations(vars), conditions, first, offset)

	var result struct {
		Recipes []*models.Recipe `json:"recipes"`
	}

	err := c.Query(ctx, q, vars, &result)
	if err != nil {
		return nil, err
	}

	return result.Recipes, nil
}

// CreateRecipe creates a new recipe
func (c *DGraphClient) CreateRecipe(ctx context.Context, recipe *models.Recipe) error {
	_, err := c.Mutate(ctx, recipe)
	if err != nil {
		return err
	}

	return nil
}

// UpdateRecipe updates an existing recipe
func (c *DGraphClient) UpdateRecipe(ctx context.Context, recipe *models.Recipe) error {
	_, err := c.Mutate(ctx, recipe)
	if err != nil {
		return err
	}

	return nil
}

// DeleteRecipe deletes a recipe
func (c *DGraphClient) DeleteRecipe(ctx context.Context, recipeID string) error {
	if recipeID == "" {
		return ErrInvalidID
	}

	txn := c.client.NewTxn()
	defer txn.Discard(ctx)

	// Delete recipe
	d := map[string]string{"uid": recipeID}
	deleteJson, err := json.Marshal(d)
	if err != nil {
		return err
	}

	mu := &api.Mutation{
		DeleteJson: deleteJson,
		CommitNow:  true,
	}

	_, err = txn.Mutate(ctx, mu)
	if err != nil {
		return err
	}

	return nil
}

// buildVarDeclarations creates a string of variable declarations for DGraph queries
func buildVarDeclarations(vars map[string]string) string {
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
