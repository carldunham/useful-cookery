package dbtypes

import (
	"errors"
)

// Common database errors.
var (
	// ErrNotFound is returned when an entity is not found.
	ErrNotFound = errors.New("entity not found")

	// ErrInvalidID is returned when an invalid ID is provided.
	ErrInvalidID = errors.New("invalid ID")

	// ErrAlreadyExists is returned when an entity already exists.
	ErrAlreadyExists = errors.New("entity already exists")

	// ErrEmailRequired is returned when an email is required but not provided.
	ErrEmailRequired = errors.New("email is required")
)
