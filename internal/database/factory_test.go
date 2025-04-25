package database_test

import (
	"testing"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

func TestNewDGraphDatabase(t *testing.T) {
	t.Parallel()
	// Test creating a DGraph database
	dgraphDB, err := database.NewDGraphDatabase(dbtypes.DatabaseOptions{
		ConnectionString: "dgraph://localhost:9080",
	})
	if err != nil {
		t.Fatalf("Failed to create DGraph database: %v", err)
	}
	if dgraphDB == nil {
		t.Fatal("DGraph database should not be nil")
	}
}

func TestNewInMemoryDatabase(t *testing.T) {
	t.Parallel()
	// Test creating an in-memory database
	memoryDB, err := database.NewInMemoryDatabase(dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	if memoryDB == nil {
		t.Fatal("In-memory database should not be nil")
	}
}
