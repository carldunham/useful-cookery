package troff

import (
	"strings"
	"testing"
)

func TestLexer_TokenizeCommand(t *testing.T) {
	input := ".TL Test Recipe"
	lexer := NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	if len(tokens) != 3 { // TL command + "Test" param + "Recipe" param + EOF
		t.Errorf("Expected 3 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TokenCommand || tokens[0].Value != "TL" {
		t.Errorf("Expected first token to be Command(TL), got %v", tokens[0])
	}

	if tokens[1].Type != TokenParam || tokens[1].Value != "Test" {
		t.Errorf("Expected second token to be Param(Test), got %v", tokens[1])
	}

	if tokens[2].Type != TokenParam || tokens[2].Value != "Recipe" {
		t.Errorf("Expected third token to be Param(Recipe), got %v", tokens[2])
	}
}

func TestLexer_TokenizeMultipleLines(t *testing.T) {
	input := `.TL Test Recipe
.AU John Doe
.AB
This is a description
of a recipe.
.AE`

	lexer := NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	// Check for the correct number of tokens
	expectedCount := 9 // TL + Test + Recipe + AU + John + Doe + AB + description param + AE (EOF is dropped)
	if len(tokens) != expectedCount {
		t.Errorf("Expected %d tokens, got %d", expectedCount, len(tokens))
	}

	// Check for specific tokens
	if tokens[0].Type != TokenCommand || tokens[0].Value != "TL" {
		t.Errorf("Expected first token to be Command(TL), got %v", tokens[0])
	}

	if tokens[3].Type != TokenCommand || tokens[3].Value != "AU" {
		t.Errorf("Expected fourth token to be Command(AU), got %v", tokens[3])
	}

	if tokens[6].Type != TokenCommand || tokens[6].Value != "AB" {
		t.Errorf("Expected seventh token to be Command(AB), got %v", tokens[6])
	}

	// Check description tokens
	if tokens[7].Type != TokenParam || tokens[7].Value != "This is a description\nof a recipe." {
		t.Errorf("Expected description token to be multi-line Param(This is a description\nof a recipe.), got %v", tokens[7])
	}
}

func TestLexer_TokenizeQuotedParams(t *testing.T) {
	input := `.IG "1 1/2 cups" "sugar" "300 g"`
	lexer := NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	if len(tokens) != 4 { // IG + 3 params + EOF
		t.Errorf("Expected 4 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TokenCommand || tokens[0].Value != "IG" {
		t.Errorf("Expected first token to be Command(IG), got %v", tokens[0])
	}

	if tokens[1].Type != TokenParam || tokens[1].Value != "1 1/2 cups" {
		t.Errorf("Expected second token to be Param(1 1/2 cups), got %v", tokens[1])
	}

	if tokens[2].Type != TokenParam || tokens[2].Value != "sugar" {
		t.Errorf("Expected third token to be Param(sugar), got %v", tokens[2])
	}

	if tokens[3].Type != TokenParam || tokens[3].Value != "300 g" {
		t.Errorf("Expected fourth token to be Param(300 g), got %v", tokens[3])
	}
}

func TestLexer_TokenizeComplexInput(t *testing.T) {
	input := `.TL Advokaat
.SH CATEGORY
Drinks, Desserts
.AU Laurent Siklossy
.AB
Advokaat is the Dutch word for "egg cognac".
It is highly recommended for A. I. (Alcohol Imbibing) meetings.
.AE
.IG "1 1/2 cups" "sugar" "300 g"
.IG "2 Tbsp" "vanilla sugar" "25 g"
.SK 1
Mix sugars.
.SK 2
Boil milk with half of sugars for two minutes.`

	lexer := NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	// Just check that we have a reasonable number of tokens
	if len(tokens) < 20 {
		t.Errorf("Expected at least 20 tokens, got %d", len(tokens))
	}

	// Check a few key tokens
	foundTL := false
	foundSH := false
	foundAU := false
	foundAB := false
	foundIG := false
	foundSK := false

	for _, token := range tokens {
		if token.Type == TokenCommand {
			switch token.Value {
			case "TL":
				foundTL = true
			case "SH":
				foundSH = true
			case "AU":
				foundAU = true
			case "AB":
				foundAB = true
			case "IG":
				foundIG = true
			case "SK":
				foundSK = true
			}
		}
	}

	if !foundTL || !foundSH || !foundAU || !foundAB || !foundIG || !foundSK {
		t.Errorf("Not all expected commands were found")
	}
}

func TestLexer_TokenizeMultilineParams(t *testing.T) {
	input := `.TL Test Recipe
.AB
This is a multi-line description
that spans multiple lines
and should be treated as separate tokens
but all belonging to the AB command.
.AE`

	lexer := NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	if err != nil {
		t.Fatalf("Tokenize() error = %v", err)
	}

	// Find the AB command and check that all following parameters have the AB command
	abIndex := -1
	for i, token := range tokens {
		if token.Type == TokenCommand && token.Value == "AB" {
			abIndex = i
			break
		}
	}

	if abIndex == -1 {
		t.Fatalf("AB command not found in tokens")
	}

	// Check that all parameters between AB and AE have the AB command
	aeIndex := -1
	for i := abIndex + 1; i < len(tokens); i++ {
		if tokens[i].Type == TokenCommand && tokens[i].Value == "AE" {
			aeIndex = i
			break
		}

		if tokens[i].Type == TokenParam {
			if tokens[i].Command != "AB" {
				t.Errorf("Expected parameter token %d to have command AB, got %s", i, tokens[i].Command)
			}
		}
	}

	if aeIndex == -1 {
		t.Fatalf("AE command not found in tokens")
	}

	// Check the content of each description line
	paramIndex := abIndex + 1
	if paramIndex < aeIndex && tokens[paramIndex].Value[:4] != "This" {
		t.Errorf("Expected first word of description to be 'This', got '%s'", tokens[paramIndex].Value[:4])
	}
}
