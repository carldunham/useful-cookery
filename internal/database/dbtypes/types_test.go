package dbtypes_test

import (
	"testing"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

func TestDatabaseTypes(t *testing.T) {
	t.Parallel()
	// Test that database type constants are defined
	if dbtypes.DatabaseTypeDGraph != "dgraph" {
		t.Errorf("DatabaseTypeDGraph has unexpected value: %s", dbtypes.DatabaseTypeDGraph)
	}
	if dbtypes.DatabaseTypeInMemory != "memory" {
		t.Errorf("DatabaseTypeInMemory has unexpected value: %s", dbtypes.DatabaseTypeInMemory)
	}

	// Test that database type constants are different
	if dbtypes.DatabaseTypeDGraph == dbtypes.DatabaseTypeInMemory {
		t.Error("DatabaseTypeDGraph and DatabaseTypeInMemory should be different")
	}
}
