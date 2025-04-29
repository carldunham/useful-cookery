package model

import (
	"time"

	"github.com/dgraph-io/dgo/v240/protos/api"
	"github.com/google/uuid"
)

// User represents a user in the system.
type User struct {
	ID             string           `dgraph:"uid"            json:"uid"`
	Name           string           `dgraph:"name"           json:"name"`
	Email          string           `dgraph:"email"          json:"email"`
	Password       string           `dgraph:"password"       json:"-"`
	Role           Role             `dgraph:"role"           json:"role"`
	Preferences    *UserPreferences `dgraph:"preferences"    json:"preferences,omitempty"`
	SavedRecipes   []Recipe         `dgraph:"savedRecipes"   json:"savedRecipes,omitempty"`
	CreatedRecipes []Recipe         `dgraph:"createdRecipes" json:"createdRecipes,omitempty"`
	Reviews        []Review         `dgraph:"reviews"        json:"reviews,omitempty"`
	CreatedAt      time.Time        `dgraph:"createdAt"      json:"createdAt"`
	UpdatedAt      time.Time        `dgraph:"updatedAt"      json:"updatedAt"`
}

type Role string

const (
	AdminRole Role = "ADMIN"
	UserRole  Role = "USER"
)

// UserPreferences represents a user's preferences.
type UserPreferences struct {
	DietaryRestrictions []string `dgraph:"dietaryRestrictions" json:"dietaryRestrictions,omitempty"`
	FavoriteIngredients []string `dgraph:"favoriteIngredients" json:"favoriteIngredients,omitempty"`
	DislikedIngredients []string `dgraph:"dislikedIngredients" json:"dislikedIngredients,omitempty"`
	SkillLevel          string   `dgraph:"skillLevel"          json:"skillLevel,omitempty"`
	CuisinePreferences  []string `dgraph:"cuisinePreferences"  json:"cuisinePreferences,omitempty"`
}

// Category represents a recipe category.
type Category struct {
	ID          string    `dgraph:"uid"         json:"uid"`
	Name        string    `dgraph:"name"        json:"name"`
	Description string    `dgraph:"description" json:"description,omitempty"`
	Recipes     []Recipe  `dgraph:"recipes"     json:"recipes,omitempty"`
	CreatedAt   time.Time `dgraph:"createdAt"   json:"createdAt"`
	UpdatedAt   time.Time `dgraph:"updatedAt"   json:"updatedAt"`
}

// Recipe represents a recipe.
type Recipe struct {
	ID             string               `dgraph:"uid"            json:"uid"`
	OriginalID     string               `dgraph:"originalID"     json:"originalID,omitempty"`
	Title          string               `dgraph:"title"          json:"title"`
	Description    string               `dgraph:"description"    json:"description,omitempty"`
	Notes          string               `dgraph:"notes"          json:"notes,omitempty"`
	Author         *User                `dgraph:"author"         json:"author,omitempty"`
	Categories     []Category           `dgraph:"categories"     json:"categories,omitempty"`
	Cuisine        string               `dgraph:"cuisine"        json:"cuisine,omitempty"`
	PrepTime       int                  `dgraph:"prepTime"       json:"prepTime,omitempty"`
	CookTime       int                  `dgraph:"cookTime"       json:"cookTime,omitempty"`
	Servings       int                  `dgraph:"servings"       json:"servings,omitempty"`
	Difficulty     string               `dgraph:"difficulty"     json:"difficulty,omitempty"`
	DifficultyText string               `dgraph:"difficultyText" json:"difficultyText,omitempty"`
	Ingredients    []DetailedIngredient `dgraph:"ingredients"    json:"ingredients,omitempty"`
	Steps          []Step               `dgraph:"steps"          json:"steps,omitempty"`
	NutritionInfo  *NutritionInfo       `dgraph:"nutritionInfo"  json:"nutritionInfo,omitempty"`
	Images         []Image              `dgraph:"images"         json:"images,omitempty"`
	Tags           []string             `dgraph:"tags"           json:"tags,omitempty"`
	Likes          int                  `dgraph:"likes"          json:"likes,omitempty"`
	Reviews        []Review             `dgraph:"reviews"        json:"reviews,omitempty"`
	AverageRating  float64              `dgraph:"averageRating"  json:"averageRating,omitempty"`
	SavedBy        []User               `dgraph:"savedBy"        json:"savedBy,omitempty"`
	Embeddings     []float32            `dgraph:"embeddings"     json:"-"`
	CreatedAt      time.Time            `dgraph:"createdAt"      json:"createdAt"`
	UpdatedAt      time.Time            `dgraph:"updatedAt"      json:"updatedAt"`
}

// IngredientUnit represents a measurement unit for an ingredient.
type IngredientUnit struct {
	ID     string  `dgraph:"uid"    json:"uid"`
	System string  `dgraph:"system" json:"system"` // "imperial", "metric", "weight", etc.
	Value  float64 `dgraph:"value"  json:"value"`  // Numeric value
	Unit   string  `dgraph:"unit"   json:"unit"`   // The unit name (cup, g, oz, etc.)
	IsMain bool    `dgraph:"isMain" json:"isMain"` // Whether this is the main unit for display
}

// DetailedIngredient represents a detailed ingredient in a recipe.
type DetailedIngredient struct {
	ID          string           `dgraph:"uid"         json:"uid"`
	Name        string           `dgraph:"name"        json:"name"`
	Units       []IngredientUnit `dgraph:"units"       json:"units,omitempty"`
	Quantity    float64          `dgraph:"quantity"    json:"quantity,omitempty"` // Deprecated field
	Unit        string           `dgraph:"unit"        json:"unit,omitempty"`     // Deprecated field
	Preparation string           `dgraph:"preparation" json:"preparation,omitempty"`
	Substitutes []string         `dgraph:"substitutes" json:"substitutes,omitempty"`
	IsOptional  bool             `dgraph:"isOptional"  json:"isOptional,omitempty"`
}

// Step represents a step in a recipe.
type Step struct {
	ID           string `dgraph:"uid"          json:"uid"`
	OrderIndex   int    `dgraph:"orderIndex"   json:"orderIndex"`
	Description  string `dgraph:"description"  json:"description"`
	Image        *Image `dgraph:"image"        json:"image,omitempty"`
	TimeEstimate int    `dgraph:"timeEstimate" json:"timeEstimate,omitempty"`
}

// NutritionInfo represents nutritional information for a recipe.
type NutritionInfo struct {
	Calories int     `dgraph:"calories" json:"calories,omitempty"`
	Protein  float64 `dgraph:"protein"  json:"protein,omitempty"`
	Carbs    float64 `dgraph:"carbs"    json:"carbs,omitempty"`
	Fat      float64 `dgraph:"fat"      json:"fat,omitempty"`
	Fiber    float64 `dgraph:"fiber"    json:"fiber,omitempty"`
	Sugar    float64 `dgraph:"sugar"    json:"sugar,omitempty"`
	Sodium   int     `dgraph:"sodium"   json:"sodium,omitempty"`
}

// Image represents an image.
type Image struct {
	ID     string `dgraph:"uid"    json:"uid"`
	URL    string `dgraph:"url"    json:"url"`
	Alt    string `dgraph:"alt"    json:"alt,omitempty"`
	Width  int    `dgraph:"width"  json:"width,omitempty"`
	Height int    `dgraph:"height" json:"height,omitempty"`
}

// Review represents a review for a recipe.
type Review struct {
	ID        string    `dgraph:"uid"       json:"uid"`
	Recipe    *Recipe   `dgraph:"recipe"    json:"recipe"`
	Author    *User     `dgraph:"author"    json:"author"`
	Rating    int       `dgraph:"rating"    json:"rating"`
	Comment   string    `dgraph:"comment"   json:"comment,omitempty"`
	CreatedAt time.Time `dgraph:"createdAt" json:"createdAt"`
	UpdatedAt time.Time `dgraph:"updatedAt" json:"updatedAt"`
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

// DgraphQuery builds a DGraph query.
func DgraphQuery(q string, vars map[string]string) *api.Request {
	req := &api.Request{
		Query: q,
	}

	if vars != nil {
		req.Vars = vars
	}

	return req
}
