// Package database provides database access and abstraction.
package database

import (
	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

// Re-export common types and errors from dbtypes.
type (
	// Options contains options for configuring a database connection.
	Options = dbtypes.DatabaseOptions

	// Type represents the type of database to use.
	Type = dbtypes.DatabaseType
)

// Database type constants.
const (
	// DatabaseTypeDGraph represents a DGraph database.
	DatabaseTypeDGraph = dbtypes.DatabaseTypeDGraph

	// DatabaseTypeInMemory represents an in-memory database (for testing).
	DatabaseTypeInMemory = dbtypes.DatabaseTypeInMemory
)

// Common errors.
var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = dbtypes.ErrNotFound

	// ErrInvalidID is returned when an invalid ID is provided.
	ErrInvalidID = dbtypes.ErrInvalidID

	// ErrAlreadyExists is returned when an entity already exists.
	ErrAlreadyExists = dbtypes.ErrAlreadyExists

	// ErrEmailRequired is returned when an email is required but not provided.
	ErrEmailRequired = dbtypes.ErrEmailRequired
)
