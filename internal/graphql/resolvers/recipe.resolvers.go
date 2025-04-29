package resolvers

import (
	"context"
	"errors"

	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// Sentinel errors for recipe resolvers.
var (
	ErrInvalidSkillLevel = errors.New("invalid skill level")
)

// getDifficulty converts the string difficulty from the domain model to the GraphQL enum type.
func (r *Resolver) difficulty(_ context.Context, recipe *domainmodel.Recipe) (*gqlmodel.SkillLevel, error) {
	if recipe.Difficulty == "" {
		return nil, ErrInvalidSkillLevel
	}

	var skillLevel gqlmodel.SkillLevel
	switch recipe.Difficulty {
	case "BEGINNER":
		skillLevel = gqlmodel.SkillLevelBeginner
	case "INTERMEDIATE":
		skillLevel = gqlmodel.SkillLevelIntermediate
	case "ADVANCED":
		skillLevel = gqlmodel.SkillLevelAdvanced
	default:
		return nil, ErrInvalidSkillLevel
	}

	return &skillLevel, nil
}

// difficultyText is the resolver for the difficultyText field.
func (r *recipeResolver) difficultyText(_ context.Context, obj *domainmodel.Recipe) (*string, error) {
	// If DifficultyText is empty but Difficulty is set, we can return the Difficulty as a fallback
	if obj.DifficultyText == "" && obj.Difficulty != "" {
		return &obj.Difficulty, nil
	}

	if obj.DifficultyText == "" {
		return nil, nil
	}

	return &obj.DifficultyText, nil
}

// skillLevel converts the string skill level from the domain model to the GraphQL enum type.
func (r *Resolver) skillLevel(_ context.Context, prefs *domainmodel.UserPreferences) (*gqlmodel.SkillLevel, error) {
	if prefs.SkillLevel == "" {
		return nil, ErrInvalidSkillLevel
	}

	var skillLevel gqlmodel.SkillLevel
	switch prefs.SkillLevel {
	case "BEGINNER":
		skillLevel = gqlmodel.SkillLevelBeginner
	case "INTERMEDIATE":
		skillLevel = gqlmodel.SkillLevelIntermediate
	case "ADVANCED":
		skillLevel = gqlmodel.SkillLevelAdvanced
	default:
		return nil, ErrInvalidSkillLevel
	}

	return &skillLevel, nil
}
