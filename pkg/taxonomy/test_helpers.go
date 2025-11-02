package taxonomy

import (
	"sync"
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"gopkg.in/yaml.v3"
)

// Test-only helpers to avoid duplication across test files

var (
	testValidator *validation.SchemaValidator
	validatorOnce sync.Once
	validatorErr  error
)

// getTestValidator returns a cached schema validator for testing
// This avoids re-reading and parsing schemas for every test
func getTestValidator(t *testing.T) *validation.SchemaValidator {
	t.Helper()

	validatorOnce.Do(func() {
		testValidator, validatorErr = validation.NewSchemaValidator("../../schema")
	})

	if validatorErr != nil {
		t.Fatalf("Failed to initialize test schema validator: %v", validatorErr)
	}

	return testValidator
}

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

// expectValidatorError expects validator creation to fail
// Useful for testing invalid schema directories
func expectValidatorError(t *testing.T, schemaPath string) {
	t.Helper()

	_, err := validation.NewSchemaValidator(schemaPath)
	if err == nil {
		t.Errorf("Expected validator creation to fail for path %s, but it succeeded", schemaPath)
	}
}

// assertValidationPasses validates that data passes schema validation
// Automatically marshals data to YAML before validation
func assertValidationPasses(t *testing.T, data interface{}, schemaFile string) {
	t.Helper()

	validator := getTestValidator(t)
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal test data to YAML: %v", err)
	}

	err = validator.ValidateData(yamlData, schemaFile)
	if err != nil {
		t.Errorf("Expected validation to pass for schema %s, got error: %v", schemaFile, err)
	}
}

// assertValidationFails validates that data fails schema validation
// Automatically marshals data to YAML before validation
func assertValidationFails(t *testing.T, data interface{}, schemaFile string) {
	t.Helper()

	validator := getTestValidator(t)
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal test data to YAML: %v", err)
	}

	err = validator.ValidateData(yamlData, schemaFile)
	if err == nil {
		t.Errorf("Expected validation to fail for schema %s, but it passed", schemaFile)
	}
}
