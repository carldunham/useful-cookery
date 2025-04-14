package troff

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParse_BasicRecipe(t *testing.T) {
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "TEST RECIPE" "A simple test recipe"
This is a test recipe description.
It has multiple lines.
.IH "4 servings"
.IG "1 cup" "sugar" "200g"
.IG "2" "eggs"
.PH
.SK 1
Mix ingredients.
.SK 2
Bake for 30 minutes.
.NX
Some notes about the recipe.
.WR
John Doe
john@example.com
Test Organization, Test City`

	recipe, err := Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check ID
	assert.Equal(t, "RECIPE-ID", recipe.ID, "Recipe ID should match")

	// Check title
	assert.Equal(t, "TEST RECIPE", recipe.Title, "Recipe title should match")

	// Check author
	require.NotNil(t, recipe.Author, "Recipe author should not be nil")
	assert.Contains(t, recipe.Author.Name, "John Doe", "Author name should contain 'John Doe'")

	// Check description
	assert.Contains(t, recipe.Description, "A simple test recipe", "Description should contain RZ description")
	assert.Contains(t, recipe.Description, "This is a test recipe description", "Description should contain introductory text")

	// Check categories - D is for Dessert
	assert.Len(t, recipe.Categories, 1, "Should have 1 category")
	assert.Equal(t, "Dessert", recipe.Categories[0].Name, "Category should be 'Dessert'")

	// Check ingredients
	assert.Len(t, recipe.Ingredients, 2, "Should have 2 ingredients")
	assert.Equal(t, "sugar", recipe.Ingredients[0].Name, "First ingredient should be 'sugar'")
	assert.Contains(t, recipe.Ingredients[0].Unit, "1 cup", "First ingredient unit should contain '1 cup'")
	assert.Equal(t, "eggs", recipe.Ingredients[1].Name, "Second ingredient should be 'eggs'")
	assert.Contains(t, recipe.Ingredients[1].Unit, "2", "Second ingredient unit should contain '2'")

	// Check steps
	assert.Len(t, recipe.Steps, 2, "Should have 2 steps")
	assert.Equal(t, 0, recipe.Steps[0].OrderIndex, "First step should have index 0")
	assert.Equal(t, "Mix ingredients.", recipe.Steps[0].Description, "First step description should match")
	assert.Equal(t, 1, recipe.Steps[1].OrderIndex, "Second step should have index 1")
	assert.Equal(t, "Bake for 30 minutes.", recipe.Steps[1].Description, "Second step description should match")
}

func TestParse_RealRecipe(t *testing.T) {
	// Using a simplified version of the Advokaat recipe
	input := `.RH MOD.RECIPES-SOURCE ADVOKAAT L "22 Dec 83"
.RZ "ADVOKAAT" "Dutch egg cognac"
Advokaat is the Dutch word for "egg cognac".
It is highly recommended for A. I. (Alcohol Imbibing) meetings.
.IH "1 bottle" "750 ml"
.IG "1 1/2 cups" "sugar" "300 g"
.IG "2 Tbsp" "vanilla sugar" "25 g"
.IG "2 cups" "milk" "500 ml"
.IG "9" "egg yolks"
.IG "1 1/2 cups" "95% grain alcohol" "350 ml"
.PH
.SK 1
Mix sugars.
.SK 2
Boil milk with half of sugars for two minutes.
.SK 3
Mix well the yolks with the other half of the sugars.
.SK 4
Using mixer add to milk. Then add alcohol slowly.
.SK 5
Bottle and let rest for two weeks to let the mixture thicken.
.NX
This recipe is from the Netherlands.
.WR
Laurent Siklossy
laurent@example.com
Test University, Amsterdam`

	recipe, err := Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check ID
	assert.Equal(t, "ADVOKAAT", recipe.ID, "Recipe ID should match")

	// Check title
	assert.Equal(t, "ADVOKAAT", recipe.Title, "Recipe title should match")

	// Check author
	require.NotNil(t, recipe.Author, "Recipe author should not be nil")
	assert.Contains(t, recipe.Author.Name, "Laurent Siklossy", "Author name should contain 'Laurent Siklossy'")

	// Check description
	assert.Contains(t, recipe.Description, "Dutch egg cognac", "Description should contain RZ description")
	assert.Contains(t, recipe.Description, "Advokaat is the Dutch word", "Description should contain introductory text")

	// Check categories - L is for Beverage (Liquid)
	assert.Len(t, recipe.Categories, 1, "Should have 1 category")
	assert.Equal(t, "Beverage", recipe.Categories[0].Name, "Category should be 'Beverage'")

	// Check ingredients
	expectedIngredients := []struct {
		name             string
		quantityContains string
	}{
		{"sugar", "1 1/2 cups"},
		{"vanilla sugar", "2 Tbsp"},
		{"milk", "2 cups"},
		{"egg yolks", "9"},
		{"95% grain alcohol", "1 1/2 cups"},
	}

	assert.Len(t, recipe.Ingredients, len(expectedIngredients), "Should have correct number of ingredients")

	for i, ing := range expectedIngredients {
		assert.Equal(t, ing.name, recipe.Ingredients[i].Name, "Ingredient name should match")
		assert.Contains(t, recipe.Ingredients[i].Unit, ing.quantityContains, "Ingredient unit should contain quantity")
	}

	// Check steps
	assert.Len(t, recipe.Steps, 5, "Should have 5 steps")
	// Just check the first and last steps
	assert.Equal(t, "Mix sugars.", recipe.Steps[0].Description, "First step description should match")
	assert.Equal(t, "Bottle and let rest for two weeks to let the mixture thicken.", recipe.Steps[4].Description, "Last step description should match")
}

func TestParse_EmptyInput(t *testing.T) {
	recipe, err := Parse(strings.NewReader(""))
	require.NoError(t, err, "Parse() should not return an error for empty input")

	// Check that we get an empty recipe
	assert.Empty(t, recipe.Title, "Title should be empty")
	assert.Empty(t, recipe.Ingredients, "Ingredients should be empty")
	assert.Empty(t, recipe.Steps, "Steps should be empty")
}

func TestParse_InvalidInput(t *testing.T) {
	// Test with a nil reader
	_, err := Parse(nil)
	assert.Error(t, err, "Parse() should return an error with nil reader")

	// Test with malformed input (missing closing quotes)
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "TEST RECIPE" "A simple test recipe"
.IG "1 cup sugar`
	_, err = Parse(strings.NewReader(input))
	// This should not error, but should handle the unclosed quote gracefully
	assert.NoError(t, err, "Parse() should handle unclosed quote gracefully")
}
