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

// ErrNilReader is returned when a nil reader is provided.
var ErrNilReader = errors.New("nil reader provided")

// Parse converts a TROFF formatted recipe into a Recipe struct using the lexer.
func Parse(reader io.Reader) (model.Recipe, error) {
	if reader == nil {
		return model.Recipe{}, ErrNilReader
	}

	// Create a lexer for the input
	lexer := NewLexer(reader)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to tokenize content: %w", err)
	}

	recipe := model.Recipe{}

	// Process tokens to build the recipe
	processTokens(tokens, &recipe)

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
//
//nolint:gocognit,nestif,gocyclo,cyclop,funlen,maintidx // TODO: simplify.
func processTokens(tokens []Token, recipe *model.Recipe) {
	// Current state tracking
	var currentIngredientSet []model.DetailedIngredient
	var currentSteps []model.Step
	var inIntroduction bool
	var introductionText string
	var inNotes bool
	var notesText string
	var stepIndex int

	// Process each token using a traditional for loop so we can control the index
	for tokenIndex := 0; tokenIndex < len(tokens); tokenIndex++ {
		token := tokens[tokenIndex]

		switch {
		case token.Type == TokenCommand:
			cmd := token.Value

			// Handle different commands
			switch cmd {
			case "RH": // Recipe Header - must be first line with 4 arguments
				if tokenIndex+4 < len(tokens) &&
					tokens[tokenIndex+1].Type == TokenParam &&
					tokens[tokenIndex+2].Type == TokenParam &&
					tokens[tokenIndex+3].Type == TokenParam &&
					tokens[tokenIndex+4].Type == TokenParam {
					// First param is source, second is recipe ID, third is category code, fourth is date
					// We'll store the recipe ID and category
					recipe.OriginalID = tokens[tokenIndex+2].Value

					// Parse category code
					categoryCode := tokens[tokenIndex+3].Value
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

					tokenIndex += 4 // Skip the parameters we just processed
				}

			case "RZ": // Recipe Title and Description
				if tokenIndex+2 < len(tokens) &&
					tokens[tokenIndex+1].Type == TokenParam &&
					tokens[tokenIndex+2].Type == TokenParam {
					recipe.Title = tokens[tokenIndex+1].Value
					recipe.Description = tokens[tokenIndex+2].Value
					tokenIndex += 2 // Skip the parameters we just processed

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
				if tokenIndex+1 < len(tokens) && tokens[tokenIndex+1].Type == TokenParam {
					// Store yield information if needed
					// For now we'll just skip it
					tokenIndex++
				}

			case "IG": // Ingredient
				if tokenIndex+2 < len(tokens) &&
					tokens[tokenIndex+1].Type == TokenParam &&
					tokens[tokenIndex+2].Type == TokenParam {
					quantity := tokens[tokenIndex+1].Value
					name := tokens[tokenIndex+2].Value
					tokenIndex += 2 // Skip the parameters we just processed

					// Check for optional metric quantity
					var metricQty string
					if tokenIndex+1 < len(tokens) && tokens[tokenIndex+1].Type == TokenParam {
						metricQty = tokens[tokenIndex+1].Value
						tokenIndex++ // Skip the metric parameter
					}

					// Create ingredient with proper types
					// If we have metric quantity, append it to the unit field
					unitValue := quantity
					if metricQty != "" {
						unitValue += " (metric: " + metricQty + ")"
					}

					// Convert TROFF codes in the name and unit
					processedName := ProcessText(name)
					processedUnit := ProcessText(unitValue)

					ingredient := model.DetailedIngredient{
						Name: processedName,
						Unit: processedUnit,
					}

					currentIngredientSet = append(currentIngredientSet, ingredient)
				}

			case "PH": // Procedure Header
				// Nothing specific to do here, just marks the start of procedure steps

			case "SK": // Step
				if tokenIndex+1 < len(tokens) && tokens[tokenIndex+1].Type == TokenParam {
					// The parameter is the step number
					stepNumber := tokens[tokenIndex+1].Value
					tokenIndex++ // Skip the step number parameter

					// Collect the step description from text following this command
					var stepDescription string
					for nestedIndex := tokenIndex + 1; nestedIndex < len(tokens); nestedIndex++ {
						if tokens[nestedIndex].Type == TokenCommand {
							break
						}
						if tokens[nestedIndex].Type == TokenParam {
							if stepDescription != "" {
								stepDescription += " "
							}
							stepDescription += tokens[nestedIndex].Value
						}
					}

					// Try to parse step number
					var orderIndex int
					if num, err := strconv.Atoi(stepNumber); err == nil {
						orderIndex = num - 1 // Convert to 0-based index
					} else {
						orderIndex = stepIndex
					}

					// Convert TROFF codes in the step description
					processedDescription := ProcessText(stepDescription)

					step := model.Step{
						OrderIndex:  orderIndex,
						Description: processedDescription,
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
					recipe.Notes = notesText
				}

				// Process author information from text following WR
				var authorInfo string
				for nestedIndex := tokenIndex + 1; nestedIndex < len(tokens); nestedIndex++ {
					if tokens[nestedIndex].Type == TokenCommand {
						break
					}
					if tokens[nestedIndex].Type == TokenParam {
						if authorInfo != "" {
							authorInfo += " "
						}
						authorInfo += tokens[nestedIndex].Value
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
				// Convert TROFF codes in the introduction text
				introductionText += ConvertCodes(token.Value)
			} else if inNotes {
				if notesText != "" {
					notesText += " "
				}
				// Convert TROFF codes in the notes text
				notesText += ConvertCodes(token.Value)
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
}
