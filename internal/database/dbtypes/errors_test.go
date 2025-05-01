package dbtypes_test

import (
	"errors"
	"testing"

	"github.com/carldunham/useful-cookery/internal/database/dbtypes"
)

func TestErrorsExistence(t *testing.T) {
	t.Parallel()
	// Test that error variables are defined
	if dbtypes.ErrNotFound == nil {
		t.Error("ErrNotFound should not be nil")
	}
	if dbtypes.ErrInvalidID == nil {
		t.Error("ErrInvalidID should not be nil")
	}
	if dbtypes.ErrAlreadyExists == nil {
		t.Error("ErrAlreadyExists should not be nil")
	}
	if dbtypes.ErrEmailRequired == nil {
		t.Error("ErrEmailRequired should not be nil")
	}
}

func TestErrorMessages(t *testing.T) {
	t.Parallel()
	// Test error messages
	if dbtypes.ErrNotFound.Error() != "entity not found" {
		t.Errorf("ErrNotFound has unexpected message: %s", dbtypes.ErrNotFound.Error())
	}
	if dbtypes.ErrInvalidID.Error() != "invalid ID" {
		t.Errorf("ErrInvalidID has unexpected message: %s", dbtypes.ErrInvalidID.Error())
	}
	if dbtypes.ErrAlreadyExists.Error() != "entity already exists" {
		t.Errorf("ErrAlreadyExists has unexpected message: %s", dbtypes.ErrAlreadyExists.Error())
	}
	if dbtypes.ErrEmailRequired.Error() != "email is required" {
		t.Errorf("ErrEmailRequired has unexpected message: %s", dbtypes.ErrEmailRequired.Error())
	}
}

func TestErrorComparison(t *testing.T) {
	t.Parallel()
	// Test error comparison with errors.Is
	err := dbtypes.ErrNotFound
	if !errors.Is(err, dbtypes.ErrNotFound) {
		t.Error("errors.Is should return true for the same error")
	}
	if errors.Is(err, dbtypes.ErrInvalidID) {
		t.Error("errors.Is should return false for different errors")
	}
}
