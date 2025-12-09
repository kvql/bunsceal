package schemaValidation

import (
	"testing"
)

// TestSchemaPath is the relative path from pkg/taxonomy/application/validation to the schema directory
const TestSchemaPath = "../../../pkg/domain/schemas"

func MustCreateValidator(t *testing.T) *SchemaValidator {
	t.Helper()
	return MustCreateValidatorWithPath(t, TestSchemaPath)
}

func MustCreateValidatorWithPath(t *testing.T, schemaPath string) *SchemaValidator {
	t.Helper()
	validator, err := NewSchemaValidator(schemaPath)
	if err != nil {
		t.Fatalf("Failed to create schema validator with path %s: %v", schemaPath, err)
	}
	return validator
}
