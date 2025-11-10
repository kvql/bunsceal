package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
)

// testSchemaPath is the relative path from pkg/taxonomy to the schema directory
const testSchemaPath = "../../pkg/domain/schemas"

// mustCreateValidator creates a validator with the default test schema path
func mustCreateValidator(t *testing.T) *validation.SchemaValidator {
	t.Helper()
	return mustCreateValidatorWithPath(t, testSchemaPath)
}

// mustCreateValidatorWithPath creates a validator with a custom schema path
// Use this when you need a validator with a different schema directory
func mustCreateValidatorWithPath(t *testing.T, schemaPath string) *validation.SchemaValidator {
	t.Helper()

	validator, err := validation.NewSchemaValidator(schemaPath)
	if err != nil {
		t.Fatalf("Failed to create schema validator with path %s: %v", schemaPath, err)
	}

	return validator
}
