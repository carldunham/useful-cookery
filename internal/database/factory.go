package database

import (
	"errors"
	"fmt"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/database/strategies"
)

// ErrUnsupportedDatabaseType is returned when an unsupported database type is specified.
var ErrUnsupportedDatabaseType = errors.New("unsupported database type")

// NewDGraphDatabase creates a new DGraph database instance.
func NewDGraphDatabase(options dbtypes.DatabaseOptions) (*strategies.DGraphDatabase, error) {
	db, err := strategies.NewDGraphDatabase(options)
	if err != nil {
		return nil, fmt.Errorf("creating DGraph database: %w", err)
	}
	return db, nil
}

// NewInMemoryDatabase creates a new in-memory database instance.
func NewInMemoryDatabase(options dbtypes.DatabaseOptions) (*strategies.InMemoryDatabase, error) {
	db, err := strategies.NewInMemoryDatabase(options)
	if err != nil {
		return nil, fmt.Errorf("creating in-memory database: %w", err)
	}
	return db, nil
}
