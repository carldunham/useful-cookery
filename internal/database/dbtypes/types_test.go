package dbtypes_test

import (
	"testing"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

func TestDatabaseTypes(t *testing.T) {
	t.Parallel()
	// Test that database type constants are defined
	if dbtypes.DatabaseTypePostgres != "postgres" {
		t.Errorf("DatabaseTypePostgres has unexpected value: %s", dbtypes.DatabaseTypePostgres)
	}
	if dbtypes.DatabaseTypeInMemory != "memory" {
		t.Errorf("DatabaseTypeInMemory has unexpected value: %s", dbtypes.DatabaseTypeInMemory)
	}

	// Test that database type constants are different
	if dbtypes.DatabaseTypePostgres == dbtypes.DatabaseTypeInMemory {
		t.Error("DatabaseTypePostgres and DatabaseTypeInMemory should be different")
	}
}
