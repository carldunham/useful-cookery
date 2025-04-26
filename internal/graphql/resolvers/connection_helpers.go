package resolvers

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gqlmodel "github.com/carldunham/useful-cookery/internal/graphql/model"
	domainmodel "github.com/carldunham/useful-cookery/internal/model"
)

// EncodeCursor encodes an offset into a cursor string.
func EncodeCursor(offset int) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("offset:%d", offset)))
}

// ErrInvalidCursorFormat is returned when a cursor has an invalid format.
var ErrInvalidCursorFormat = errors.New("invalid cursor format")

// DecodeCursor decodes a cursor string into an offset.
func DecodeCursor(cursor string) (int, error) {
	bytes, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return 0, fmt.Errorf("decoding cursor: %w", err)
	}

	str := string(bytes)
	if !strings.HasPrefix(str, "offset:") {
		return 0, ErrInvalidCursorFormat
	}

	offsetStr := strings.TrimPrefix(str, "offset:")
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		return 0, fmt.Errorf("parsing offset: %w", err)
	}

	return offset, nil
}

// CreateRecipeConnection creates a RecipeConnection from a slice of recipes.
func CreateRecipeConnection(recipes []*domainmodel.Recipe, totalCount, _, offset int) *gqlmodel.RecipeConnection {
	edges := make([]*gqlmodel.RecipeEdge, len(recipes))
	for i, recipe := range recipes {
		edges[i] = &gqlmodel.RecipeEdge{
			Node:   recipe,
			Cursor: EncodeCursor(offset + i),
		}
	}

	var startCursor, endCursor *string
	if len(edges) > 0 {
		sc := edges[0].Cursor
		ec := edges[len(edges)-1].Cursor
		startCursor = &sc
		endCursor = &ec
	}

	return &gqlmodel.RecipeConnection{
		Edges: edges,
		PageInfo: &gqlmodel.PageInfo{
			HasNextPage:     offset+len(recipes) < totalCount,
			HasPreviousPage: offset > 0,
			StartCursor:     startCursor,
			EndCursor:       endCursor,
		},
		TotalCount: totalCount,
	}
}
