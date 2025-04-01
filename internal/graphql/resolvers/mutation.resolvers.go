package resolvers

import (
	"context"
	"errors"
	"time"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/graphql/models"
	domainModels "github.com/carldunham/useful-cookery/internal/model"
)

// Register registers a new user
func (r *Resolver) Register(ctx context.Context, input models.RegisterInput) (*models.AuthPayload, error) {
	user, err := r.UserResolver.RegisterUser(ctx, input)
	if err != nil {
		return nil, err
	}

	token, err := r.AuthService.generateToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &models.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

// Login authenticates a user
func (r *Resolver) Login(ctx context.Context, email string, password string) (*models.AuthPayload, error) {
	token, user, err := r.AuthService.Login(ctx, email, password)
	if err != nil {
		return nil, err
	}

	return &models.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

// RegisterUser registers a new user
func (r *UserResolver) RegisterUser(ctx context.Context, input models.RegisterInput) (*domainModels.User, error) {
	user, err := r.AuthService.RegisterUser(ctx, input.Name, input.Email, input.Password)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// UpdateUser updates the current user's information
func (r *Resolver) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*domainModels.User, error) {
	return r.UserResolver.UpdateUser(ctx, input)
}

// UpdateUser updates a user's information
func (r *UserResolver) UpdateUser(ctx context.Context, input models.UpdateUserInput) (*domainModels.User, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Update user fields
	if input.Name != nil {
		currentUser.Name = *input.Name
	}
	if input.Email != nil {
		currentUser.Email = *input.Email
	}
	if input.Password != nil {
		// Hash the new password
		hashedPassword, err := r.AuthService.HashPassword(*input.Password)
		if err != nil {
			return nil, err
		}
		currentUser.Password = hashedPassword
	}

	// Update preferences if provided
	if input.Preferences != nil {
		if currentUser.Preferences == nil {
			currentUser.Preferences = &domainModels.UserPreferences{}
		}

		prefs := currentUser.Preferences

		if len(input.Preferences.DietaryRestrictions) > 0 {
			prefs.DietaryRestrictions = input.Preferences.DietaryRestrictions
		}

		if len(input.Preferences.FavoriteIngredients) > 0 {
			prefs.FavoriteIngredients = input.Preferences.FavoriteIngredients
		}

		if len(input.Preferences.DislikedIngredients) > 0 {
			prefs.DislikedIngredients = input.Preferences.DislikedIngredients
		}

		if input.Preferences.SkillLevel != nil {
			prefs.SkillLevel = *input.Preferences.SkillLevel
		}

		if len(input.Preferences.CuisinePreferences) > 0 {
			prefs.CuisinePreferences = input.Preferences.CuisinePreferences
		}
	}

	// Update timestamp
	currentUser.UpdatedAt = time.Now()

	// Save to database
	err := r.DB.UpdateUser(ctx, currentUser)
	if err != nil {
		return nil, err
	}

	return currentUser, nil
}

// CreateRecipe creates a new recipe
func (r *Resolver) CreateRecipe(ctx context.Context, input models.RecipeInput) (*domainModels.Recipe, error) {
	return r.RecipeResolver.CreateRecipe(ctx, input)
}

// CreateRecipe creates a new recipe
func (r *RecipeResolver) CreateRecipe(ctx context.Context, input models.RecipeInput) (*domainModels.Recipe, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Create recipe
	recipe := &domainModels.Recipe{
		ID:          domainModels.NewID(),
		Title:       input.Title,
		Author:      currentUser,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Ingredients: []domainModels.DetailedIngredient{},
		Steps:       []domainModels.Step{},
	}

	// Set optional fields
	if input.Description != nil {
		recipe.Description = *input.Description
	}
	if input.Cuisine != nil {
		recipe.Cuisine = *input.Cuisine
	}
	if input.PrepTime != nil {
		recipe.PrepTime = *input.PrepTime
	}
	if input.CookTime != nil {
		recipe.CookTime = *input.CookTime
	}
	if input.Servings != nil {
		recipe.Servings = *input.Servings
	}
	if input.Difficulty != nil {
		recipe.Difficulty = *input.Difficulty
	}
	if len(input.Tags) > 0 {
		recipe.Tags = input.Tags
	}

	// Process ingredients
	for _, ingInput := range input.Ingredients {
		ingredient := domainModels.DetailedIngredient{
			ID:   domainModels.NewID(),
			Name: ingInput.Name,
		}

		if ingInput.Quantity != nil {
			ingredient.Quantity = *ingInput.Quantity
		}
		if ingInput.Unit != nil {
			ingredient.Unit = *ingInput.Unit
		}
		if ingInput.Preparation != nil {
			ingredient.Preparation = *ingInput.Preparation
		}
		if len(ingInput.Substitutes) > 0 {
			ingredient.Substitutes = ingInput.Substitutes
		}
		if ingInput.IsOptional != nil {
			ingredient.IsOptional = *ingInput.IsOptional
		}

		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Process steps
	for _, stepInput := range input.Steps {
		step := domainModels.Step{
			ID:          domainModels.NewID(),
			OrderIndex:  stepInput.OrderIndex,
			Description: stepInput.Description,
		}

		if stepInput.TimeEstimate != nil {
			step.TimeEstimate = *stepInput.TimeEstimate
		}

		recipe.Steps = append(recipe.Steps, step)
	}

	// Process categories
	if len(input.CategoryIDs) > 0 {
		categories := make([]domainModels.Category, 0, len(input.CategoryIDs))
		for _, catID := range input.CategoryIDs {
			cat, err := r.DB.GetCategory(ctx, catID)
			if err != nil {
				continue // Skip if category not found
			}
			categories = append(categories, *cat)
		}
		recipe.Categories = categories
	}

	// Save to database
	err := r.DB.CreateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

// UpdateRecipe updates an existing recipe
func (r *Resolver) UpdateRecipe(ctx context.Context, id string, input models.RecipeInput) (*domainModels.Recipe, error) {
	return r.RecipeResolver.UpdateRecipe(ctx, id, input)
}

// UpdateRecipe updates an existing recipe
func (r *RecipeResolver) UpdateRecipe(ctx context.Context, id string, input models.RecipeInput) (*domainModels.Recipe, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Get existing recipe
	recipe, err := r.DB.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("recipe not found")
		}
		return nil, err
	}

	// Check if user has permission to update (author or admin)
	if recipe.Author != nil && recipe.Author.ID != currentUser.ID && currentUser.Role != "ADMIN" {
		return nil, errors.New("permission denied")
	}

	// Update basic fields
	recipe.Title = input.Title
	recipe.UpdatedAt = time.Now()

	// Update optional fields
	if input.Description != nil {
		recipe.Description = *input.Description
	}
	if input.Cuisine != nil {
		recipe.Cuisine = *input.Cuisine
	}
	if input.PrepTime != nil {
		recipe.PrepTime = *input.PrepTime
	}
	if input.CookTime != nil {
		recipe.CookTime = *input.CookTime
	}
	if input.Servings != nil {
		recipe.Servings = *input.Servings
	}
	if input.Difficulty != nil {
		recipe.Difficulty = *input.Difficulty
	}
	if len(input.Tags) > 0 {
		recipe.Tags = input.Tags
	}

	// Process ingredients (replace all)
	recipe.Ingredients = []domainModels.DetailedIngredient{}
	for _, ingInput := range input.Ingredients {
		ingredient := domainModels.DetailedIngredient{
			ID:   domainModels.NewID(),
			Name: ingInput.Name,
		}

		if ingInput.Quantity != nil {
			ingredient.Quantity = *ingInput.Quantity
		}
		if ingInput.Unit != nil {
			ingredient.Unit = *ingInput.Unit
		}
		if ingInput.Preparation != nil {
			ingredient.Preparation = *ingInput.Preparation
		}
		if len(ingInput.Substitutes) > 0 {
			ingredient.Substitutes = ingInput.Substitutes
		}
		if ingInput.IsOptional != nil {
			ingredient.IsOptional = *ingInput.IsOptional
		}

		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	// Process steps (replace all)
	recipe.Steps = []domainModels.Step{}
	for _, stepInput := range input.Steps {
		step := domainModels.Step{
			ID:          domainModels.NewID(),
			OrderIndex:  stepInput.OrderIndex,
			Description: stepInput.Description,
		}

		if stepInput.TimeEstimate != nil {
			step.TimeEstimate = *stepInput.TimeEstimate
		}

		recipe.Steps = append(recipe.Steps, step)
	}

	// Process categories
	if len(input.CategoryIDs) > 0 {
		categories := make([]domainModels.Category, 0, len(input.CategoryIDs))
		for _, catID := range input.CategoryIDs {
			cat, err := r.DB.GetCategory(ctx, catID)
			if err != nil {
				continue // Skip if category not found
			}
			categories = append(categories, *cat)
		}
		recipe.Categories = categories
	}

	// Save to database
	err = r.DB.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

// DeleteRecipe deletes a recipe
func (r *Resolver) DeleteRecipe(ctx context.Context, id string) (bool, error) {
	return r.RecipeResolver.DeleteRecipe(ctx, id)
}

// DeleteRecipe deletes a recipe
func (r *RecipeResolver) DeleteRecipe(ctx context.Context, id string) (bool, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return false, errors.New("not authenticated")
	}

	// Get recipe to check ownership
	recipe, err := r.DB.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return false, errors.New("recipe not found")
		}
		return false, err
	}

	// Check if user has permission to delete (author or admin)
	if recipe.Author != nil && recipe.Author.ID != currentUser.ID && currentUser.Role != "ADMIN" {
		return false, errors.New("permission denied")
	}

	// Delete from database
	err = r.DB.DeleteRecipe(ctx, id)
	if err != nil {
		return false, err
	}

	return true, nil
}

// LikeRecipe adds a like to a recipe
func (r *Resolver) LikeRecipe(ctx context.Context, id string) (*domainModels.Recipe, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Get recipe
	recipe, err := r.DB.GetRecipe(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("recipe not found")
		}
		return nil, err
	}

	// Increment likes
	recipe.Likes++
	recipe.UpdatedAt = time.Now()

	// Save to database
	err = r.DB.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, err
	}

	return recipe, nil
}

// AddReview adds a review to a recipe
func (r *Resolver) AddReview(ctx context.Context, recipeId string, rating int, comment *string) (*domainModels.Review, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Get recipe
	recipe, err := r.DB.GetRecipe(ctx, recipeId)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("recipe not found")
		}
		return nil, err
	}

	// Create review
	review := &domainModels.Review{
		ID:        domainModels.NewID(),
		Recipe:    recipe,
		Author:    currentUser,
		Rating:    rating,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if comment != nil {
		review.Comment = *comment
	}

	// Save to database
	err = r.DB.CreateReview(ctx, review)
	if err != nil {
		return nil, err
	}

	// Update recipe's average rating
	r.updateRecipeRating(ctx, recipe)

	return review, nil
}

// updateRecipeRating calculates and updates a recipe's average rating
func (r *RecipeResolver) updateRecipeRating(ctx context.Context, recipe *domainModels.Recipe) error {
	// Implementation would recalculate the average rating based on all reviews
	return nil
}
