package validation

import (
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
)

func TestUniquenessValidator_EmptySlice(t *testing.T) {
	var segments []domain.SegL1

	validations := UniquenessValidator(segments)

	if len(validations) != 0 {
		t.Errorf("Expected no validations for empty slice, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_SingleItem(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 0 {
		t.Errorf("Expected no validations for single item, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_UniqueItems(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-2", Name: "Test 2"},
		{ID: "test-3", Name: "Test 3"},
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 0 {
		t.Errorf("Expected no validations for unique items, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_DuplicateID(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-1", Name: "Test 2"}, // Duplicate ID
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 1 {
		t.Fatalf("Expected 1 validation error, got %d: %v", len(validations), validations)
	}

	expectedSubstring := "ID for Test 2 is not unique: test-1"
	if !strings.Contains(validations[0], expectedSubstring) {
		t.Errorf("Expected validation message to contain '%s', got: %s", expectedSubstring, validations[0])
	}
}

func TestUniquenessValidator_DuplicateName(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test"},
		{ID: "test-2", Name: "Test"}, // Duplicate Name
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 1 {
		t.Fatalf("Expected 1 validation error, got %d: %v", len(validations), validations)
	}

	expectedSubstring := "Name is not unique: Test"
	if !strings.Contains(validations[0], expectedSubstring) {
		t.Errorf("Expected validation message to contain '%s', got: %s", expectedSubstring, validations[0])
	}
}

func TestUniquenessValidator_MultipleDuplicateIDs(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-1", Name: "Test 2"}, // Duplicate ID
		{ID: "test-1", Name: "Test 3"}, // Another duplicate ID
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 2 {
		t.Errorf("Expected 2 validation errors for duplicate IDs, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_MultipleDuplicateNames(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test"},
		{ID: "test-2", Name: "Test"}, // Duplicate Name
		{ID: "test-3", Name: "Test"}, // Another duplicate Name
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 2 {
		t.Errorf("Expected 2 validation errors for duplicate names, got %d: %v", len(validations), validations)
	}
}

func TestUniquenessValidator_BothIDAndNameDuplicates(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-1", Name: "Test 2"}, // Duplicate ID
		{ID: "test-2", Name: "Test 2"}, // Duplicate Name
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 2 {
		t.Fatalf("Expected 2 validation errors (1 ID, 1 Name), got %d: %v", len(validations), validations)
	}

	// Verify we got both types of errors
	hasIDError := false
	hasNameError := false

	for _, validation := range validations {
		if strings.Contains(validation, "ID for") && strings.Contains(validation, "is not unique") {
			hasIDError = true
		}
		if strings.Contains(validation, "Name is not unique") {
			hasNameError = true
		}
	}

	if !hasIDError {
		t.Error("Expected to find ID validation error")
	}
	if !hasNameError {
		t.Error("Expected to find Name validation error")
	}
}

func TestUniquenessValidator_SegL2(t *testing.T) {
	// Test that the generic validator works with SegL2 as well
	segments := []domain.SegL2{
		{ID: "test-1", Name: "Test 1"},
		{ID: "test-1", Name: "Test 2"}, // Duplicate ID
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 1 {
		t.Fatalf("Expected 1 validation error for SegL2, got %d: %v", len(validations), validations)
	}

	expectedSubstring := "ID for Test 2 is not unique: test-1"
	if !strings.Contains(validations[0], expectedSubstring) {
		t.Errorf("Expected validation message to contain '%s', got: %s", expectedSubstring, validations[0])
	}
}

func TestUniquenessValidator_ErrorMessageFormat(t *testing.T) {
	segments := []domain.SegL1{
		{ID: "unique-id-123", Name: "Original Name"},
		{ID: "unique-id-123", Name: "Duplicate ID Name"},
		{ID: "different-id", Name: "Duplicate ID Name"}, // Duplicate Name
	}

	validations := UniquenessValidator(segments)

	if len(validations) != 2 {
		t.Fatalf("Expected 2 validation errors, got %d: %v", len(validations), validations)
	}

	// Verify ID error message includes both name and ID
	idError := validations[0]
	if !strings.Contains(idError, "Duplicate ID Name") {
		t.Errorf("ID error should include the name of the duplicate item: %s", idError)
	}
	if !strings.Contains(idError, "unique-id-123") {
		t.Errorf("ID error should include the duplicate ID: %s", idError)
	}

	// Verify Name error message includes the name
	nameError := validations[1]
	if !strings.Contains(nameError, "Duplicate ID Name") {
		t.Errorf("Name error should include the duplicate name: %s", nameError)
	}
}
