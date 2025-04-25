package strategies_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
	"github.com/carldunham/useful-cookery/internal/database/strategies"
)

// TestNewDGraphDatabase tests the NewDGraphDatabase function.
func TestNewDGraphDatabase(t *testing.T) {
	t.Parallel()
	// Skip this test in short mode
	if testing.Short() {
		t.Skip("Skipping DGraph database test in short mode")
	}

	// Test with valid options
	options := dbtypes.DatabaseOptions{
		ConnectionString: "dgraph://localhost:9080",
	}
	db, err := strategies.NewDGraphDatabase(options)
	require.NoError(t, err)
	assert.NotNil(t, db)

	// Test with invalid connection string
	options.ConnectionString = "invalid"
	_, err = strategies.NewDGraphDatabase(options)
	assert.Error(t, err)
}

// TestBuildVarDeclarations tests the buildVarDeclarations function.
func TestBuildVarDeclarations(t *testing.T) {
	t.Parallel()
	// Test with empty vars
	result := strategies.BuildVarDeclarations(nil)
	assert.Empty(t, result)

	result = strategies.BuildVarDeclarations(map[string]string{})
	assert.Empty(t, result)

	// Test with one var
	vars := map[string]string{
		"var1": "value1",
	}
	result = strategies.BuildVarDeclarations(vars)
	assert.Equal(t, "$var1: string", result)

	// Test with multiple vars
	vars = map[string]string{
		"var1": "value1",
		"var2": "value2",
		"var3": "value3",
	}
	result = strategies.BuildVarDeclarations(vars)
	// The order of vars in the result is not guaranteed, so we need to check for each var
	assert.Contains(t, result, "$var1: string")
	assert.Contains(t, result, "$var2: string")
	assert.Contains(t, result, "$var3: string")
	// Count the number of commas by counting the occurrences of ", "
	commaCount := 0
	for i, j := 0, len(result)-1; i < j; i++ {
		if result[i:i+2] == ", " {
			commaCount++
		}
	}
	assert.Equal(t, 2, commaCount) // There should be 2 commas for 3 variables
}
