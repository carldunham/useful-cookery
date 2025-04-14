package troff

import (
	"strings"
	"testing"
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
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check ID
	if recipe.ID != "RECIPE-ID" {
		t.Errorf("Expected ID 'RECIPE-ID', got '%s'", recipe.ID)
	}

	// Check title
	if recipe.Title != "TEST RECIPE" {
		t.Errorf("Expected title 'TEST RECIPE', got '%s'", recipe.Title)
	}

	// Check author
	if recipe.Author == nil || !strings.Contains(recipe.Author.Name, "John Doe") {
		t.Errorf("Expected author containing 'John Doe', got '%v'", recipe.Author)
	}

	// Check description
	if !strings.Contains(recipe.Description, "A simple test recipe") ||
		!strings.Contains(recipe.Description, "This is a test recipe description") {
		t.Errorf("Expected description to contain both the RZ description and introductory text, got '%s'", recipe.Description)
	}

	// Check categories - D is for Dessert
	if len(recipe.Categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(recipe.Categories))
	} else {
		if recipe.Categories[0].Name != "Dessert" {
			t.Errorf("Expected category 'Dessert', got '%s'", recipe.Categories[0].Name)
		}
	}

	// Check ingredients
	if len(recipe.Ingredients) != 2 {
		t.Errorf("Expected 2 ingredients, got %d", len(recipe.Ingredients))
	} else {
		if recipe.Ingredients[0].Name != "sugar" || !strings.Contains(recipe.Ingredients[0].Unit, "1 cup") {
			t.Errorf("Expected first ingredient 'sugar' with unit containing '1 cup', got '%s' with unit '%s'",
				recipe.Ingredients[0].Name, recipe.Ingredients[0].Unit)
		}
		if recipe.Ingredients[1].Name != "eggs" || !strings.Contains(recipe.Ingredients[1].Unit, "2") {
			t.Errorf("Expected second ingredient 'eggs' with unit containing '2', got '%s' with unit '%s'",
				recipe.Ingredients[1].Name, recipe.Ingredients[1].Unit)
		}
	}

	// Check steps
	if len(recipe.Steps) != 2 {
		t.Errorf("Expected 2 steps, got %d", len(recipe.Steps))
	} else {
		if recipe.Steps[0].OrderIndex != 0 || recipe.Steps[0].Description != "Mix ingredients." {
			t.Errorf("Expected first step with index 0 and description 'Mix ingredients.', got index %d and description '%s'",
				recipe.Steps[0].OrderIndex, recipe.Steps[0].Description)
		}
		if recipe.Steps[1].OrderIndex != 1 || recipe.Steps[1].Description != "Bake for 30 minutes." {
			t.Errorf("Expected second step with index 1 and description 'Bake for 30 minutes.', got index %d and description '%s'",
				recipe.Steps[1].OrderIndex, recipe.Steps[1].Description)
		}
	}
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
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check ID
	if recipe.ID != "ADVOKAAT" {
		t.Errorf("Expected ID 'ADVOKAAT', got '%s'", recipe.ID)
	}

	// Check title
	if recipe.Title != "ADVOKAAT" {
		t.Errorf("Expected title 'ADVOKAAT', got '%s'", recipe.Title)
	}

	// Check author
	if recipe.Author == nil || !strings.Contains(recipe.Author.Name, "Laurent Siklossy") {
		t.Errorf("Expected author containing 'Laurent Siklossy', got '%v'", recipe.Author)
	}

	// Check description
	if !strings.Contains(recipe.Description, "Dutch egg cognac") ||
		!strings.Contains(recipe.Description, "Advokaat is the Dutch word") {
		t.Errorf("Expected description to contain both the RZ description and introductory text, got '%s'", recipe.Description)
	}

	// Check categories - L is for Beverage (Liquid)
	if len(recipe.Categories) != 1 {
		t.Errorf("Expected 1 category, got %d", len(recipe.Categories))
	} else {
		if recipe.Categories[0].Name != "Beverage" {
			t.Errorf("Expected category 'Beverage', got '%s'", recipe.Categories[0].Name)
		}
	}

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

	if len(recipe.Ingredients) != len(expectedIngredients) {
		t.Errorf("Expected %d ingredients, got %d", len(expectedIngredients), len(recipe.Ingredients))
	} else {
		for i, ing := range expectedIngredients {
			if recipe.Ingredients[i].Name != ing.name {
				t.Errorf("Expected ingredient '%s', got '%s'", ing.name, recipe.Ingredients[i].Name)
			}
			if !strings.Contains(recipe.Ingredients[i].Unit, ing.quantityContains) {
				t.Errorf("Expected unit for '%s' to contain '%s', got '%s'",
					ing.name, ing.quantityContains, recipe.Ingredients[i].Unit)
			}
		}
	}

	// Check steps
	if len(recipe.Steps) != 5 {
		t.Errorf("Expected 5 steps, got %d", len(recipe.Steps))
	} else {
		// Just check the first and last steps
		if recipe.Steps[0].Description != "Mix sugars." {
			t.Errorf("Expected first step 'Mix sugars.', got '%s'", recipe.Steps[0].Description)
		}
		if recipe.Steps[4].Description != "Bottle and let rest for two weeks to let the mixture thicken." {
			t.Errorf("Expected last step 'Bottle and let rest for two weeks to let the mixture thicken.', got '%s'",
				recipe.Steps[4].Description)
		}
	}
}

func TestParse_EmptyInput(t *testing.T) {
	recipe, err := Parse(strings.NewReader(""))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check that we get an empty recipe
	if recipe.Title != "" {
		t.Errorf("Expected empty title, got '%s'", recipe.Title)
	}

	if len(recipe.Ingredients) != 0 {
		t.Errorf("Expected 0 ingredients, got %d", len(recipe.Ingredients))
	}

	if len(recipe.Steps) != 0 {
		t.Errorf("Expected 0 steps, got %d", len(recipe.Steps))
	}
}

func TestParse_InvalidInput(t *testing.T) {
	// Test with a nil reader
	_, err := Parse(nil)
	if err == nil {
		t.Errorf("Expected error with nil reader, got nil")
	}

	// Test with malformed input (missing closing quotes)
	input := `.RH MOD.RECIPES-SOURCE RECIPE-ID D "22 Dec 83"
.RZ "TEST RECIPE" "A simple test recipe"
.IG "1 cup sugar`
	_, err = Parse(strings.NewReader(input))
	// This should not error, but should handle the unclosed quote gracefully
	if err != nil {
		t.Errorf("Expected no error with unclosed quote, got %v", err)
	}
}
