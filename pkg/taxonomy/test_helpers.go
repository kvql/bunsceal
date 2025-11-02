package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
)

// mustCreateValidator creates a validator with a custom schema path
// Use this when you need a validator with a different schema directory
func mustCreateValidator(t *testing.T, schemaPath string) *validation.SchemaValidator {
	t.Helper()

	validator, err := validation.NewSchemaValidator(schemaPath)
	if err != nil {
		t.Fatalf("Failed to create schema validator with path %s: %v", schemaPath, err)
	}

	return validator
}
