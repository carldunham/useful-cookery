package database

import (
	"errors"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

// ErrUnsupportedDatabaseType is returned when an unsupported database type is specified.
var ErrUnsupportedDatabaseType = errors.New("unsupported database type")

// CreateDatabase creates a database instance based on the specified type.
func CreateDatabase(dbType dbtypes.DatabaseType, options dbtypes.DatabaseOptions) (any, error) {
	switch dbType {
	case dbtypes.DatabaseTypeDGraph, "":
		return NewDGraphDatabase(options)
	case dbtypes.DatabaseTypeInMemory:
		return NewInMemoryDatabase(options)
	case dbtypes.DatabaseTypePostgres:
		return NewPostgresDatabase(options)
	default:
		return nil, ErrUnsupportedDatabaseType
	}
}
