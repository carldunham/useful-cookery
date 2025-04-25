package database_test

import (
	"errors"
	"testing"
	"time"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/model"
)

func TestInMemoryDatabaseCreation(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Check that the database is not nil
	if db == nil {
		t.Fatal("In-memory database should not be nil")
	}

	// No need to check interface implementation since we're returning concrete type
}

//nolint:cyclop,funlen // Test functions are allowed to have higher complexity and length
func TestInMemoryDatabase_User(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Create a test user
	user := &model.User{
		ID:        "test-user-id",
		Name:      "Test User",
		Email:     "test@example.com",
		Password:  "password",
		Role:      model.UserRole,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test CreateUser
	err = db.CreateUser(t.Context(), user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Test GetUser
	retrievedUser, err := db.GetUser(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("Failed to get user: %v", err)
	}
	if retrievedUser.ID != user.ID {
		t.Errorf("Retrieved user ID %s does not match original ID %s", retrievedUser.ID, user.ID)
	}
	if retrievedUser.Email != user.Email {
		t.Errorf("Retrieved user email %s does not match original email %s", retrievedUser.Email, user.Email)
	}

	// Test GetUserByEmail
	retrievedUserByEmail, err := db.GetUserByEmail(t.Context(), user.Email)
	if err != nil {
		t.Fatalf("Failed to get user by email: %v", err)
	}
	if retrievedUserByEmail.ID != user.ID {
		t.Errorf("Retrieved user ID %s does not match original ID %s", retrievedUserByEmail.ID, user.ID)
	}

	// Test GetUsers
	users, err := db.GetUsers(t.Context(), 10, 0)
	if err != nil {
		t.Fatalf("Failed to get users: %v", err)
	}
	if len(users) != 1 {
		t.Errorf("Expected 1 user, got %d", len(users))
	}

	// Test UpdateUser
	user.Name = "Updated User"
	err = db.UpdateUser(t.Context(), user)
	if err != nil {
		t.Fatalf("Failed to update user: %v", err)
	}
	retrievedUser, err = db.GetUser(t.Context(), user.ID)
	if err != nil {
		t.Fatalf("Failed to get updated user: %v", err)
	}
	if retrievedUser.Name != "Updated User" {
		t.Errorf("Retrieved user name %s does not match updated name %s", retrievedUser.Name, "Updated User")
	}

	// Test error cases
	_, err = db.GetUser(t.Context(), "non-existent-id")
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	_, err = db.GetUserByEmail(t.Context(), "non-existent@example.com")
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	err = db.CreateUser(t.Context(), &model.User{Email: user.Email})
	if !errors.Is(err, dbtypes.ErrAlreadyExists) {
		t.Errorf("Expected ErrAlreadyExists, got %v", err)
	}

	err = db.CreateUser(t.Context(), &model.User{})
	if !errors.Is(err, dbtypes.ErrEmailRequired) {
		t.Errorf("Expected ErrEmailRequired, got %v", err)
	}

	err = db.UpdateUser(t.Context(), &model.User{ID: "non-existent-id"})
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

//nolint:cyclop,funlen // Test functions are allowed to have higher complexity and length
func TestInMemoryDatabase_Category(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Create a test category
	category := &model.Category{
		ID:          "test-category-id",
		Name:        "Test Category",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test CreateCategory
	err = db.CreateCategory(t.Context(), category)
	if err != nil {
		t.Fatalf("Failed to create category: %v", err)
	}

	// Test GetCategory
	retrievedCategory, err := db.GetCategory(t.Context(), category.ID)
	if err != nil {
		t.Fatalf("Failed to get category: %v", err)
	}
	if retrievedCategory.ID != category.ID {
		t.Errorf("Retrieved category ID %s does not match original ID %s", retrievedCategory.ID, category.ID)
	}
	if retrievedCategory.Name != category.Name {
		t.Errorf("Retrieved category name %s does not match original name %s", retrievedCategory.Name, category.Name)
	}

	// Test UpdateCategory
	category.Name = "Updated Category"
	err = db.UpdateCategory(t.Context(), category)
	if err != nil {
		t.Fatalf("Failed to update category: %v", err)
	}
	retrievedCategory, err = db.GetCategory(t.Context(), category.ID)
	if err != nil {
		t.Fatalf("Failed to get updated category: %v", err)
	}
	if retrievedCategory.Name != "Updated Category" {
		t.Errorf("Retrieved category name %s does not match updated name %s", retrievedCategory.Name, "Updated Category")
	}

	// Test DeleteCategory
	err = db.DeleteCategory(t.Context(), category.ID)
	if err != nil {
		t.Fatalf("Failed to delete category: %v", err)
	}
	_, err = db.GetCategory(t.Context(), category.ID)
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}

	// Test error cases
	_, err = db.GetCategory(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.UpdateCategory(t.Context(), &model.Category{ID: "non-existent-id"})
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	err = db.DeleteCategory(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.DeleteCategory(t.Context(), "non-existent-id")
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

//nolint:cyclop,funlen // Test functions are allowed to have higher complexity and length
func TestInMemoryDatabase_Recipe(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Create a test recipe
	recipe := &model.Recipe{
		ID:          "test-recipe-id",
		Title:       "Test Recipe",
		Description: "Test Description",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Test CreateRecipe
	err = db.CreateRecipe(t.Context(), recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Test GetRecipe
	retrievedRecipe, err := db.GetRecipe(t.Context(), recipe.ID)
	if err != nil {
		t.Fatalf("Failed to get recipe: %v", err)
	}
	if retrievedRecipe.ID != recipe.ID {
		t.Errorf("Retrieved recipe ID %s does not match original ID %s", retrievedRecipe.ID, recipe.ID)
	}
	if retrievedRecipe.Title != recipe.Title {
		t.Errorf("Retrieved recipe title %s does not match original title %s", retrievedRecipe.Title, recipe.Title)
	}

	// Test GetRecipes
	recipes, err := db.GetRecipes(t.Context(), nil, 10, 0)
	if err != nil {
		t.Fatalf("Failed to get recipes: %v", err)
	}
	if len(recipes) != 1 {
		t.Errorf("Expected 1 recipe, got %d", len(recipes))
	}

	// Test UpdateRecipe
	recipe.Title = "Updated Recipe"
	err = db.UpdateRecipe(t.Context(), recipe)
	if err != nil {
		t.Fatalf("Failed to update recipe: %v", err)
	}
	retrievedRecipe, err = db.GetRecipe(t.Context(), recipe.ID)
	if err != nil {
		t.Fatalf("Failed to get updated recipe: %v", err)
	}
	if retrievedRecipe.Title != "Updated Recipe" {
		t.Errorf("Retrieved recipe title %s does not match updated title %s", retrievedRecipe.Title, "Updated Recipe")
	}

	// Test DeleteRecipe
	err = db.DeleteRecipe(t.Context(), recipe.ID)
	if err != nil {
		t.Fatalf("Failed to delete recipe: %v", err)
	}
	_, err = db.GetRecipe(t.Context(), recipe.ID)
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}

	// Test error cases
	_, err = db.GetRecipe(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.UpdateRecipe(t.Context(), &model.Recipe{ID: "non-existent-id"})
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	err = db.DeleteRecipe(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.DeleteRecipe(t.Context(), "non-existent-id")
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

//nolint:cyclop,funlen // Test functions are allowed to have higher complexity and length
func TestInMemoryDatabase_Review(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Create a test recipe for the review
	recipe := &model.Recipe{
		ID:    "test-recipe-id",
		Title: "Test Recipe",
	}
	err = db.CreateRecipe(t.Context(), recipe)
	if err != nil {
		t.Fatalf("Failed to create recipe: %v", err)
	}

	// Create a test review
	review := &model.Review{
		ID:        "test-review-id",
		Recipe:    recipe,
		Rating:    5,
		Comment:   "Great recipe!",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test CreateReview
	err = db.CreateReview(t.Context(), review)
	if err != nil {
		t.Fatalf("Failed to create review: %v", err)
	}

	// Test GetReview
	retrievedReview, err := db.GetReview(t.Context(), review.ID)
	if err != nil {
		t.Fatalf("Failed to get review: %v", err)
	}
	if retrievedReview.ID != review.ID {
		t.Errorf("Retrieved review ID %s does not match original ID %s", retrievedReview.ID, review.ID)
	}
	if retrievedReview.Rating != review.Rating {
		t.Errorf("Retrieved review rating %d does not match original rating %d", retrievedReview.Rating, review.Rating)
	}

	// Test UpdateReview
	review.Rating = 4
	err = db.UpdateReview(t.Context(), review)
	if err != nil {
		t.Fatalf("Failed to update review: %v", err)
	}
	retrievedReview, err = db.GetReview(t.Context(), review.ID)
	if err != nil {
		t.Fatalf("Failed to get updated review: %v", err)
	}
	if retrievedReview.Rating != 4 {
		t.Errorf("Retrieved review rating %d does not match updated rating %d", retrievedReview.Rating, 4)
	}

	// Test DeleteReview
	err = db.DeleteReview(t.Context(), review.ID)
	if err != nil {
		t.Fatalf("Failed to delete review: %v", err)
	}
	_, err = db.GetReview(t.Context(), review.ID)
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound after deletion, got %v", err)
	}

	// Test error cases
	_, err = db.GetReview(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.UpdateReview(t.Context(), &model.Review{ID: "non-existent-id"})
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}

	err = db.DeleteReview(t.Context(), "")
	if !errors.Is(err, dbtypes.ErrInvalidID) {
		t.Errorf("Expected ErrInvalidID, got %v", err)
	}

	err = db.DeleteReview(t.Context(), "non-existent-id")
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Errorf("Expected ErrNotFound, got %v", err)
	}
}

func TestInMemoryDatabase_Query(t *testing.T) {
	t.Parallel()
	// Create a new in-memory database
	db, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}

	// Test Query (should return an error since it's not supported)
	err = db.Query(t.Context(), "", nil, nil)
	if err == nil {
		t.Error("Expected an error for Query, got nil")
	}
}
