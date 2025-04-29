package database_test

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carldunham/useful-cookery/internal/database"
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

// createMockDB creates a mock database for testing.
//
//nolint:ireturn // Third-party interface
func createMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

// TestPostgresDatabaseCreation tests the NewPostgresDatabase function.
func TestPostgresDatabaseCreation(t *testing.T) {
	t.Parallel()
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping PostgreSQL database test in short mode")
	}

	// Test with invalid connection string
	options := dbtypes.DatabaseOptions{
		ConnectionString: "invalid",
	}
	_, err := database.NewPostgresDatabase(options)
	assert.Error(t, err)
}

// TestPostgresInputValidation tests the input validation logic for various methods.
func TestPostgresInputValidation(t *testing.T) {
	t.Parallel()

	// Create a mock database
	mockDB, _ := createMockDB(t)
	t.Cleanup(func() {
		mockDB.Close()
	})

	// Create a PostgresDatabase with the mock
	db := &database.PostgresDatabase{
		DB: mockDB,
	}

	// Test GetUser with empty user ID
	t.Run("GetUser", func(t *testing.T) {
		t.Parallel()
		_, err := db.GetUser(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})

	// Test GetUserByEmail with empty email
	t.Run("GetUserByEmail", func(t *testing.T) {
		t.Parallel()
		_, err := db.GetUserByEmail(t.Context(), "")
		assert.Equal(t, dbtypes.ErrEmailRequired, err)
	})

	// Test GetCategory with empty category ID
	t.Run("GetCategory", func(t *testing.T) {
		t.Parallel()
		_, err := db.GetCategory(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})

	// Test DeleteCategory with empty category ID
	t.Run("DeleteCategory", func(t *testing.T) {
		t.Parallel()
		err := db.DeleteCategory(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})

	// Test GetReview with empty review ID
	t.Run("GetReview", func(t *testing.T) {
		t.Parallel()
		_, err := db.GetReview(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})

	// Test DeleteReview with empty review ID
	t.Run("DeleteReview", func(t *testing.T) {
		t.Parallel()
		err := db.DeleteReview(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})

	// Test GetRecipe with empty recipe ID
	t.Run("GetRecipe", func(t *testing.T) {
		t.Parallel()
		_, err := db.GetRecipe(t.Context(), "")
		assert.Equal(t, dbtypes.ErrInvalidID, err)
	})
}

// TestPostgresQuery tests the Query function.
func TestPostgresQuery(t *testing.T) {
	t.Parallel()

	// Create a mock database
	mockDB, _ := createMockDB(t)
	t.Cleanup(func() {
		mockDB.Close()
	})

	// Create a PostgresDatabase with the mock
	db := &database.PostgresDatabase{
		DB: mockDB,
	}

	// Test with nil result
	err := db.Query(t.Context(), "query", nil, nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "raw queries not supported")
}
