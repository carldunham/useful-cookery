package troff

import (
	"strings"
	"testing"
)

func TestParse_BasicRecipe(t *testing.T) {
	input := `.TL Test Recipe
.AU John Doe
.AB
This is a test recipe description.
It has multiple lines.
.AE
.SH CATEGORY
Dessert, Test
.IG "1 cup" "sugar" "200g"
.IG "2" "eggs"
.SK 1
Mix ingredients.
.SK 2
Bake for 30 minutes.`

	recipe, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check title
	if recipe.Title != "Test Recipe" {
		t.Errorf("Expected title 'Test Recipe', got '%s'", recipe.Title)
	}

	// Check author
	if recipe.Author == nil || recipe.Author.Name != "John Doe" {
		t.Errorf("Expected author 'John Doe', got '%v'", recipe.Author)
	}

	// Check description
	expectedDesc := "This is a test recipe description. It has multiple lines."
	if recipe.Description != expectedDesc {
		t.Errorf("Expected description '%s', got '%s'", expectedDesc, recipe.Description)
	}

	// Check categories
	if len(recipe.Categories) != 2 {
		t.Errorf("Expected 2 categories, got %d", len(recipe.Categories))
	} else {
		if recipe.Categories[0].Name != "Dessert" {
			t.Errorf("Expected first category 'Dessert', got '%s'", recipe.Categories[0].Name)
		}
		if recipe.Categories[1].Name != "Test" {
			t.Errorf("Expected second category 'Test', got '%s'", recipe.Categories[1].Name)
		}
	}

	// Check ingredients
	if len(recipe.Ingredients) != 2 {
		t.Errorf("Expected 2 ingredients, got %d", len(recipe.Ingredients))
	} else {
		if recipe.Ingredients[0].Name != "sugar" || recipe.Ingredients[0].Unit != "200g" {
			t.Errorf("Expected first ingredient 'sugar' with unit '200g', got '%s' with unit '%s'",
				recipe.Ingredients[0].Name, recipe.Ingredients[0].Unit)
		}
		if recipe.Ingredients[1].Name != "eggs" {
			t.Errorf("Expected second ingredient 'eggs', got '%s'", recipe.Ingredients[1].Name)
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
	input := `.TL Advokaat
.AU Laurent Siklossy
.AB
Advokaat is the Dutch word for "egg cognac".
It is highly recommended for A. I. (Alcohol Imbibing) meetings.
.AE
.SH CATEGORY
Drinks, Desserts
.IG "1 1/2 cups" "sugar" "300 g"
.IG "2 Tbsp" "vanilla sugar" "25 g"
.IG "2 cups" "milk" "500 ml"
.IG "9" "egg yolks"
.IG "1 1/2 cups" "95% grain alcohol" "350 ml"
.SK 1
Mix sugars.
.SK 2
Boil milk with half of sugars for two minutes.
.SK 3
Mix well the yolks with the other half of the sugars.
.SK 4
Using mixer add to milk. Then add alcohol slowly.
.SK 5
Bottle and let rest for two weeks to let the mixture thicken.`

	recipe, err := Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	// Check basic fields
	if recipe.Title != "Advokaat" {
		t.Errorf("Expected title 'Advokaat', got '%s'", recipe.Title)
	}

	if recipe.Author == nil || recipe.Author.Name != "Laurent Siklossy" {
		t.Errorf("Expected author 'Laurent Siklossy', got '%v'", recipe.Author)
	}

	// Check categories
	expectedCategories := []string{"Drinks", "Desserts"}
	if len(recipe.Categories) != len(expectedCategories) {
		t.Errorf("Expected %d categories, got %d", len(expectedCategories), len(recipe.Categories))
	} else {
		for i, cat := range expectedCategories {
			if recipe.Categories[i].Name != cat {
				t.Errorf("Expected category '%s', got '%s'", cat, recipe.Categories[i].Name)
			}
		}
	}

	// Check ingredients
	expectedIngredients := []struct {
		name string
		unit string
	}{
		{"sugar", "300 g"},
		{"vanilla sugar", "25 g"},
		{"milk", "500 ml"},
		{"egg yolks", ""},
		{"95% grain alcohol", "350 ml"},
	}

	if len(recipe.Ingredients) != len(expectedIngredients) {
		t.Errorf("Expected %d ingredients, got %d", len(expectedIngredients), len(recipe.Ingredients))
	} else {
		for i, ing := range expectedIngredients {
			if recipe.Ingredients[i].Name != ing.name {
				t.Errorf("Expected ingredient '%s', got '%s'", ing.name, recipe.Ingredients[i].Name)
			}
			if recipe.Ingredients[i].Unit != ing.unit {
				t.Errorf("Expected unit '%s' for ingredient '%s', got '%s'",
					ing.unit, ing.name, recipe.Ingredients[i].Unit)
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
	input := `.TL Test Recipe
.IG "1 cup sugar`
	_, err = Parse(strings.NewReader(input))
	// This should not error, but should handle the unclosed quote gracefully
	if err != nil {
		t.Errorf("Expected no error with unclosed quote, got %v", err)
	}
}
