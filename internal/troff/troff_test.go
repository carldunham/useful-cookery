package troff_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carldunham/useful-cookery/internal/troff"
)

//nolint:funlen // Test function with many assertions to verify parsing correctness.
func TestParse_BasicRecipe(t *testing.T) {
	t.Parallel()
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

	recipe, err := troff.Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check OriginalID
	assert.Equal(t, "RECIPE-ID", recipe.OriginalID, "Recipe OriginalID should match")

	// Check title
	assert.Equal(t, "TEST RECIPE", recipe.Title, "Recipe title should match")

	// Check author
	require.NotNil(t, recipe.Author, "Recipe author should not be nil")
	assert.Contains(t, recipe.Author.Name, "John Doe", "Author name should contain 'John Doe'")

	// Check description
	assert.Contains(t, recipe.Description, "A simple test recipe", "Description should contain RZ description")
	assert.Contains(t, recipe.Description, "This is a test recipe description",
		"Description should contain introductory text")

	// Check notes
	assert.Equal(t, "Some notes about the recipe.", recipe.Notes, "Notes should be stored in the Notes field")

	// Check categories - D is for Dessert
	assert.Len(t, recipe.Categories, 1, "Should have 1 category")
	assert.Equal(t, "Dessert", recipe.Categories[0].Name, "Category should be 'Dessert'")

	// Check ingredients
	assert.Len(t, recipe.Ingredients, 2, "Should have 2 ingredients")

	// Check first ingredient (sugar with imperial and metric units)
	assert.Equal(t, "sugar", recipe.Ingredients[0].Name, "First ingredient should be 'sugar'")
	assert.Contains(t, recipe.Ingredients[0].Unit, "1 cup", "First ingredient unit should contain '1 cup'")
	assert.InEpsilon(t, 1.0, recipe.Ingredients[0].Quantity, 0.001, "First ingredient quantity should be 1.0")

	// Check units for first ingredient
	assert.Len(t, recipe.Ingredients[0].Units, 2, "First ingredient should have 2 units")
	assert.Equal(t, "imperial", recipe.Ingredients[0].Units[0].System, "First unit should be imperial")
	assert.Equal(t, "1 cup", recipe.Ingredients[0].Units[0].Unit, "First unit should be '1 cup'")
	assert.InEpsilon(t, 1.0, recipe.Ingredients[0].Units[0].Value, 0.001, "First unit value should be 1.0")
	assert.True(t, recipe.Ingredients[0].Units[0].IsMain, "First unit should be marked as main")

	assert.Equal(t, "metric", recipe.Ingredients[0].Units[1].System, "Second unit should be metric")
	assert.Equal(t, "200g", recipe.Ingredients[0].Units[1].Unit, "Second unit should be '200g'")
	assert.False(t, recipe.Ingredients[0].Units[1].IsMain, "Second unit should not be marked as main")

	// Check second ingredient (eggs with only imperial unit)
	assert.Equal(t, "eggs", recipe.Ingredients[1].Name, "Second ingredient should be 'eggs'")
	assert.Contains(t, recipe.Ingredients[1].Unit, "2", "Second ingredient unit should contain '2'")
	assert.InEpsilon(t, 2.0, recipe.Ingredients[1].Quantity, 0.001, "Second ingredient quantity should be 2.0")

	// Check units for second ingredient
	assert.Len(t, recipe.Ingredients[1].Units, 1, "Second ingredient should have 1 unit")
	assert.Equal(t, "imperial", recipe.Ingredients[1].Units[0].System, "Unit should be imperial")
	assert.Equal(t, "2", recipe.Ingredients[1].Units[0].Unit, "Unit should be '2'")
	assert.InEpsilon(t, 2.0, recipe.Ingredients[1].Units[0].Value, 0.001, "Unit value should be 2.0")
	assert.True(t, recipe.Ingredients[1].Units[0].IsMain, "Unit should be marked as main")

	// Check steps
	assert.Len(t, recipe.Steps, 2, "Should have 2 steps")
	assert.Equal(t, 0, recipe.Steps[0].OrderIndex, "First step should have index 0")
	assert.Equal(t, "Mix ingredients.", recipe.Steps[0].Description, "First step description should match")
	assert.Equal(t, 1, recipe.Steps[1].OrderIndex, "Second step should have index 1")
	assert.Equal(t, "Bake for 30 minutes.", recipe.Steps[1].Description, "Second step description should match")
}

//nolint:funlen // TODO: simplify.
func TestParse_RealRecipe(t *testing.T) {
	t.Parallel()
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

	recipe, err := troff.Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check OriginalID
	assert.Equal(t, "ADVOKAAT", recipe.OriginalID, "Recipe OriginalID should match")

	// Check title
	assert.Equal(t, "ADVOKAAT", recipe.Title, "Recipe title should match")

	// Check author
	require.NotNil(t, recipe.Author, "Recipe author should not be nil")
	assert.Contains(t, recipe.Author.Name, "Laurent Siklossy", "Author name should contain 'Laurent Siklossy'")

	// Check description
	assert.Contains(t, recipe.Description, "Dutch egg cognac", "Description should contain RZ description")
	assert.Contains(t, recipe.Description, "Advokaat is the Dutch word", "Description should contain introductory text")

	// Check notes
	assert.Equal(t, "This recipe is from the Netherlands.", recipe.Notes, "Notes should be stored in the Notes field")

	// Check categories - L is for Beverage (Liquid)
	assert.Len(t, recipe.Categories, 1, "Should have 1 category")
	assert.Equal(t, "Beverage", recipe.Categories[0].Name, "Category should be 'Beverage'")

	// Check ingredients
	expectedIngredients := []struct {
		name          string
		imperialUnit  string
		imperialValue float64
		hasMetric     bool
		metricUnit    string
	}{
		{"sugar", "1 1/2 cups", 1.5, true, "300 g"},
		{"vanilla sugar", "2 Tbsp", 2.0, true, "25 g"},
		{"milk", "2 cups", 2.0, true, "500 ml"},
		{"egg yolks", "9", 9.0, false, ""},
		{"95% grain alcohol", "1 1/2 cups", 1.5, true, "350 ml"},
	}

	assert.Len(t, recipe.Ingredients, len(expectedIngredients), "Should have correct number of ingredients")

	for i, ing := range expectedIngredients {
		assert.Equal(t, ing.name, recipe.Ingredients[i].Name, "Ingredient name should match")
		assert.Contains(t, recipe.Ingredients[i].Unit, ing.imperialUnit, "Ingredient unit should contain imperial quantity")

		// Check imperial unit
		assert.Equal(t, "imperial", recipe.Ingredients[i].Units[0].System, "First unit should be imperial")
		assert.Equal(t, ing.imperialUnit, recipe.Ingredients[i].Units[0].Unit, "Imperial unit should match")
		assert.InEpsilon(t, ing.imperialValue, recipe.Ingredients[i].Units[0].Value, 0.001, "Imperial value should match")
		assert.True(t, recipe.Ingredients[i].Units[0].IsMain, "Imperial unit should be marked as main")

		// Check metric unit if present
		if ing.hasMetric {
			assert.Len(t, recipe.Ingredients[i].Units, 2, "Ingredient should have 2 units")
			assert.Equal(t, "metric", recipe.Ingredients[i].Units[1].System, "Second unit should be metric")
			assert.Equal(t, ing.metricUnit, recipe.Ingredients[i].Units[1].Unit, "Metric unit should match")
			assert.False(t, recipe.Ingredients[i].Units[1].IsMain, "Metric unit should not be marked as main")
		} else {
			assert.Len(t, recipe.Ingredients[i].Units, 1, "Ingredient should have 1 unit")
		}
	}

	// Check steps
	assert.Len(t, recipe.Steps, 5, "Should have 5 steps")
	// Just check the first and last steps
	assert.Equal(t, "Mix sugars.", recipe.Steps[0].Description, "First step description should match")
	assert.Equal(t, "Bottle and let rest for two weeks to let the mixture thicken.",
		recipe.Steps[4].Description, "Last step description should match")
}

// TestParse_RecipeWithRatings tests parsing a recipe with ratings information.
//
//nolint:funlen // Test function with many assertions to verify parsing correctness.
func TestParse_RecipeWithRatings(t *testing.T) {
	t.Parallel()
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "CHOCOLATE CAKE" "A delicious chocolate cake"
This is a classic chocolate cake recipe.
.SH "RATING"
<i>Difficulty</i>
Easy to moderate
<i>Time</i>
45 minutes preparation, 35 minutes baking
<i>Precision</i>
Measure carefully
.IH "8 servings"
.IG "2 cups" "flour" "250g"
.IG "1 cup" "sugar" "200g"
.IG "1/2 cup" "cocoa powder" "50g"
.PH
.SK 1
Preheat oven to 350°F.
.SK 2
Mix dry ingredients.
.SK 3
Bake for 35 minutes.
.NX
Best served warm with ice cream.
.WR
Jane Smith
jane@example.com`

	recipe, err := troff.Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check basic fields
	assert.Equal(t, "CHOCOLATE CAKE", recipe.Title, "Recipe title should match")
	assert.Equal(t, "RECIPE-ID", recipe.OriginalID, "Recipe OriginalID should match")

	// Check ratings fields
	assert.Equal(t, "Easy to moderate", recipe.Difficulty, "Difficulty should be parsed correctly")
	assert.Equal(t, 45, recipe.PrepTime, "Prep time should be parsed correctly")
	assert.Equal(t, 8, recipe.Servings, "Servings should be parsed correctly")

	// Check that precision is stored in tags
	assert.Contains(t, recipe.Tags, "Precision: Measure carefully", "Precision should be stored in tags")

	// Check ingredients with quantity parsing
	assert.Len(t, recipe.Ingredients, 3, "Should have 3 ingredients")

	// Check first ingredient with numeric quantity
	assert.Equal(t, "flour", recipe.Ingredients[0].Name, "First ingredient should be 'flour'")
	assert.InEpsilon(t, 2.0, recipe.Ingredients[0].Quantity, 0.001, "First ingredient quantity should be 2.0")

	// Check units for first ingredient
	assert.Len(t, recipe.Ingredients[0].Units, 2, "First ingredient should have 2 units")
	assert.Equal(t, "imperial", recipe.Ingredients[0].Units[0].System, "First unit should be imperial")
	assert.Equal(t, "2 cups", recipe.Ingredients[0].Units[0].Unit, "First unit should be '2 cups'")
	assert.InEpsilon(t, 2.0, recipe.Ingredients[0].Units[0].Value, 0.001, "First unit value should be 2.0")
	assert.Equal(t, "metric", recipe.Ingredients[0].Units[1].System, "Second unit should be metric")
	assert.Equal(t, "250g", recipe.Ingredients[0].Units[1].Unit, "Second unit should be '250g'")

	// Check ingredient with fraction
	assert.Equal(t, "cocoa powder", recipe.Ingredients[2].Name, "Third ingredient should be 'cocoa powder'")
	assert.InEpsilon(t, 0.5, recipe.Ingredients[2].Quantity, 0.001, "Third ingredient quantity should be 0.5")

	// Check units for third ingredient
	assert.Len(t, recipe.Ingredients[2].Units, 2, "Third ingredient should have 2 units")
	assert.Equal(t, "imperial", recipe.Ingredients[2].Units[0].System, "First unit should be imperial")
	assert.Equal(t, "1/2 cup", recipe.Ingredients[2].Units[0].Unit, "First unit should be '1/2 cup'")
	assert.InEpsilon(t, 0.5, recipe.Ingredients[2].Units[0].Value, 0.001, "First unit value should be 0.5")
	assert.Equal(t, "metric", recipe.Ingredients[2].Units[1].System, "Second unit should be metric")
	assert.Equal(t, "50g", recipe.Ingredients[2].Units[1].Unit, "Second unit should be '50g'")
}

func TestParse_EmptyInput(t *testing.T) {
	t.Parallel()
	recipe, err := troff.Parse(strings.NewReader(""))
	require.NoError(t, err, "Parse() should not return an error for empty input")

	// Check that we get an empty recipe
	assert.Empty(t, recipe.Title, "Title should be empty")
	assert.Empty(t, recipe.Ingredients, "Ingredients should be empty")
	assert.Empty(t, recipe.Steps, "Steps should be empty")
}

func TestParse_InvalidInput(t *testing.T) {
	t.Parallel()
	// Test with a nil reader
	_, err := troff.Parse(nil)
	require.Error(t, err, "Parse() should return an error with nil reader")

	// Test with malformed input (missing closing quotes)
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "TEST RECIPE" "A simple test recipe"
.IG "1 cup sugar`
	_, err = troff.Parse(strings.NewReader(input))
	// This should not error, but should handle the unclosed quote gracefully
	assert.NoError(t, err, "Parse() should handle unclosed quote gracefully")
}

// TestParse_TokenParamHandling specifically tests the fix for the bug where
// TokenParams were being incorrectly accumulated to the description for .RZ and .NX commands.
func TestParse_TokenParamHandling(t *testing.T) {
	t.Parallel()
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "TEST RECIPE" "A simple test recipe"
This is introduction text.
.IH "4 servings"
.IG "1 cup" "sugar"
.PH
.SK 1
First step.
.NX
These are notes.
.WR
Author Name`

	recipe, err := troff.Parse(strings.NewReader(input))
	require.NoError(t, err, "Parse() should not return an error")

	// Check that the description contains only what it should
	assert.Equal(t, "A simple test recipe\n\nThis is introduction text.",
		recipe.Description, "Description should only contain the RZ description and introduction text")

	// Check that notes are stored in the Notes field
	assert.Equal(t, "These are notes.", recipe.Notes,
		"Notes should be stored in the Notes field")

	// Verify that the description doesn't contain any command parameters that should have been skipped
	assert.NotContains(t, recipe.Description, "TEST RECIPE",
		"Description should not contain the title parameter from RZ command")
	assert.NotContains(t, recipe.Description, "4 servings",
		"Description should not contain the yield parameter from IH command")
	assert.NotContains(t, recipe.Description, "Author Name",
		"Description should not contain the author information from WR command")
}
