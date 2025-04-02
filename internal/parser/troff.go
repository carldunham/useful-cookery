package parser

import (
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/carldunham/useful-cookery/internal/model"
)

// Common TROFF formatting patterns.
var (
	titleRegex       = regexp.MustCompile(`\.TL\s+(.+)`)
	categoryRegex    = regexp.MustCompile(`\.SH\s+CATEGORY\s*\n(.+)`)
	ingredientRegex  = regexp.MustCompile(`\.IG\s+(.+)`)
	stepRegex        = regexp.MustCompile(`\.SK\s*\n(.+)`)
	descriptionRegex = regexp.MustCompile(`\.AB\s*\n([\s\S]+?)\n\.AE`)
	authorRegex      = regexp.MustCompile(`\.AU\s+(.+)`)
)

// TroffParser implements the parsing of TROFF formatted recipe files.
type TroffParser struct {
	// Configuration options could be added here
}

// NewTroffParser creates a new TROFF parser.
func NewTroffParser() *TroffParser {
	return &TroffParser{}
}

// Parse converts a TROFF formatted recipe into a Recipe struct.
func (p *TroffParser) Parse(r io.Reader) (model.Recipe, error) {
	if r == nil {
		return model.Recipe{}, errors.New("nil reader provided")
	}

	content, err := io.ReadAll(r)
	if err != nil {
		return model.Recipe{}, fmt.Errorf("failed to read content: %w", err)
	}

	strContent := string(content)
	recipe := model.Recipe{}

	// Extract title
	titleMatch := titleRegex.FindStringSubmatch(strContent)
	if len(titleMatch) > 1 {
		recipe.Title = strings.TrimSpace(titleMatch[1])
	}

	// Extract categories
	categories := []model.Category{}
	categoryMatch := categoryRegex.FindStringSubmatch(strContent)
	if len(categoryMatch) > 1 {
		categoryList := strings.Split(categoryMatch[1], ",")
		for _, cat := range categoryList {
			cat = strings.TrimSpace(cat)
			if cat != "" {
				categories = append(categories, model.Category{Name: cat})
			}
		}
	}
	recipe.Categories = categories

	// Extract author
	authorMatch := authorRegex.FindStringSubmatch(strContent)
	if len(authorMatch) > 1 {
		recipe.Author = &model.User{Name: strings.TrimSpace(authorMatch[1])}
	}

	// Extract description
	descMatch := descriptionRegex.FindStringSubmatch(strContent)
	if len(descMatch) > 1 {
		recipe.Description = strings.TrimSpace(descMatch[1])
	}

	// Extract ingredients
	ingredientMatches := ingredientRegex.FindAllStringSubmatch(strContent, -1)
	ingredients := []model.DetailedIngredient{}
	for _, match := range ingredientMatches {
		if len(match) > 1 {
			ingredient := parseIngredient(match[1])
			ingredients = append(ingredients, ingredient)
		}
	}
	recipe.Ingredients = ingredients

	// Extract steps
	stepMatches := stepRegex.FindAllStringSubmatch(strContent, -1)
	steps := []model.Step{}
	for i, match := range stepMatches {
		if len(match) > 1 {
			step := model.Step{
				OrderIndex:  i,
				Description: strings.TrimSpace(match[1]),
			}
			steps = append(steps, step)
		}
	}
	recipe.Steps = steps

	return recipe, nil
}

// ParseFile parses a TROFF file and returns a Recipe.
func (p *TroffParser) ParseFile(filename string) (model.Recipe, error) {
	// Implementation would open the file and call Parse
	return model.Recipe{}, errors.New("not implemented")
}

// parseIngredient parses an ingredient line into a structured ingredient.
func parseIngredient(line string) model.DetailedIngredient {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)

	ingredient := model.DetailedIngredient{
		Name: line, // Default to full line if we can't parse
	}

	if len(parts) >= 3 {
		// Try to parse quantity and unit
		// This is a simplified approach - real implementation would be more robust
		quantity := parts[0]
		unit := parts[1]
		name := strings.Join(parts[2:], " ")

		// Check if the first part could be a quantity
		if regexQuantity.MatchString(quantity) {
			ingredient.Name = name
			ingredient.Unit = unit
			// Actual quantity parsing would go here
		}
	}

	return ingredient
}

// Regex for matching common quantity formats.
var regexQuantity = regexp.MustCompile(`^(\d+/?\.?\d*)$`)

// ParseToJSON converts a TROFF file to a JSON string.
func (p *TroffParser) ParseToJSON(r io.Reader) (string, error) {
	// Implementation would parse the TROFF and convert to JSON
	return "", errors.New("not implemented")
}
