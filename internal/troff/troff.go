package troff

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/carldunham/useful-cookery/internal/model"
)

// Parser is the interface for parsing recipe files.
type Parser interface {
	Parse(r io.Reader) (model.Recipe, error)
	ParseFile(filename string) (model.Recipe, error)
}

// TROFFParser implements the Parser interface for TROFF formatted recipe files.
type TROFFParser struct{}

// NewTROFFParser creates a new TROFF parser.
func NewTROFFParser() *TROFFParser {
	return &TROFFParser{}
}

// Parse converts a TROFF formatted recipe into a Recipe struct.
func (p *TROFFParser) Parse(r io.Reader) (model.Recipe, error) {
	return Parse(r)
}

// ParseFile parses a TROFF file and returns a Recipe.
func (p *TROFFParser) ParseFile(filename string) (model.Recipe, error) {
	return ParseFile(filename)
}

// Parse converts a TROFF formatted recipe into a Recipe struct using the lexer.
func Parse(r io.Reader) (model.Recipe, error) {
	if r == nil {
		return model.Recipe{}, errors.New("nil reader provided")
	}

	// Create a lexer for the input
	lexer := NewLexer(r)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to tokenize content: %w", err)
	}

	recipe := model.Recipe{}

	// Process tokens to build the recipe
	err = processTokens(tokens, &recipe)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to process tokens: %w", err)
	}

	return recipe, nil
}

// ParseFile parses a TROFF file and returns a Recipe.
func ParseFile(filename string) (model.Recipe, error) {
	file, err := os.Open(filename)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	return Parse(file)
}

// processTokens processes the tokens and builds a Recipe struct.
func processTokens(tokens []Token, recipe *model.Recipe) error {
	// Current state tracking
	var currentIngredientSet []model.DetailedIngredient
	var currentSteps []model.Step
	var inIntroduction bool
	var introductionText string
	var inNotes bool
	var notesText string
	var stepIndex int

	// Process each token
	for i, token := range tokens {
		switch {
		case token.Type == TokenCommand:
			cmd := token.Value

			// Handle different commands
			switch cmd {
			case "RH": // Recipe Header - must be first line with 4 arguments
				if i+3 < len(tokens) &&
					tokens[i+1].Type == TokenParam &&
					tokens[i+2].Type == TokenParam &&
					tokens[i+3].Type == TokenParam &&
					tokens[i+4].Type == TokenParam {
					// First param is source, second is recipe ID, third is category code, fourth is date
					// We'll store the recipe ID and category
					recipe.ID = tokens[i+2].Value

					// Parse category code
					categoryCode := tokens[i+3].Value
					for _, code := range strings.Split(categoryCode, "") {
						switch code {
						case "M":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Main dish"})
						case "A":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Appetizer or snack"})
						case "B":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Bread/pasta"})
						case "L":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Beverage"})
						case "C":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Cookie or cake"})
						case "S":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Sauce"})
						case "SL":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Salad"})
						case "SP":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Soup"})
						case "D":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Dessert"})
						case "V":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Vegetable dish"})
						case "O":
							recipe.Categories = append(recipe.Categories, model.Category{Name: "Other"})
						}
					}

					i += 4 // Skip the parameters we just processed
				}

			case "RZ": // Recipe Title and Description
				if i+2 < len(tokens) && tokens[i+1].Type == TokenParam && tokens[i+2].Type == TokenParam {
					recipe.Title = tokens[i+1].Value
					recipe.Description = tokens[i+2].Value
					i += 2 // Skip the parameters we just processed

					// After RZ, introductory comments begin
					inIntroduction = true
					introductionText = ""
				}

			case "IH": // Ingredients Header
				inIntroduction = false // End of introduction
				if len(introductionText) > 0 {
					if recipe.Description != "" {
						recipe.Description += "\n\n"
					}
					recipe.Description += introductionText
				}

				// IH can have yield information
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					// Store yield information if needed
					// For now we'll just skip it
					i++
				}

			case "IG": // Ingredient
				if i+2 < len(tokens) && tokens[i+1].Type == TokenParam && tokens[i+2].Type == TokenParam {
					quantity := tokens[i+1].Value
					name := tokens[i+2].Value
					i += 2 // Skip the parameters we just processed

					// Check for optional metric quantity
					var metricQty string
					if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
						metricQty = tokens[i+1].Value
						i++ // Skip the metric parameter
					}

					// Create ingredient with proper types
					// If we have metric quantity, append it to the unit field
					unitValue := quantity
					if metricQty != "" {
						unitValue += " (metric: " + metricQty + ")"
					}

					ingredient := model.DetailedIngredient{
						Name: name,
						Unit: unitValue,
					}

					currentIngredientSet = append(currentIngredientSet, ingredient)
				}

			case "PH": // Procedure Header
				// Nothing specific to do here, just marks the start of procedure steps

			case "SK": // Step
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					// The parameter is the step number
					stepNumber := tokens[i+1].Value
					i++ // Skip the step number parameter

					// Collect the step description from text following this command
					var stepDescription string
					for j := i + 1; j < len(tokens); j++ {
						if tokens[j].Type == TokenCommand {
							break
						}
						if tokens[j].Type == TokenParam {
							if stepDescription != "" {
								stepDescription += " "
							}
							stepDescription += tokens[j].Value
						}
					}

					// Try to parse step number
					var orderIndex int
					if num, err := strconv.Atoi(stepNumber); err == nil {
						orderIndex = num - 1 // Convert to 0-based index
					} else {
						orderIndex = stepIndex
					}

					step := model.Step{
						OrderIndex:  orderIndex,
						Description: stepDescription,
					}
					currentSteps = append(currentSteps, step)
					stepIndex++
				}

			case "NX": // Notes Header
				inNotes = true
				notesText = ""

			case "WR": // Wrapup
				inNotes = false
				if len(notesText) > 0 {
					if recipe.Description != "" {
						recipe.Description += "\n\n"
					}
					recipe.Description += "NOTES: " + notesText
				}

				// Process author information from text following WR
				var authorInfo string
				for j := i + 1; j < len(tokens); j++ {
					if tokens[j].Type == TokenCommand {
						break
					}
					if tokens[j].Type == TokenParam {
						if authorInfo != "" {
							authorInfo += " "
						}
						authorInfo += tokens[j].Value
					}
				}

				if authorInfo != "" {
					recipe.Author = &model.User{Name: authorInfo}
				}
			}

		case token.Type == TokenParam:
			// Handle parameters based on context
			if inIntroduction {
				if introductionText != "" {
					introductionText += " "
				}
				introductionText += token.Value
			} else if inNotes {
				if notesText != "" {
					notesText += " "
				}
				notesText += token.Value
			}
		}
	}

	// Add any remaining ingredients and steps to the recipe
	if len(currentIngredientSet) > 0 {
		recipe.Ingredients = currentIngredientSet
	}

	if len(currentSteps) > 0 {
		recipe.Steps = currentSteps
	}

	return nil
}

// Helper function to split a string and trim each part
func splitAndTrim(s, sep string) []string {
	parts := make([]string, 0)
	for _, part := range strings.Split(s, sep) {
		parts = append(parts, strings.TrimSpace(part))
	}
	return parts
}
