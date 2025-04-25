package database_test

import (
	"errors"
	"testing"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

func TestTypeAliases(t *testing.T) {
	t.Parallel()
	// Test that the type aliases are defined correctly
	var _ database.Options = dbtypes.DatabaseOptions{}
	var _ database.Type = dbtypes.DatabaseType("")

	// Test that the constants are defined correctly
	if database.DatabaseTypeDGraph != dbtypes.DatabaseTypeDGraph {
		t.Errorf("DatabaseTypeDGraph has unexpected value: %s", database.DatabaseTypeDGraph)
	}
	if database.DatabaseTypeInMemory != dbtypes.DatabaseTypeInMemory {
		t.Errorf("DatabaseTypeInMemory has unexpected value: %s", database.DatabaseTypeInMemory)
	}

	// Test that the error variables are defined correctly
	if !errors.Is(database.ErrNotFound, dbtypes.ErrNotFound) {
		t.Errorf("ErrNotFound has unexpected value: %v", database.ErrNotFound)
	}
	if !errors.Is(database.ErrInvalidID, dbtypes.ErrInvalidID) {
		t.Errorf("ErrInvalidID has unexpected value: %v", database.ErrInvalidID)
	}
	if !errors.Is(database.ErrAlreadyExists, dbtypes.ErrAlreadyExists) {
		t.Errorf("ErrAlreadyExists has unexpected value: %v", database.ErrAlreadyExists)
	}
	if !errors.Is(database.ErrEmailRequired, dbtypes.ErrEmailRequired) {
		t.Errorf("ErrEmailRequired has unexpected value: %v", database.ErrEmailRequired)
	}
}
