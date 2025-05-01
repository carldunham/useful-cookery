package resolvers

import (
	"context"
	"errors"
	"fmt"

	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// Sentinel errors for recipe resolvers.
var (
	ErrMissingSkillLevel = errors.New("missing skill level")
	ErrInvalidSkillLevel = errors.New("invalid skill level")
)

// getDifficulty converts the string skillLevel from the domain model to the GraphQL enum type.
func (r *Resolver) skillLevel(_ context.Context, recipe *domainmodel.Recipe) (*gqlmodel.SkillLevel, error) {
	if recipe.SkillLevel == "" {
		return nil, ErrMissingSkillLevel
	}

	var skillLevel gqlmodel.SkillLevel
	switch recipe.SkillLevel {
	case "BEGINNER":
		skillLevel = gqlmodel.SkillLevelBeginner
	case "INTERMEDIATE":
		skillLevel = gqlmodel.SkillLevelIntermediate
	case "ADVANCED":
		skillLevel = gqlmodel.SkillLevelAdvanced
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidSkillLevel, recipe.SkillLevel)
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

// skillLevelPref converts the string skill level from the domain model to the GraphQL enum type.
func (r *Resolver) skillLevelPref(_ context.Context, prefs *domainmodel.UserPreferences) (*gqlmodel.SkillLevel, error) {
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
		return nil, fmt.Errorf("%w: %s", ErrInvalidSkillLevel, prefs.SkillLevel)
	}

	return &skillLevel, nil
}
