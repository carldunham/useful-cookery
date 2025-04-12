package troff

import (
	"errors"
	"fmt"
	"io"
	"os"
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
	var inDescription bool
	var descriptionText string
	var stepIndex int

	// Process each token
	for i := 0; i < len(tokens); i++ {
		token := tokens[i]

		switch {
		case token.Type == TokenCommand:
			cmd := token.Value

			// Handle different commands
			switch cmd {
			case "TL": // Title
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					recipe.Title = tokens[i+1].Value
					i++ // Skip the parameter we just processed
				}

			case "SH": // Section Header
				// Check if we have a parameter for the section
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					sectionName := tokens[i+1].Value
					i++ // Skip the parameter we just processed

					// Handle special sections
					if sectionName == "CATEGORY" {
						// Process category parameters
						categories := []model.Category{}
						for j := i + 1; j < len(tokens); j++ {
							if tokens[j].Type == TokenParam && tokens[j].Command == token.Value {
								categoryNames := tokens[j].Value
								for _, name := range splitAndTrim(categoryNames, ",") {
									if name != "" {
										categories = append(categories, model.Category{Name: name})
									}
								}
								i = j // Skip the parameters we just processed
								break
							} else if tokens[j].Type == TokenCommand {
								break
							}
						}
						recipe.Categories = categories
					}
				}

			case "AU": // Author
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					recipe.Author = &model.User{Name: tokens[i+1].Value}
					i++ // Skip the parameter we just processed
				}

			case "AB": // Abstract (Description) start
				inDescription = true
				descriptionText = ""

			case "AE": // Abstract (Description) end
				inDescription = false
				recipe.Description = descriptionText

			case "IG": // Ingredient
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					_ = tokens[i+1].Value // quantity (stored for future use if needed)
					i++                   // Skip the quantity parameter

					// Get the ingredient name
					if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
						name := tokens[i+1].Value
						i++ // Skip the name parameter

						// Check for optional metric quantity
						var unit string
						if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
							unit = tokens[i+1].Value
							i++ // Skip the unit parameter
						}

						ingredient := model.DetailedIngredient{
							Name: name,
							Unit: unit,
						}

						// Try to parse the quantity as a float
						currentIngredientSet = append(currentIngredientSet, ingredient)
					}
				}

			case "SK": // Step
				if i+1 < len(tokens) && tokens[i+1].Type == TokenParam {
					// The parameter is the step number
					i++ // Skip the step number parameter

					// Collect the step description from subsequent parameters
					var stepDescription string
					for j := i + 1; j < len(tokens); j++ {
						if tokens[j].Type == TokenParam && tokens[j].Command == token.Value {
							if stepDescription != "" {
								stepDescription += " "
							}
							stepDescription += tokens[j].Value
							i = j // Skip the parameters we just processed
						} else if tokens[j].Type == TokenCommand {
							break
						}
					}

					step := model.Step{
						OrderIndex:  stepIndex,
						Description: stepDescription,
					}
					currentSteps = append(currentSteps, step)
					stepIndex++
				}
			}

		case token.Type == TokenParam:
			// Handle parameters based on context
			if inDescription {
				if descriptionText != "" {
					descriptionText += " "
				}
				descriptionText += token.Value
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
