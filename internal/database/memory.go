package database

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/model"
)

// InMemoryDatabase implements database operations using in-memory storage.
// This is primarily used for testing.
type InMemoryDatabase struct {
	users        map[string]*model.User
	usersByEmail map[string]*model.User
	categories   map[string]*model.Category
	recipes      map[string]*model.Recipe
	reviews      map[string]*model.Review

	mu sync.RWMutex
}

// NewInMemoryDatabase creates a new in-memory database instance.
func NewInMemoryDatabase(_ dbtypes.DatabaseOptions) (*InMemoryDatabase, error) {
	return &InMemoryDatabase{
		users:        make(map[string]*model.User),
		usersByEmail: make(map[string]*model.User),
		categories:   make(map[string]*model.Category),
		recipes:      make(map[string]*model.Recipe),
		reviews:      make(map[string]*model.Review),
	}, nil
}

// ErrRawQueriesNotSupported is returned when attempting to execute a raw query on the in-memory database.
var ErrRawQueriesNotSupported = errors.New("raw queries not supported in in-memory database")

// Query executes a query against the in-memory database.
// This is a simplified implementation that doesn't support all DGraph query features.
func (db *InMemoryDatabase) Query(_ context.Context, _ string, _ map[string]string, _ any) error {
	return ErrRawQueriesNotSupported
}

// GetUser fetches a user by ID.
func (db *InMemoryDatabase) GetUser(_ context.Context, userID string) (*model.User, error) {
	if userID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.users[userID]
	if !ok {
		return nil, dbtypes.ErrNotFound
	}

	return user, nil
}

// GetUserByEmail fetches a user by email.
func (db *InMemoryDatabase) GetUserByEmail(_ context.Context, email string) (*model.User, error) {
	if email == "" {
		return nil, dbtypes.ErrEmailRequired
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	user, ok := db.usersByEmail[email]
	if !ok {
		return nil, dbtypes.ErrNotFound
	}

	return user, nil
}

// CreateUser creates a new user.
func (db *InMemoryDatabase) CreateUser(_ context.Context, user *model.User) error {
	if user.Email == "" {
		return dbtypes.ErrEmailRequired
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if user already exists
	if _, ok := db.usersByEmail[user.Email]; ok {
		return dbtypes.ErrAlreadyExists
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	user.UpdatedAt = now

	// Store user
	db.users[user.ID] = user
	db.usersByEmail[user.Email] = user

	return nil
}

// GetUsers fetches a list of users with pagination.
func (db *InMemoryDatabase) GetUsers(_ context.Context, limit, offset int) ([]*model.User, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Convert map to slice
	users := make([]*model.User, 0, len(db.users))
	for _, user := range db.users {
		users = append(users, user)
	}

	// Apply pagination
	if limit <= 0 {
		limit = len(users)
	}
	if offset < 0 {
		offset = 0
	}

	// Check bounds
	if offset >= len(users) {
		return []*model.User{}, nil
	}

	end := offset + limit
	if end > len(users) {
		end = len(users)
	}

	return users[offset:end], nil
}

// UpdateUser updates a user.
func (db *InMemoryDatabase) UpdateUser(_ context.Context, user *model.User) error {
	if user.ID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if user exists
	existingUser, ok := db.users[user.ID]
	if !ok {
		return dbtypes.ErrNotFound
	}

	// Update email index if email changed
	if existingUser.Email != user.Email {
		delete(db.usersByEmail, existingUser.Email)
		db.usersByEmail[user.Email] = user
	}

	// Update timestamp
	user.UpdatedAt = time.Now()

	// Store updated user
	db.users[user.ID] = user

	return nil
}

// GetCategory fetches a category by ID.
func (db *InMemoryDatabase) GetCategory(_ context.Context, categoryID string) (*model.Category, error) {
	if categoryID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	category, ok := db.categories[categoryID]
	if !ok {
		return nil, dbtypes.ErrNotFound
	}

	return category, nil
}

// CreateCategory creates a new category.
func (db *InMemoryDatabase) CreateCategory(_ context.Context, category *model.Category) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Generate ID if not provided
	if category.ID == "" {
		category.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if category.CreatedAt.IsZero() {
		category.CreatedAt = now
	}
	category.UpdatedAt = now

	// Store category
	db.categories[category.ID] = category

	return nil
}

// UpdateCategory updates an existing category.
func (db *InMemoryDatabase) UpdateCategory(_ context.Context, category *model.Category) error {
	if category.ID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if category exists
	if _, ok := db.categories[category.ID]; !ok {
		return dbtypes.ErrNotFound
	}

	// Update timestamp
	category.UpdatedAt = time.Now()

	// Store updated category
	db.categories[category.ID] = category

	return nil
}

// DeleteCategory deletes a category.
func (db *InMemoryDatabase) DeleteCategory(_ context.Context, categoryID string) error {
	if categoryID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if category exists
	if _, ok := db.categories[categoryID]; !ok {
		return dbtypes.ErrNotFound
	}

	// Delete category
	delete(db.categories, categoryID)

	return nil
}

// CreateReview creates a new review.
func (db *InMemoryDatabase) CreateReview(_ context.Context, review *model.Review) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Generate ID if not provided
	if review.ID == "" {
		review.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if review.CreatedAt.IsZero() {
		review.CreatedAt = now
	}
	review.UpdatedAt = now

	// Store review
	db.reviews[review.ID] = review

	// Update recipe's average rating
	updateRecipeRating(db, review)

	return nil
}

// GetReview fetches a review by ID.
func (db *InMemoryDatabase) GetReview(_ context.Context, reviewID string) (*model.Review, error) {
	if reviewID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	review, ok := db.reviews[reviewID]
	if !ok {
		return nil, dbtypes.ErrNotFound
	}

	return review, nil
}

// UpdateReview updates an existing review.
func (db *InMemoryDatabase) UpdateReview(_ context.Context, review *model.Review) error {
	if review.ID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if review exists
	existingReview, ok := db.reviews[review.ID]
	if !ok {
		return dbtypes.ErrNotFound
	}

	// Update timestamp
	review.UpdatedAt = time.Now()

	// Store updated review
	db.reviews[review.ID] = review

	// Update recipe's average rating if rating changed
	updateReviewRating(db, review, existingReview)

	return nil
}

// DeleteReview deletes a review.
func (db *InMemoryDatabase) DeleteReview(_ context.Context, reviewID string) error {
	if reviewID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if review exists
	review, ok := db.reviews[reviewID]
	if !ok {
		return dbtypes.ErrNotFound
	}

	// Delete review
	delete(db.reviews, reviewID)

	// Update recipe's average rating
	removeReviewFromRecipe(db, review, reviewID)

	return nil
}

// GetRecipe fetches a recipe by ID.
func (db *InMemoryDatabase) GetRecipe(_ context.Context, recipeID string) (*model.Recipe, error) {
	if recipeID == "" {
		return nil, dbtypes.ErrInvalidID
	}

	db.mu.RLock()
	defer db.mu.RUnlock()

	recipe, ok := db.recipes[recipeID]
	if !ok {
		return nil, dbtypes.ErrNotFound
	}

	return recipe, nil
}

// GetRecipes fetches recipes based on filters.
func (db *InMemoryDatabase) GetRecipes(
	_ context.Context, filter map[string]string, first, offset int,
) ([]*model.Recipe, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	// Convert map to slice
	allRecipes := make([]*model.Recipe, 0, len(db.recipes))
	for _, recipe := range db.recipes {
		allRecipes = append(allRecipes, recipe)
	}

	// Apply filters
	filteredRecipes := make([]*model.Recipe, 0)
	for _, recipe := range allRecipes {
		if matchesFilter(recipe, filter) {
			filteredRecipes = append(filteredRecipes, recipe)
		}
	}

	// Apply pagination
	if first <= 0 {
		first = len(filteredRecipes)
	}
	if offset < 0 {
		offset = 0
	}

	// Check bounds
	if offset >= len(filteredRecipes) {
		return []*model.Recipe{}, nil
	}

	end := offset + first
	if end > len(filteredRecipes) {
		end = len(filteredRecipes)
	}

	return filteredRecipes[offset:end], nil
}

// updateReviewRating updates the recipe's average rating when a review is updated.
func updateReviewRating(db *InMemoryDatabase, review, existingReview *model.Review) {
	// Skip if review doesn't have a recipe or recipe ID, or if rating hasn't changed
	if review.Recipe == nil || review.Recipe.ID == "" || existingReview.Rating == review.Rating {
		return
	}

	// Get the recipe
	recipe, ok := db.recipes[review.Recipe.ID]
	if !ok || len(recipe.Reviews) == 0 {
		return
	}

	// Calculate new average rating
	totalRating := recipe.AverageRating * float64(len(recipe.Reviews))
	totalRating -= float64(existingReview.Rating)
	totalRating += float64(review.Rating)
	recipe.AverageRating = totalRating / float64(len(recipe.Reviews))

	// Update review in recipe
	for i, r := range recipe.Reviews {
		if r.ID == review.ID {
			recipe.Reviews[i] = *review
			break
		}
	}

	// Update recipe
	db.recipes[recipe.ID] = recipe
}

// removeReviewFromRecipe removes a review from a recipe and updates the average rating.
func removeReviewFromRecipe(db *InMemoryDatabase, review *model.Review, reviewID string) {
	// Skip if review doesn't have a recipe or recipe ID
	if review.Recipe == nil || review.Recipe.ID == "" {
		return
	}

	// Get the recipe
	recipe, ok := db.recipes[review.Recipe.ID]
	if !ok || len(recipe.Reviews) == 0 {
		return
	}

	// Calculate new average rating
	totalRating := recipe.AverageRating * float64(len(recipe.Reviews))
	totalRating -= float64(review.Rating)

	// Remove review from recipe
	newReviews := make([]model.Review, 0, len(recipe.Reviews)-1)
	for _, r := range recipe.Reviews {
		if r.ID != reviewID {
			newReviews = append(newReviews, r)
		}
	}
	recipe.Reviews = newReviews

	// Update average rating
	if len(recipe.Reviews) > 0 {
		recipe.AverageRating = totalRating / float64(len(recipe.Reviews))
	} else {
		recipe.AverageRating = 0
	}

	// Update recipe
	db.recipes[recipe.ID] = recipe
}

// updateRecipeRating updates the recipe's average rating when a review is created.
func updateRecipeRating(db *InMemoryDatabase, review *model.Review) {
	// Skip if review doesn't have a recipe or recipe ID
	if review.Recipe == nil || review.Recipe.ID == "" {
		return
	}

	// Get the recipe
	recipe, ok := db.recipes[review.Recipe.ID]
	if !ok {
		return
	}

	// Calculate new average rating
	totalRating := recipe.AverageRating * float64(len(recipe.Reviews))
	totalRating += float64(review.Rating)
	recipe.AverageRating = totalRating / float64(len(recipe.Reviews)+1)

	// Add review to recipe
	recipe.Reviews = append(recipe.Reviews, *review)

	// Update recipe
	db.recipes[recipe.ID] = recipe
}

// matchesFilter checks if a recipe matches the given filter.
//
//nolint:cyclop // This function necessarily has high cyclomatic complexity due to filter matching logic
func matchesFilter(recipe *model.Recipe, filter map[string]string) bool {
	for key, value := range filter {
		if value == "" {
			continue
		}

		switch key {
		case "title":
			if recipe.Title != value {
				return false
			}
		case "category":
			found := false
			for _, category := range recipe.Categories {
				if category.Name == value {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		case "cuisine":
			if recipe.Cuisine != value {
				return false
			}
		case "authorID":
			if recipe.Author == nil || recipe.Author.ID != value {
				return false
			}
		}
	}

	return true
}

// CreateRecipe creates a new recipe.
func (db *InMemoryDatabase) CreateRecipe(_ context.Context, recipe *model.Recipe) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	// Generate ID if not provided
	if recipe.ID == "" {
		recipe.ID = model.NewID()
	}

	// Set timestamps
	now := time.Now()
	if recipe.CreatedAt.IsZero() {
		recipe.CreatedAt = now
	}
	recipe.UpdatedAt = now

	// Store recipe
	db.recipes[recipe.ID] = recipe

	return nil
}

// UpdateRecipe updates an existing recipe.
func (db *InMemoryDatabase) UpdateRecipe(_ context.Context, recipe *model.Recipe) error {
	if recipe.ID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if recipe exists
	if _, ok := db.recipes[recipe.ID]; !ok {
		return dbtypes.ErrNotFound
	}

	// Update timestamp
	recipe.UpdatedAt = time.Now()

	// Store updated recipe
	db.recipes[recipe.ID] = recipe

	return nil
}

// DeleteRecipe deletes a recipe.
func (db *InMemoryDatabase) DeleteRecipe(_ context.Context, recipeID string) error {
	if recipeID == "" {
		return dbtypes.ErrInvalidID
	}

	db.mu.Lock()
	defer db.mu.Unlock()

	// Check if recipe exists
	if _, ok := db.recipes[recipeID]; !ok {
		return dbtypes.ErrNotFound
	}

	// Delete recipe
	delete(db.recipes, recipeID)

	return nil
}
