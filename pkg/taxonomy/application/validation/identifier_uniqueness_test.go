package validation

import (
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
)

func TestUniquenessValidator_UniqueItems(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-2", Name: "Test 2"},
		{ID: "test-3", Name: "Test 3"},
	}

	validations := IdentifierUniquenessValidation(segments)

	if len(validations) != 0 {
		t.Errorf("Expected no validations for unique items, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_DuplicateID(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-1", Name: "Test 2"}, // Duplicate ID
	}

	validations := IdentifierUniquenessValidation(segments)

	if len(validations) != 1 {
		t.Fatalf("Expected 1 validation error, got %d: %v", len(validations), validations)
	}

	expectedSubstring := "ID for Test 2 is not unique: test-1"
	if !strings.Contains(validations[0], expectedSubstring) {
		t.Errorf("Expected validation message to contain '%s', got: %s", expectedSubstring, validations[0])
	}
}
