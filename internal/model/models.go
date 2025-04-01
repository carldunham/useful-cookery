package models

import (
	"time"

	"github.com/dgraph-io/dgo/v2/protos/api"
	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID              string           `json:"id" dgraph:"uid"`
	Name            string           `json:"name" dgraph:"name"`
	Email           string           `json:"email" dgraph:"email"`
	Password        string           `json:"-" dgraph:"password"`
	Role            string           `json:"role" dgraph:"role"`
	Preferences     *UserPreferences `json:"preferences,omitempty" dgraph:"preferences"`
	SavedRecipes    []Recipe         `json:"savedRecipes,omitempty" dgraph:"savedRecipes"`
	CreatedRecipes  []Recipe         `json:"createdRecipes,omitempty" dgraph:"createdRecipes"`
	Reviews         []Review         `json:"reviews,omitempty" dgraph:"reviews"`
	CreatedAt       time.Time        `json:"createdAt" dgraph:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt" dgraph:"updatedAt"`
}

// UserPreferences represents a user's preferences
type UserPreferences struct {
	DietaryRestrictions []string   `json:"dietaryRestrictions,omitempty" dgraph:"dietaryRestrictions"`
	FavoriteIngredients []string   `json:"favoriteIngredients,omitempty" dgraph:"favoriteIngredients"`
	DislikedIngredients []string   `json:"dislikedIngredients,omitempty" dgraph:"dislikedIngredients"`
	SkillLevel          string     `json:"skillLevel,omitempty" dgraph:"skillLevel"`
	CuisinePreferences  []string   `json:"cuisinePreferences,omitempty" dgraph:"cuisinePreferences"`
}

// Category represents a recipe category
type Category struct {
	ID          string    `json:"id" dgraph:"uid"`
	Name        string    `json:"name" dgraph:"name"`
	Description string    `json:"description,omitempty" dgraph:"description"`
	Recipes     []Recipe  `json:"recipes,omitempty" dgraph:"recipes"`
}

// Recipe represents a recipe
type Recipe struct {
	ID            string               `json:"id" dgraph:"uid"`
	Title         string               `json:"title" dgraph:"title"`
	Description   string               `json:"description,omitempty" dgraph:"description"`
	Author        *User                `json:"author,omitempty" dgraph:"author"`
	Categories    []Category           `json:"categories,omitempty" dgraph:"categories"`
	Cuisine       string               `json:"cuisine,omitempty" dgraph:"cuisine"`
	PrepTime      int                  `json:"prepTime,omitempty" dgraph:"prepTime"`
	CookTime      int                  `json:"cookTime,omitempty" dgraph:"cookTime"`
	Servings      int                  `json:"servings,omitempty" dgraph:"servings"`
	Difficulty    string               `json:"difficulty,omitempty" dgraph:"difficulty"`
	Ingredients   []DetailedIngredient `json:"ingredients,omitempty" dgraph:"ingredients"`
	Steps         []Step               `json:"steps,omitempty" dgraph:"steps"`
	NutritionInfo *NutritionInfo       `json:"nutritionInfo,omitempty" dgraph:"nutritionInfo"`
	Images        []Image              `json:"images,omitempty" dgraph:"images"`
	Tags          []string             `json:"tags,omitempty" dgraph:"tags"`
	Likes         int                  `json:"likes,omitempty" dgraph:"likes"`
	Reviews       []Review             `json:"reviews,omitempty" dgraph:"reviews"`
	AverageRating float64              `json:"averageRating,omitempty" dgraph:"averageRating"`
	SavedBy       []User               `json:"savedBy,omitempty" dgraph:"savedBy"`
	Embeddings    []float32            `json:"-" dgraph:"embeddings"`
	CreatedAt     time.Time            `json:"createdAt" dgraph:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt" dgraph:"updatedAt"`
}

// DetailedIngredient represents a detailed ingredient in a recipe
type DetailedIngredient struct {
	ID          string   `json:"id" dgraph:"uid"`
	Name        string   `json:"name" dgraph:"name"`
	Quantity    float64  `json:"quantity,omitempty" dgraph:"quantity"`
	Unit        string   `json:"unit,omitempty" dgraph:"unit"`
	Preparation string   `json:"preparation,omitempty" dgraph:"preparation"`
	Substitutes []string `json:"substitutes,omitempty" dgraph:"substitutes"`
	IsOptional  bool     `json:"isOptional,omitempty" dgraph:"isOptional"`
}

// Step represents a step in a recipe
type Step struct {
	ID           string `json:"id" dgraph:"uid"`
	OrderIndex   int    `json:"orderIndex" dgraph:"orderIndex"`
	Description  string `json:"description" dgraph:"description"`
	Image        *Image `json:"image,omitempty" dgraph:"image"`
	TimeEstimate int    `json:"timeEstimate,omitempty" dgraph:"timeEstimate"`
}

// NutritionInfo represents nutritional information for a recipe
type NutritionInfo struct {
	Calories int     `json:"calories,omitempty" dgraph:"calories"`
	Protein  float64 `json:"protein,omitempty" dgraph:"protein"`
	Carbs    float64 `json:"carbs,omitempty" dgraph:"carbs"`
	Fat      float64 `json:"fat,omitempty" dgraph:"fat"`
	Fiber    float64 `json:"fiber,omitempty" dgraph:"fiber"`
	Sugar    float64 `json:"sugar,omitempty" dgraph:"sugar"`
	Sodium   int     `json:"sodium,omitempty" dgraph:"sodium"`
}

// Image represents an image
type Image struct {
	ID     string `json:"id" dgraph:"uid"`
	URL    string `json:"url" dgraph:"url"`
	Alt    string `json:"alt,omitempty" dgraph:"alt"`
	Width  int    `json:"width,omitempty" dgraph:"width"`
	Height int    `json:"height,omitempty" dgraph:"height"`
}

// Review represents a review for a recipe
type Review struct {
	ID        string    `json:"id" dgraph:"uid"`
	Recipe    *Recipe   `json:"recipe" dgraph:"recipe"`
	Author    *User     `json:"author" dgraph:"author"`
	Rating    int       `json:"rating" dgraph:"rating"`
	Comment   string    `json:"comment,omitempty" dgraph:"comment"`
	CreatedAt time.Time `json:"createdAt" dgraph:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" dgraph:"updatedAt"`
}

// SearchParams represents parameters for recipe search
type SearchParams struct {
	Ingredients         []string `json:"ingredients,omitempty"`
	ExcludedIngredients []string `json:"excludedIngredients,omitempty"`
	Categories          []string `json:"categories,omitempty"`
	Cuisine             string   `json:"cuisine,omitempty"`
	DietaryRestrictions []string `json:"dietaryRestrictions,omitempty"`
	MaxPrepTime         int      `json:"maxPrepTime,omitempty"`
	Difficulty          string   `json:"difficulty,omitempty"`
}

// NewID generates a new UUID string
func NewID() string {
	return uuid.New().String()
}

// DgraphQuery builds a DGraph query
func DgraphQuery(q string, vars map[string]string) *api.Request {
	req := &api.Request{
		Query: q,
	}
	
	if vars != nil {
		req.Vars = vars
	}
	
	return req
}
