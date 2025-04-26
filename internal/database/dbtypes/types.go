package dbtypes

import (
	"context"
	"time"
)

// DatabaseOptions contains options for configuring a database connection.
type DatabaseOptions struct {
	// ConnectionString is the connection string for the database.
	ConnectionString string

	// CacheEnabled determines whether caching is enabled.
	CacheEnabled bool

	// CacheTTL is the time-to-live for cached items.
	CacheTTL time.Duration

	// Cache is an optional cache implementation.
	Cache Cache
}

// Cache interface defines common caching operations.
type Cache interface {
	Get(ctx context.Context, key string) ([]byte, error)
	Set(ctx context.Context, key string, value []byte, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// DatabaseType represents the type of database to use.
type DatabaseType string

const (
	// DatabaseTypeDGraph represents a DGraph database.
	DatabaseTypeDGraph DatabaseType = "dgraph"

	// DatabaseTypeInMemory represents an in-memory database (for testing).
	DatabaseTypeInMemory DatabaseType = "memory"

	// DatabaseTypePostgres represents a PostgreSQL database.
	DatabaseTypePostgres DatabaseType = "postgres"

	// Add more database types as needed.
)
