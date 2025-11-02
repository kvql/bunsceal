package validation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/testdata"
	"gopkg.in/yaml.v3"
)

// Test helpers
func mustCreateValidator(t *testing.T, schemaPath string) *SchemaValidator {
	t.Helper()
	validator, err := NewSchemaValidator(schemaPath)
	if err != nil {
		t.Fatalf("Failed to create schema validator with path %s: %v", schemaPath, err)
	}
	return validator
}

func expectValidatorError(t *testing.T, schemaPath string) {
	t.Helper()
	_, err := NewSchemaValidator(schemaPath)
	if err == nil {
		t.Errorf("Expected validator creation to fail for path %s, but it succeeded", schemaPath)
	}
}

func assertValidationPasses(t *testing.T, data interface{}, schemaFile string) {
	t.Helper()
	validator := mustCreateValidator(t, "../../../schema")
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal test data to YAML: %v", err)
	}
	err = validator.ValidateData(yamlData, schemaFile)
	if err != nil {
		t.Errorf("Expected validation to pass for schema %s, got error: %v", schemaFile, err)
	}
}

func assertValidationFails(t *testing.T, data interface{}, schemaFile string) {
	t.Helper()
	validator := mustCreateValidator(t, "../../../schema")
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal test data to YAML: %v", err)
	}
	err = validator.ValidateData(yamlData, schemaFile)
	if err == nil {
		t.Errorf("Expected validation to fail for schema %s, but it passed", schemaFile)
	}
}

// Helper type and function to eliminate duplication in SegL2 tests
type segL2TestData struct {
	Version     string                          `yaml:"version"`
	Name        string                          `yaml:"name,omitempty"`
	ID          string                          `yaml:"id"`
	Description string                          `yaml:"description"`
	L1Overrides map[string]testdata.L1Overrides `yaml:"l1_overrides,omitempty"`
}

func marshalSegL2(seg testdata.SegL2) ([]byte, error) {
	return yaml.Marshal(segL2TestData{
		Version:     "1.0",
		Name:        seg.Name,
		ID:          seg.ID,
		Description: seg.Description,
		L1Overrides: seg.L1Overrides,
	})
}

func TestNewSchemaValidator(t *testing.T) {
	t.Run("Successfully creates validator with valid schema directory", func(t *testing.T) {
		validator := mustCreateValidator(t, "../../../schema")
		if validator == nil {
			t.Fatal("Expected validator, got nil")
		}
		if len(validator.schemas) == 0 {
			t.Error("Expected schemas to be loaded")
		}
	})

	t.Run("Fails with non-existent schema directory", func(t *testing.T) {
		expectValidatorError(t, "/non/existent/path")
	})
}

func TestValidateData_Config(t *testing.T) {
	t.Run("Valid config level keys", func(t *testing.T) {
		assertValidationPasses(t, testdata.ValidConfigSchema, "config.json")
	})

	t.Run("Invalid level key fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidConfigSchema, "config.json")
	})
}

func TestValidateData_SegL1(t *testing.T) {
	t.Run("Valid SegL1 Production passes validation", func(t *testing.T) {
		assertValidationPasses(t, testdata.ValidSegL1Production, "seg-level1.json")
	})

	t.Run("Valid SegL1 Staging passes validation", func(t *testing.T) {
		assertValidationPasses(t, testdata.ValidSegL1Staging, "seg-level1.json")
	})

	t.Run("Valid SegL1 SharedService passes validation", func(t *testing.T) {
		assertValidationPasses(t, testdata.ValidSegL1SharedService, "seg-level1.json")
	})

	t.Run("Missing required name field fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_MissingName, "seg-level1.json")
	})

	t.Run("Invalid ID pattern fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_InvalidID, "seg-level1.json")
	})

	t.Run("Short description fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_ShortDescription, "seg-level1.json")
	})

	t.Run("Invalid sensitivity enum fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_InvalidSensitivity, "seg-level1.json")
	})

	t.Run("Invalid criticality enum fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_InvalidCriticality, "seg-level1.json")
	})

	t.Run("Short rationale fails validation", func(t *testing.T) {
		assertValidationFails(t, testdata.InvalidSegL1_ShortRationale, "seg-level1.json")
	})
}

func TestValidateData_SegL2(t *testing.T) {
	validator, err := NewSchemaValidator("../../../schema")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Valid SegL2 Security passes validation", func(t *testing.T) {
		data, err := marshalSegL2(testdata.ValidSegL2Security)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "seg-level2.json")
		if err != nil {
			t.Errorf("Expected valid data to pass, got error: %v", err)
		}
	})

	t.Run("Valid SegL2 Application passes validation", func(t *testing.T) {
		data, err := marshalSegL2(testdata.ValidSegL2Application)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "seg-level2.json")
		if err != nil {
			t.Errorf("Expected valid data to pass, got error: %v", err)
		}
	})

	t.Run("Missing required name field fails validation", func(t *testing.T) {
		data, err := marshalSegL2(testdata.InvalidSegL2_MissingName)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "seg-level2.json")
		if err == nil {
			t.Error("Expected validation to fail for missing name")
		}
	})

	t.Run("Invalid ID pattern fails validation", func(t *testing.T) {
		data, err := marshalSegL2(testdata.InvalidSegL2_InvalidID)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "seg-level2.json")
		if err == nil {
			t.Error("Expected validation to fail for invalid ID pattern")
		}
	})

	t.Run("Empty environment details fails validation", func(t *testing.T) {
		data, err := marshalSegL2(testdata.InvalidSegL2_NoL1Overrides)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "seg-level2.json")
		if err == nil {
			t.Error("Expected validation to fail for empty l1_overrides")
		}
	})
}

func TestValidateData_CompReqs(t *testing.T) {
	validator, err := NewSchemaValidator("../../../schema")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Valid compliance requirements pass validation", func(t *testing.T) {
		data, err := yaml.Marshal(testdata.ValidCompReqs)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "compliance-reqs.json")
		if err != nil {
			t.Errorf("Expected valid data to pass, got error: %v", err)
		}
	})

	t.Run("Missing required fields fails validation", func(t *testing.T) {
		invalid := map[string]domain.CompReq{
			"test": {
				Name: "Test",
				// Missing Description - required field
				ReqsLink: "https://example.com",
			},
		}
		data, err := yaml.Marshal(invalid)
		if err != nil {
			t.Fatalf("Failed to marshal fixture: %v", err)
		}

		err = validator.ValidateData(data, "compliance-reqs.json")
		if err == nil {
			t.Error("Expected validation to fail for missing required field")
		}
	})
}

func TestValidateData_JSON(t *testing.T) {
	validator, err := NewSchemaValidator("../../../schema")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Valid JSON data passes validation", func(t *testing.T) {
		jsonData := []byte(`{
			"name": "Test Environment",
			"id": "test",
			"description": "This is a test environment description that is long enough to meet the minimum character requirement for validation purposes.",
			"sensitivity": "A",
			"sensitivity_rationale": "This is a test sensitivity rationale that is long enough to meet the minimum requirement.",
			"criticality": "1",
			"criticality_rationale": "This is a test criticality rationale that is long enough to meet the minimum requirement."
		}`)

		err = validator.ValidateData(jsonData, "seg-level1.json")
		if err != nil {
			t.Errorf("Expected valid JSON to pass, got error: %v", err)
		}
	})
}

func TestValidateData_ErrorHandling(t *testing.T) {
	validator, err := NewSchemaValidator("../../../schema")
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Malformed data fails parsing", func(t *testing.T) {
		malformed := []byte("this is not valid yaml or json{[")

		err = validator.ValidateData(malformed, "seg-level1.json")
		if err == nil {
			t.Error("Expected error for malformed data")
		}
	})

	t.Run("Non-existent schema file fails", func(t *testing.T) {
		data, _ := yaml.Marshal(testdata.ValidSegL1Production)

		err = validator.ValidateData(data, "non-existent.json")
		if err == nil {
			t.Error("Expected error for non-existent schema file")
		}
	})
}

func TestConvertYAMLToJSON(t *testing.T) {
	t.Run("Converts map[interface{}]interface{} to map[string]interface{}", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"key1": "value1",
			"key2": 123,
			"key3": true,
		}

		result := convertYAMLToJSON(input)

		resultMap, ok := result.(map[string]interface{})
		if !ok {
			t.Error("Expected result to be map[string]interface{}")
		}
		if resultMap["key1"] != "value1" {
			t.Errorf("Expected key1='value1', got %v", resultMap["key1"])
		}
		if resultMap["key2"] != 123 {
			t.Errorf("Expected key2=123, got %v", resultMap["key2"])
		}
	})

	t.Run("Recursively converts nested maps", func(t *testing.T) {
		input := map[interface{}]interface{}{
			"outer": map[interface{}]interface{}{
				"inner": "value",
			},
		}

		result := convertYAMLToJSON(input)

		resultMap := result.(map[string]interface{})
		innerMap, ok := resultMap["outer"].(map[string]interface{})
		if !ok {
			t.Error("Expected nested map to be converted")
		}
		if innerMap["inner"] != "value" {
			t.Errorf("Expected inner='value', got %v", innerMap["inner"])
		}
	})

	t.Run("Converts array elements recursively", func(t *testing.T) {
		input := []interface{}{
			map[interface{}]interface{}{
				"key": "value",
			},
			"string",
			123,
		}

		result := convertYAMLToJSON(input)

		resultArr, ok := result.([]interface{})
		if !ok {
			t.Error("Expected result to be array")
		}
		if len(resultArr) != 3 {
			t.Errorf("Expected 3 elements, got %d", len(resultArr))
		}

		firstMap, ok := resultArr[0].(map[string]interface{})
		if !ok {
			t.Error("Expected first element to be converted map")
		}
		if firstMap["key"] != "value" {
			t.Errorf("Expected key='value', got %v", firstMap["key"])
		}
	})

	t.Run("Returns primitives unchanged", func(t *testing.T) {
		inputs := []interface{}{
			"string",
			123,
			true,
			nil,
		}

		for _, input := range inputs {
			result := convertYAMLToJSON(input)
			if result != input {
				t.Errorf("Expected %v to remain unchanged, got %v", input, result)
			}
		}
	})
}

func TestFormatValidationError(t *testing.T) {
	t.Run("Formats jsonschema.ValidationError", func(t *testing.T) {
		validator, err := NewSchemaValidator("../../../schema")
		if err != nil {
			t.Fatalf("Failed to create validator: %v", err)
		}

		// Create invalid data to trigger validation error
		invalidData := []byte(`{"name": "test"}`)

		err = validator.ValidateData(invalidData, "seg-level1.json")
		if err == nil {
			t.Fatal("Expected validation error")
		}

		// Check that error message is formatted
		errMsg := err.Error()
		if !strings.Contains(errMsg, "schema validation failed") {
			t.Errorf("Expected formatted error message, got: %s", errMsg)
		}
	})

	t.Run("Returns non-ValidationError unchanged", func(t *testing.T) {
		testErr := formatValidationError(errors.New("test error"))
		if testErr.Error() != "test error" {
			t.Errorf("Expected unchanged error message, got: %s", testErr.Error())
		}
	})
}

func TestNewSchemaValidator_ErrorCases(t *testing.T) {
	t.Run("Fails with invalid schema JSON", func(t *testing.T) {
		// Create temp directory with invalid schema
		tmpDir, err := os.MkdirTemp("", "schema-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create an invalid JSON file
		invalidJSON := `{"this is not valid json`
		if err := os.WriteFile(filepath.Join(tmpDir, "common.json"), []byte(invalidJSON), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		_, err = NewSchemaValidator(tmpDir)
		if err == nil {
			t.Error("Expected error for invalid JSON schema")
		}
	})
}
