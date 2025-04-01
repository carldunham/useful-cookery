package models

// RegisterInput represents the input for user registration
type RegisterInput struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginInput represents the input for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserInput represents the input for updating a user
type UpdateUserInput struct {
	Name        *string               `json:"name,omitempty"`
	Email       *string               `json:"email,omitempty"`
	Password    *string               `json:"password,omitempty"`
	Preferences *UserPreferencesInput `json:"preferences,omitempty"`
}

// UserPreferencesInput represents the input for user preferences
type UserPreferencesInput struct {
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`
	FavoriteIngredients []string `json:"favoriteIngredients,omitempty"`
	DislikedIngredients []string `json:"dislikedIngredients,omitempty"`
	SkillLevel          *string  `json:"skillLevel,omitempty"`
	CuisinePreferences  []string `json:"cuisinePreferences,omitempty"`
}

// RecipeInput represents the input for creating or updating a recipe
type RecipeInput struct {
	Title         string              `json:"title"`
	Description   *string             `json:"description,omitempty"`
	CategoryIDs   []string            `json:"categoryIds,omitempty"`
	Cuisine       *string             `json:"cuisine,omitempty"`
	PrepTime      *int                `json:"prepTime,omitempty"`
	CookTime      *int                `json:"cookTime,omitempty"`
	Servings      *int                `json:"servings,omitempty"`
	Difficulty    *string             `json:"difficulty,omitempty"`
	Ingredients   []IngredientInput   `json:"ingredients"`
	Steps         []StepInput         `json:"steps"`
	NutritionInfo *NutritionInfoInput `json:"nutritionInfo,omitempty"`
	Tags          []string            `json:"tags,omitempty"`
	Images        []ImageInput        `json:"images,omitempty"`
}

// RecipeFilter represents the filter for recipe queries
type RecipeFilter struct {
	Search      *string  `json:"search,omitempty"`
	Categories  []string `json:"categories,omitempty"`
	Cuisine     *string  `json:"cuisine,omitempty"`
	Difficulty  *string  `json:"difficulty,omitempty"`
	MaxPrepTime *int     `json:"maxPrepTime,omitempty"`
	Ingredients []string `json:"ingredients,omitempty"`
	AuthorID    *string  `json:"authorId,omitempty"`
}

// RecipeOrder represents the ordering options for recipe queries
type RecipeOrder struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

// IngredientInput represents the input for a recipe ingredient
type IngredientInput struct {
	Name        string   `json:"name"`
	Quantity    *float64 `json:"quantity,omitempty"`
	Unit        *string  `json:"unit,omitempty"`
	Preparation *string  `json:"preparation,omitempty"`
	Substitutes []string `json:"substitutes,omitempty"`
	IsOptional  *bool    `json:"isOptional,omitempty"`
}

// StepInput represents the input for a recipe step
type StepInput struct {
	OrderIndex   int         `json:"orderIndex"`
	Description  string      `json:"description"`
	Image        *ImageInput `json:"image,omitempty"`
	TimeEstimate *int        `json:"timeEstimate,omitempty"`
}

// NutritionInfoInput represents the input for nutrition information
type NutritionInfoInput struct {
	Calories *int     `json:"calories,omitempty"`
	Protein  *float64 `json:"protein,omitempty"`
	Carbs    *float64 `json:"carbs,omitempty"`
	Fat      *float64 `json:"fat,omitempty"`
	Fiber    *float64 `json:"fiber,omitempty"`
	Sugar    *float64 `json:"sugar,omitempty"`
	Sodium   *int     `json:"sodium,omitempty"`
}

// ImageInput represents the input for an image
type ImageInput struct {
	URL    string  `json:"url"`
	Alt    *string `json:"alt,omitempty"`
	Width  *int    `json:"width,omitempty"`
	Height *int    `json:"height,omitempty"`
}

// AuthPayload represents the authentication response
type AuthPayload struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

// SearchParams represents parameters for AI search
type SearchParams struct {
	Ingredients         []string `json:"ingredients,omitempty"`
	ExcludedIngredients []string `json:"excludedIngredients,omitempty"`
	Categories          []string `json:"categories,omitempty"`
	Cuisine             string   `json:"cuisine,omitempty"`
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`
	MaxPrepTime         int      `json:"maxPrepTime,omitempty"`
	Difficulty          string   `json:"difficulty,omitempty"`
}
