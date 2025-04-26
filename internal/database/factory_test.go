package database_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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

func TestCreatePostgresDatabase(t *testing.T) {
	t.Parallel()

	// Test with invalid connection string
	options := dbtypes.DatabaseOptions{
		ConnectionString: "invalid",
	}
	_, err := database.CreateDatabase(dbtypes.DatabaseTypePostgres, options)
	assert.Error(t, err, "Expected error for invalid connection string")
}

func TestCreateDatabase(t *testing.T) {
	t.Parallel()

	// Test with DGraph type
	db, err := database.CreateDatabase(dbtypes.DatabaseTypeDGraph, dbtypes.DatabaseOptions{
		ConnectionString: "dgraph://localhost:9080",
	})
	if err != nil {
		t.Fatalf("Failed to create DGraph database: %v", err)
	}
	if db == nil {
		t.Fatal("DGraph database should not be nil")
	}

	// Test with InMemory type
	db, err = database.CreateDatabase(dbtypes.DatabaseTypeInMemory, dbtypes.DatabaseOptions{})
	if err != nil {
		t.Fatalf("Failed to create in-memory database: %v", err)
	}
	if db == nil {
		t.Fatal("In-memory database should not be nil")
	}

	// Test with unsupported type
	db, err = database.CreateDatabase("unsupported", dbtypes.DatabaseOptions{})
	if err == nil {
		t.Fatal("Expected error for unsupported database type")
	}
	if db != nil {
		t.Fatal("Database should be nil for unsupported type")
	}
}
