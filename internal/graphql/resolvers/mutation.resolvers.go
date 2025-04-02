package resolvers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/carldunham/useful-cookery/internal/auth"
	"github.com/carldunham/useful-cookery/internal/database"
	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// Register registers a new user.
func (r *Resolver) register(ctx context.Context, input gqlmodel.RegisterInput) (*gqlmodel.AuthPayload, error) {
	user, err := r.UserResolver.registerUser(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	// token, err := r.AuthService.GenerateToken(user.ID, user.Role)
	// if err != nil {
	// 	return nil, err
	// }

	return &gqlmodel.AuthPayload{
		// Token: token,
		User: user,
	}, nil
}

// Login authenticates a user.
func (r *Resolver) login(ctx context.Context, email string, password string) (*gqlmodel.AuthPayload, error) {
	token, user, err := r.AuthService.Login(ctx, email, password)
	if err != nil {
		return nil, fmt.Errorf("failed to log in: %w", err)
	}

	return &gqlmodel.AuthPayload{
		Token: token,
		User:  user,
	}, nil
}

// RegisterUser registers a new user.
func (r *UserResolver) registerUser(ctx context.Context, input gqlmodel.RegisterInput) (*domainmodel.User, error) {
	user, err := r.AuthService.RegisterUser(ctx, input.Name, input.Email, input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	return user, nil
}

// UpdateUser updates the current user's information.
func (r *Resolver) updateUser(ctx context.Context, input gqlmodel.UpdateUserInput) (*domainmodel.User, error) {
	return r.UserResolver.updateUser(ctx, input)
}

// UpdateUser updates a user's information.
func (r *UserResolver) updateUser(ctx context.Context, input gqlmodel.UpdateUserInput) (*domainmodel.User, error) {
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
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		currentUser.Password = hashedPassword
	}

	// Update preferences if provided
	if input.Preferences != nil {
		if currentUser.Preferences == nil {
			currentUser.Preferences = &domainmodel.UserPreferences{}
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
			prefs.SkillLevel = input.Preferences.SkillLevel.String()
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
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return currentUser, nil
}

// CreateRecipe creates a new recipe.
func (r *Resolver) createRecipe(ctx context.Context, input gqlmodel.RecipeInput) (*domainmodel.Recipe, error) {
	return r.RecipeResolver.createRecipe(ctx, input)
}

// CreateRecipe creates a new recipe.
func (r *RecipeResolver) createRecipe(ctx context.Context, input gqlmodel.RecipeInput) (*domainmodel.Recipe, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Create recipe
	recipe := &domainmodel.Recipe{
		Title:       input.Title,
		Description: input.Description,
		Author:      currentUser,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		Ingredients: []domainmodel.DetailedIngredient{},
		Steps:       []domainmodel.Step{},
	}

	// Set optional fields
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
		recipe.Difficulty = string(*input.Difficulty)
	}
	if len(input.Tags) > 0 {
		recipe.Tags = input.Tags
	}

	// Process ingredients
	for _, ingInput := range input.Ingredients {
		ingredient := domainmodel.DetailedIngredient{
			ID:   domainmodel.NewID(),
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
		step := domainmodel.Step{
			ID:          domainmodel.NewID(),
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
		categories := make([]domainmodel.Category, 0, len(input.CategoryIDs))
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
		return nil, fmt.Errorf("failed to create recipe: %w", err)
	}

	return recipe, nil
}

// UpdateRecipe updates an existing recipe.
func (r *Resolver) updateRecipe(ctx context.Context, id string, input gqlmodel.RecipeInput) (*domainmodel.Recipe, error) {
	return r.RecipeResolver.updateRecipe(ctx, id, input)
}

// UpdateRecipe updates an existing recipe.
func (r *RecipeResolver) updateRecipe(ctx context.Context, id string, input gqlmodel.RecipeInput) (*domainmodel.Recipe, error) {
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
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Check if user has permission to update (author or admin)
	if recipe.Author != nil && recipe.Author.ID != currentUser.ID && currentUser.Role != domainmodel.AdminRole {
		return nil, errors.New("permission denied")
	}

	// Update basic fields
	recipe.Title = input.Title
	recipe.Description = input.Description
	recipe.Difficulty = input.Difficulty.String()
	recipe.UpdatedAt = time.Now()

	// Update optional fields
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
	if len(input.Tags) > 0 {
		recipe.Tags = input.Tags
	}

	// Process ingredients (replace all)
	recipe.Ingredients = []domainmodel.DetailedIngredient{}
	for _, ingInput := range input.Ingredients {
		ingredient := domainmodel.DetailedIngredient{
			ID:   domainmodel.NewID(),
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
	recipe.Steps = []domainmodel.Step{}
	for _, stepInput := range input.Steps {
		step := domainmodel.Step{
			ID:          domainmodel.NewID(),
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
		categories := make([]domainmodel.Category, 0, len(input.CategoryIDs))
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
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

// DeleteRecipe deletes a recipe.
func (r *Resolver) deleteRecipe(ctx context.Context, id string) (bool, error) {
	return r.RecipeResolver.deleteRecipe(ctx, id)
}

// DeleteRecipe deletes a recipe.
func (r *RecipeResolver) deleteRecipe(ctx context.Context, id string) (bool, error) {
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
		return false, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Check if user has permission to delete (author or admin)
	if recipe.Author != nil && recipe.Author.ID != currentUser.ID && currentUser.Role != domainmodel.AdminRole {
		return false, errors.New("permission denied")
	}

	// Delete from database
	err = r.DB.DeleteRecipe(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete recipe: %w", err)
	}

	return true, nil
}

// LikeRecipe adds a like to a recipe.
func (r *Resolver) likeRecipe(ctx context.Context, id string) (*domainmodel.Recipe, error) {
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
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Increment likes
	recipe.Likes++
	recipe.UpdatedAt = time.Now()

	// Save to database
	err = r.DB.UpdateRecipe(ctx, recipe)
	if err != nil {
		return nil, fmt.Errorf("failed to update recipe: %w", err)
	}

	return recipe, nil
}

// SaveRecipe saves a recipe to the user's saved recipes list.
func (r *Resolver) saveRecipe(ctx context.Context, id string) (*domainmodel.User, error) {
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
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Check if recipe is already saved
	for _, savedRecipe := range currentUser.SavedRecipes {
		if savedRecipe.ID == recipe.ID {
			return currentUser, nil // Recipe already saved
		}
	}

	// Add recipe to saved recipes
	currentUser.SavedRecipes = append(currentUser.SavedRecipes, *recipe)
	currentUser.UpdatedAt = time.Now()

	// Save to database
	err = r.DB.UpdateUser(ctx, currentUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return currentUser, nil
}

// UnsaveRecipe removes a recipe from the user's saved recipes list.
func (r *Resolver) unsaveRecipe(ctx context.Context, id string) (*domainmodel.User, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Find and remove recipe from saved recipes
	savedRecipes := make([]domainmodel.Recipe, 0, len(currentUser.SavedRecipes))
	recipeFound := false

	for _, savedRecipe := range currentUser.SavedRecipes {
		if savedRecipe.ID != id {
			savedRecipes = append(savedRecipes, savedRecipe)
		} else {
			recipeFound = true
		}
	}

	if !recipeFound {
		return nil, errors.New("recipe not found in saved recipes")
	}

	currentUser.SavedRecipes = savedRecipes
	currentUser.UpdatedAt = time.Now()

	// Save to database
	err := r.DB.UpdateUser(ctx, currentUser)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return currentUser, nil
}

// AddReview adds a review to a recipe.
func (r *Resolver) addReview(ctx context.Context, recipeID string, rating int, comment *string) (*domainmodel.Review, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Get recipe
	recipe, err := r.DB.GetRecipe(ctx, recipeID)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("recipe not found")
		}
		return nil, fmt.Errorf("failed to get recipe: %w", err)
	}

	// Create review
	review := &domainmodel.Review{
		ID:        domainmodel.NewID(),
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
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	// Update recipe's average rating
	if err := r.RecipeResolver.updateRecipeRating(ctx, recipe); err != nil {
		return nil, fmt.Errorf("updating recipe rating: %w", err)
	}

	return review, nil
}

// UpdateReview updates an existing review.
func (r *Resolver) updateReview(ctx context.Context, id string, rating *int, comment *string) (*domainmodel.Review, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Get review
	review, err := r.DB.GetReview(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("review not found")
		}
		return nil, fmt.Errorf("failed to get review: %w", err)
	}

	// Check if user has permission to update (author or admin)
	if review.Author != nil && review.Author.ID != currentUser.ID && currentUser.Role != domainmodel.AdminRole {
		return nil, errors.New("permission denied")
	}

	// Update review fields
	if rating != nil {
		review.Rating = *rating
	}
	if comment != nil {
		review.Comment = *comment
	}
	review.UpdatedAt = time.Now()

	// Save to database
	err = r.DB.UpdateReview(ctx, review)
	if err != nil {
		return nil, fmt.Errorf("failed to update review: %w", err)
	}

	// Update recipe's average rating
	if err := r.RecipeResolver.updateRecipeRating(ctx, review.Recipe); err != nil {
		return nil, fmt.Errorf("updating recipe rating: %w", err)
	}

	return review, nil
}

// DeleteReview deletes a review.
func (r *Resolver) deleteReview(ctx context.Context, id string) (bool, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return false, errors.New("not authenticated")
	}

	// Get review to check ownership
	review, err := r.DB.GetReview(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return false, errors.New("review not found")
		}
		return false, fmt.Errorf("failed to get review: %w", err)
	}

	// Check if user has permission to delete (author or admin)
	if review.Author != nil && review.Author.ID != currentUser.ID && currentUser.Role != domainmodel.AdminRole {
		return false, errors.New("permission denied")
	}

	// Delete from database
	err = r.DB.DeleteReview(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete review: %w", err)
	}

	// Update recipe's average rating
	if err := r.RecipeResolver.updateRecipeRating(ctx, review.Recipe); err != nil {
		return false, fmt.Errorf("updating recipe rating: %w", err)
	}

	return true, nil
}

// updateRecipeRating calculates and updates a recipe's average rating.
func (r *RecipeResolver) updateRecipeRating(ctx context.Context, recipe *domainmodel.Recipe) error {
	// Implementation would recalculate the average rating based on all reviews
	return nil
}

// CreateCategory creates a new category.
func (r *Resolver) createCategory(ctx context.Context, name string, description *string) (*domainmodel.Category, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Check if user is admin
	if currentUser.Role != domainmodel.AdminRole {
		return nil, errors.New("permission denied")
	}

	// Create category
	category := &domainmodel.Category{
		ID:        domainmodel.NewID(),
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if description != nil {
		category.Description = *description
	}

	// Save to database
	err := r.DB.CreateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return category, nil
}

// UpdateCategory updates an existing category.
func (r *Resolver) updateCategory(ctx context.Context, id string, name *string, description *string) (*domainmodel.Category, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return nil, errors.New("not authenticated")
	}

	// Check if user is admin
	if currentUser.Role != domainmodel.AdminRole {
		return nil, errors.New("permission denied")
	}

	// Get category
	category, err := r.DB.GetCategory(ctx, id)
	if err != nil {
		if errors.Is(err, database.ErrNotFound) {
			return nil, errors.New("category not found")
		}
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Update category fields
	if name != nil {
		category.Name = *name
	}
	if description != nil {
		category.Description = *description
	}
	category.UpdatedAt = time.Now()

	// Save to database
	err = r.DB.UpdateCategory(ctx, category)
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category.
func (r *Resolver) deleteCategory(ctx context.Context, id string) (bool, error) {
	// Get current user from context
	currentUser := auth.GetUserFromContext(ctx)
	if currentUser == nil {
		return false, errors.New("not authenticated")
	}

	// Check if user is admin
	if currentUser.Role != domainmodel.AdminRole {
		return false, errors.New("permission denied")
	}

	// Delete from database
	err := r.DB.DeleteCategory(ctx, id)
	if err != nil {
		return false, fmt.Errorf("failed to delete category: %w", err)
	}

	return true, nil
}
