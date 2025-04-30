package model

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID             string           `json:"id"`
	Name           string           `json:"name"`
	Email          string           `json:"email"`
	Password       string           `json:"-"`
	Role           Role             `json:"role"`
	Preferences    *UserPreferences `json:"preferences,omitempty"`
	SavedRecipes   []Recipe         `json:"savedRecipes,omitempty"`
	CreatedRecipes []Recipe         `json:"createdRecipes,omitempty"`
	Reviews        []Review         `json:"reviews,omitempty"`
	CreatedAt      time.Time        `json:"createdAt"`
	UpdatedAt      time.Time        `json:"updatedAt"`
}

type Role string

const (
	AdminRole Role = "ADMIN"
	UserRole  Role = "USER"
)

// UserPreferences represents a user's preferences.
type UserPreferences struct {
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`
	FavoriteIngredients []string `json:"favoriteIngredients,omitempty"`
	DislikedIngredients []string `json:"dislikedIngredients,omitempty"`
	SkillLevel          string   `json:"skillLevel,omitempty"`
	CuisinePreferences  []string `json:"cuisinePreferences,omitempty"`
}

// Category represents a recipe category.
type Category struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description,omitempty"`
	Recipes     []Recipe  `json:"recipes,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Recipe represents a recipe.
type Recipe struct {
	ID             string               `json:"id"`
	OriginalID     string               `json:"originalID,omitempty"`
	Title          string               `json:"title"`
	Description    string               `json:"description,omitempty"`
	Notes          string               `json:"notes,omitempty"`
	Author         *User                `json:"author,omitempty"`
	Categories     []Category           `json:"categories,omitempty"`
	Cuisine        string               `json:"cuisine,omitempty"`
	PrepTime       int                  `json:"prepTime,omitempty"`
	CookTime       int                  `json:"cookTime,omitempty"`
	Servings       int                  `json:"servings,omitempty"`
	Difficulty     string               `json:"difficulty,omitempty"`
	SkillLevel     string               `json:"skillLevel,omitempty"`
	DifficultyText string               `json:"difficultyText,omitempty"`
	Ingredients    []DetailedIngredient `json:"ingredients,omitempty"`
	Steps          []Step               `json:"steps,omitempty"`
	NutritionInfo  *NutritionInfo       `json:"nutritionInfo,omitempty"`
	Images         []Image              `json:"images,omitempty"`
	Tags           []string             `json:"tags,omitempty"`
	Likes          int                  `json:"likes,omitempty"`
	Reviews        []Review             `json:"reviews,omitempty"`
	AverageRating  float64              `json:"averageRating,omitempty"`
	SavedBy        []User               `json:"savedBy,omitempty"`
	Embeddings     []float32            `json:"-"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}

// IngredientUnit represents a measurement unit for an ingredient.
type IngredientUnit struct {
	ID     string  `json:"id"`
	System string  `json:"system"` // "imperial", "metric", "weight", etc.
	Value  float64 `json:"value"`  // Numeric value
	Unit   string  `json:"unit"`   // The unit name (cup, g, oz, etc.)
	IsMain bool    `json:"isMain"` // Whether this is the main unit for display
}

// DetailedIngredient represents a detailed ingredient in a recipe.
type DetailedIngredient struct {
	ID          string           `json:"id"`
	Name        string           `json:"name"`
	Units       []IngredientUnit `json:"units,omitempty"`
	Quantity    float64          `json:"quantity,omitempty"` // Deprecated field
	Unit        string           `json:"unit,omitempty"`     // Deprecated field
	Preparation string           `json:"preparation,omitempty"`
	Substitutes []string         `json:"substitutes,omitempty"`
	IsOptional  bool             `json:"isOptional,omitempty"`
}

// Step represents a step in a recipe.
type Step struct {
	ID           string `json:"id"`
	OrderIndex   int    `json:"orderIndex"`
	Description  string `json:"description"`
	Image        *Image `json:"image,omitempty"`
	TimeEstimate int    `json:"timeEstimate,omitempty"`
}

// NutritionInfo represents nutritional information for a recipe.
type NutritionInfo struct {
	Calories int     `json:"calories,omitempty"`
	Protein  float64 `json:"protein,omitempty"`
	Carbs    float64 `json:"carbs,omitempty"`
	Fat      float64 `json:"fat,omitempty"`
	Fiber    float64 `json:"fiber,omitempty"`
	Sugar    float64 `json:"sugar,omitempty"`
	Sodium   int     `json:"sodium,omitempty"`
}

// Image represents an image.
type Image struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Alt    string `json:"alt,omitempty"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Review represents a review for a recipe.
type Review struct {
	ID        string    `json:"id"`
	Recipe    *Recipe   `json:"recipe"`
	Author    *User     `json:"author"`
	Rating    int       `json:"rating"`
	Comment   string    `json:"comment,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SearchParams represents parameters for recipe search.
type SearchParams struct {
	Ingredients         []string `json:"ingredients,omitempty"`
	ExcludedIngredients []string `json:"excludedIngredients,omitempty"`
	Categories          []string `json:"categories,omitempty"`
	Cuisine             string   `json:"cuisine,omitempty"`
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`
	MaxPrepTime         int      `json:"maxPrepTime,omitempty"`
	Difficulty          string   `json:"difficulty,omitempty"`
}

// NewID generates a new UUID string.
func NewID() string {
	return uuid.New().String()
}
