package resolvers

import (
	"context"

	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// getDifficulty converts the string difficulty from the domain model to the GraphQL enum type.
func (r *Resolver) difficulty(ctx context.Context, recipe *domainmodel.Recipe) (*gqlmodel.SkillLevel, error) {
	if recipe.Difficulty == "" {
		return nil, nil
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
		return nil, nil
	}

	return &skillLevel, nil
}

// skillLevel converts the string skill level from the domain model to the GraphQL enum type.
func (r *Resolver) skillLevel(ctx context.Context, prefs *domainmodel.UserPreferences) (*gqlmodel.SkillLevel, error) {
	if prefs.SkillLevel == "" {
		return nil, nil
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
		return nil, nil
	}

	return &skillLevel, nil
}
