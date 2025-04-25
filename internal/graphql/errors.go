package graphql

import (
	"errors"
)

// Common errors.
var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("entity not found")

	// ErrRecipeNotFound is returned when a recipe is not found.
	ErrRecipeNotFound = errors.New("recipe not found")

	// ErrReviewNotFound is returned when a review is not found.
	ErrReviewNotFound = errors.New("review not found")

	// ErrCategoryNotFound is returned when a category is not found.
	ErrCategoryNotFound = errors.New("category not found")
)
