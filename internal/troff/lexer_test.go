package troff_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carldunham/useful-cookery/internal/troff"
)

func TestLexer_TokenizeCommand(t *testing.T) {
	t.Parallel()
	input := ".TL Test Recipe"
	lexer := troff.NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	require.NoError(t, err, "Tokenize() should not return an error")

	assert.Len(t, tokens, 3, "Should have 3 tokens") // TL command + "Test" param + "Recipe" param + EOF

	assert.Equal(t, troff.TokenCommand, tokens[0].Type, "First token should be a Command")
	assert.Equal(t, "TL", tokens[0].Value, "First token should be TL command")

	assert.Equal(t, troff.TokenParam, tokens[1].Type, "Second token should be a Param")
	assert.Equal(t, "Test", tokens[1].Value, "Second token should be 'Test' param")

	assert.Equal(t, troff.TokenParam, tokens[2].Type, "Third token should be a Param")
	assert.Equal(t, "Recipe", tokens[2].Value, "Third token should be 'Recipe' param")
}

func TestLexer_TokenizeMultipleLines(t *testing.T) {
	t.Parallel()
	input := `.TL Test Recipe
.AU John Doe
.AB
This is a description
of a recipe.
.AE`

	lexer := troff.NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	require.NoError(t, err, "Tokenize() should not return an error")

	// Check for the correct number of tokens
	expectedCount := 9 // TL + Test + Recipe + AU + John + Doe + AB + description param + AE (EOF is dropped)
	assert.Len(t, tokens, expectedCount, "Should have correct number of tokens")

	// Check for specific tokens
	assert.Equal(t, troff.TokenCommand, tokens[0].Type, "First token should be a Command")
	assert.Equal(t, "TL", tokens[0].Value, "First token should be TL command")

	assert.Equal(t, troff.TokenCommand, tokens[3].Type, "Fourth token should be a Command")
	assert.Equal(t, "AU", tokens[3].Value, "Fourth token should be AU command")

	assert.Equal(t, troff.TokenCommand, tokens[6].Type, "Seventh token should be a Command")
	assert.Equal(t, "AB", tokens[6].Value, "Seventh token should be AB command")

	// Check description tokens
	assert.Equal(t, troff.TokenParam, tokens[7].Type, "Description token should be a Param")
	assert.Equal(t, "This is a description\nof a recipe.", tokens[7].Value,
		"Description token should contain multi-line text")
}

func TestLexer_TokenizeQuotedParams(t *testing.T) {
	t.Parallel()
	input := `.IG "1 1/2 cups" "sugar" "300 g"`
	lexer := troff.NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	require.NoError(t, err, "Tokenize() should not return an error")

	assert.Len(t, tokens, 4, "Should have 4 tokens") // IG + 3 params + EOF

	assert.Equal(t, troff.TokenCommand, tokens[0].Type, "First token should be a Command")
	assert.Equal(t, "IG", tokens[0].Value, "First token should be IG command")

	assert.Equal(t, troff.TokenParam, tokens[1].Type, "Second token should be a Param")
	assert.Equal(t, "1 1/2 cups", tokens[1].Value, "Second token should be '1 1/2 cups' param")

	assert.Equal(t, troff.TokenParam, tokens[2].Type, "Third token should be a Param")
	assert.Equal(t, "sugar", tokens[2].Value, "Third token should be 'sugar' param")

	assert.Equal(t, troff.TokenParam, tokens[3].Type, "Fourth token should be a Param")
	assert.Equal(t, "300 g", tokens[3].Value, "Fourth token should be '300 g' param")
}

func TestLexer_TokenizeComplexInput(t *testing.T) {
	t.Parallel()
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

	lexer := troff.NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	require.NoError(t, err, "Tokenize() should not return an error")

	// Just check that we have a reasonable number of tokens
	assert.GreaterOrEqual(t, len(tokens), 20, "Should have at least 20 tokens")

	// Check a few key tokens
	foundTL := false
	foundSH := false
	foundAU := false
	foundAB := false
	foundIG := false
	foundSK := false

	for _, token := range tokens {
		if token.Type == troff.TokenCommand {
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

	assert.True(t, foundTL, "Should find TL command")
	assert.True(t, foundSH, "Should find SH command")
	assert.True(t, foundAU, "Should find AU command")
	assert.True(t, foundAB, "Should find AB command")
	assert.True(t, foundIG, "Should find IG command")
	assert.True(t, foundSK, "Should find SK command")
}

func TestLexer_TokenizeMultilineParams(t *testing.T) {
	t.Parallel()
	input := `.TL Test Recipe
.AB
This is a multi-line description
that spans multiple lines
and should be treated as separate tokens
but all belonging to the AB command.
.AE`

	lexer := troff.NewLexer(strings.NewReader(input))
	tokens, err := lexer.Tokenize()
	require.NoError(t, err, "Tokenize() should not return an error")

	// Find the AB command and check that all following parameters have the AB command
	abIndex := -1
	for tokenIndex, token := range tokens {
		if token.Type == troff.TokenCommand && token.Value == "AB" {
			abIndex = tokenIndex
			break
		}
	}

	require.NotEqual(t, -1, abIndex, "AB command should be found in tokens")

	// Check that all parameters between AB and AE have the AB command
	aeIndex := -1
	for tokenIndex := abIndex + 1; tokenIndex < len(tokens); tokenIndex++ {
		if tokens[tokenIndex].Type == troff.TokenCommand && tokens[tokenIndex].Value == "AE" {
			aeIndex = tokenIndex
			break
		}

		if tokens[tokenIndex].Type == troff.TokenParam {
			assert.Equal(t, "AB", tokens[tokenIndex].Command, "Parameter token should have AB command")
		}
	}

	require.NotEqual(t, -1, aeIndex, "AE command should be found in tokens")

	// Check the content of each description line
	paramIndex := abIndex + 1
	if paramIndex < aeIndex {
		assert.GreaterOrEqual(t, len(tokens[paramIndex].Value), 4, "Description should have at least 4 characters")
		assert.Equal(t, "This", tokens[paramIndex].Value[:4], "First word of description should be 'This'")
	}
}
