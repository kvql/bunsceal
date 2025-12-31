package schemaValidation

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/testhelpers"
	"gopkg.in/yaml.v3"
)

// Config test fixtures for schema validation testing
type Config struct {
	Terminology TermConfig `yaml:"terminology"`
}
type InvalidConfig struct {
	Terminology InvalidTermConfig `yaml:"terminology"`
}
type TermConfig struct {
	L1 TermDef `yaml:"l1,omitempty"`
	L2 TermDef `yaml:"l2,omitempty"`
}
type InvalidTermConfig struct {
	L4 TermDef `yaml:"l4"`
}
type TermDef struct {
	Singular string `yaml:"singular"`
	Plural   string `yaml:"plural"`
}

func expectValidatorError(t *testing.T, schemaPath string) {
	t.Helper()
	_, err := NewSchemaValidator(schemaPath, SchemaBaseURL)
	if err == nil {
		t.Errorf("Expected validator creation to fail for path %s, but it succeeded", schemaPath)
	}
}

func assertValidationPasses(t *testing.T, data interface{}, schemaFile string) {
	t.Helper()
	validator := MustCreateValidator(t)
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
	validator := MustCreateValidator(t)
	yamlData, err := yaml.Marshal(data)
	if err != nil {
		t.Fatalf("Failed to marshal test data to YAML: %v", err)
	}
	err = validator.ValidateData(yamlData, schemaFile)
	if err == nil {
		t.Errorf("Expected validation to fail for schema %s, but it passed", schemaFile)
	}
}

// segL1WithLabels creates a valid SegL1 with labels for testing
func segL1WithLabels(labels []string) domain.Seg {
	seg := testhelpers.NewSegL1("test", "Test", "A", "1", nil)
	seg.Labels = labels
	return seg
}

// segWithLabels creates a valid L2 Seg with labels for testing
func segWithLabels(labels []string) domain.Seg {
	seg := testhelpers.NewSeg("test", "Test", map[string]domain.L1Overrides{
		"prod": testhelpers.NewL1Override("A", "1", nil),
	})
	seg.Labels = labels
	return seg
}

func TestNewSchemaValidator(t *testing.T) {
	t.Run("Successfully creates validator with valid schema directory", func(t *testing.T) {
		validator := MustCreateValidator(t)
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

func TestValidateData_SegL1(t *testing.T) {
	t.Run("Valid SegL1 passes validation", func(t *testing.T) {
		assertValidationPasses(t, testhelpers.NewSegL1("test", "Test", "A", "1", nil), "seg-level.json")
	})

	t.Run("Missing required name field fails validation", func(t *testing.T) {
		seg := testhelpers.NewSegL1("test", "", "A", "1", nil)
		seg.Name = ""
		assertValidationFails(t, seg, "seg-level.json")
	})

	t.Run("Invalid ID pattern fails validation", func(t *testing.T) {
		seg := testhelpers.NewSegL1("Invalid_ID_With_Capitals", "Test", "A", "1", nil)
		assertValidationFails(t, seg, "seg-level.json")
	})

	t.Run("Short description fails validation", func(t *testing.T) {
		seg := testhelpers.NewSegL1("test", "Test", "A", "1", nil)
		seg.Description = "Too short"
		assertValidationFails(t, seg, "seg-level.json")
	})

	// Note: Sensitivity/Criticality enum validation removed - now handled via plugin labels
}

func TestValidateData_Seg(t *testing.T) {
	t.Run("Valid Seg passes validation", func(t *testing.T) {
		seg := testhelpers.NewSeg("test", "Test", map[string]domain.L1Overrides{
			"prod": testhelpers.NewL1Override("A", "1", nil),
		})
		assertValidationPasses(t, seg, "seg-level.json")
	})

	t.Run("Missing required name field fails validation", func(t *testing.T) {
		seg := testhelpers.NewSeg("test", "", map[string]domain.L1Overrides{
			"prod": testhelpers.NewL1Override("A", "1", nil),
		})
		seg.Name = ""
		assertValidationFails(t, seg, "seg-level.json")
	})

	t.Run("Invalid ID pattern fails validation", func(t *testing.T) {
		seg := testhelpers.NewSeg("Invalid_ID!", "Test", map[string]domain.L1Overrides{
			"prod": testhelpers.NewL1Override("A", "1", nil),
		})
		assertValidationFails(t, seg, "seg-level.json")
	})
}

func TestValidateData_JSON(t *testing.T) {
	validator, err := NewSchemaValidator(TestSchemaPath, SchemaBaseURL)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Valid JSON data passes validation", func(t *testing.T) {
		jsonData := []byte(`{
			"name": "Test Environment",
			"id": "test",
			"description": "This is a test environment description that is long enough to meet the minimum character requirement for validation purposes.",
			"labels": [
				"bunsceal.plugin.classifications/sensitivity:A",
				"bunsceal.plugin.classifications/criticality:1"
			]
		}`)

		err = validator.ValidateData(jsonData, "seg-level.json")
		if err != nil {
			t.Errorf("Expected valid JSON to pass, got error: %v", err)
		}
	})
}

func TestValidateData_ErrorHandling(t *testing.T) {
	validator, err := NewSchemaValidator(TestSchemaPath, SchemaBaseURL)
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	t.Run("Malformed data fails parsing", func(t *testing.T) {
		malformed := []byte("this is not valid yaml or json{[")

		err = validator.ValidateData(malformed, "seg-level.json")
		if err == nil {
			t.Error("Expected error for malformed data")
		}
	})

	t.Run("Non-existent schema file fails", func(t *testing.T) {
		data, _ := yaml.Marshal(testhelpers.NewSegL1("test", "Test", "A", "1", nil))

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
		validator, err := NewSchemaValidator(TestSchemaPath, SchemaBaseURL)
		if err != nil {
			t.Fatalf("Failed to create validator: %v", err)
		}

		// Create invalid data to trigger validation error
		invalidData := []byte(`{"name": "test"}`)

		err = validator.ValidateData(invalidData, "seg-level.json")
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

		_, err = NewSchemaValidator(tmpDir, SchemaBaseURL)
		if err == nil {
			t.Error("Expected error for invalid JSON schema")
		}
	})
}

func TestValidateData_Labels(t *testing.T) {
	t.Run("Valid label formats in SegL1", func(t *testing.T) {
		testCases := []struct {
			name   string
			labels []string
		}{
			{"Simple key:value pairs", []string{"env:prod", "region:us-east-1"}},
			{"With hyphens", []string{"region-id:us-west-2", "app-tier:backend"}},
			{"With underscores", []string{"team_name:platform", "cost_center:engineering"}},
			{"Namespaced labels", []string{"org.example/app:backend", "company.io/team:infra"}},
			{"Values with forward slashes", []string{"url:api.example.com/v1", "path:/var/log/app.log"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertValidationPasses(t, segL1WithLabels(tc.labels), "seg-level.json")
			})
		}
	})

	t.Run("AWS-compliant special characters in values", func(t *testing.T) {
		testCases := []struct {
			name   string
			labels []string
		}{
			{"Values with spaces", []string{"env:production environment", "desc:web application tier"}},
			{"Values with plus signs", []string{"version:v1.0.0+build.123", "tag:feature+bugfix"}},
			{"Values with equals signs", []string{"formula:x=y+z", "equation:a=b"}},
			{"Values with at signs", []string{"owner:team@example.com", "contact:admin@internal"}},
			{"Values with colons", []string{"url:https://api.example.com:8080", "time:12:30:45"}},
			{"Combined special chars", []string{"deploy:user@host:/path/to/app v1.0+patch", "ref:team=platform contact@example.com"}},
			{"Trailing/leading spaces in values", []string{"env: production ", "tier: frontend"}},
			{"Multiple spaces in values", []string{"desc:multi  word   description", "note:has    gaps"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertValidationPasses(t, segL1WithLabels(tc.labels), "seg-level.json")
			})
		}
	})

	t.Run("Invalid label formats in SegL1", func(t *testing.T) {
		testCases := []struct {
			name   string
			labels []string
		}{
			{"Missing colon separator", []string{"invalid-no-colon"}},
			{"Key starts with hyphen", []string{"-invalid:value"}},
			{"Key ends with hyphen", []string{"invalid-:value"}},
			{"Empty key", []string{":value"}},
			{"Key starts with dot", []string{".invalid:value"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertValidationFails(t, segL1WithLabels(tc.labels), "seg-level.json")
			})
		}
	})

	t.Run("Invalid special characters in values", func(t *testing.T) {
		testCases := []struct {
			name   string
			labels []string
		}{
			{"Value with brackets", []string{"list:[item1,item2]"}},
			{"Value with braces", []string{"obj:{key:val}"}},
			{"Value with percent", []string{"rate:50%"}},
			{"Value with asterisk", []string{"glob:*.txt"}},
			{"Value with hash", []string{"color:#FF0000"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertValidationFails(t, segL1WithLabels(tc.labels), "seg-level.json")
			})
		}
	})

	t.Run("Valid extended punctuation in values", func(t *testing.T) {
		testCases := []struct {
			name   string
			labels []string
		}{
			{"Value with parentheses", []string{"note:(important)"}},
			{"Value with ampersand", []string{"query:foo&bar"}},
			{"Value with apostrophe", []string{"note:doesn't matter"}},
			{"Value with comma", []string{"list:one, two, three"}},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				assertValidationPasses(t, segL1WithLabels(tc.labels), "seg-level.json")
			})
		}
	})

	t.Run("Uniqueness constraint in SegL1", func(t *testing.T) {
		assertValidationFails(t, segL1WithLabels([]string{"env:prod", "env:prod"}), "seg-level.json")
	})

	t.Run("Valid labels in Seg", func(t *testing.T) {
		assertValidationPasses(t, segWithLabels([]string{"domain:application", "tier:frontend"}), "seg-level.json")
	})

	t.Run("Optional labels field", func(t *testing.T) {
		t.Run("SegL1 without labels passes", func(t *testing.T) {
			assertValidationPasses(t, testhelpers.NewSegL1("test", "Test", "A", "1", nil), "seg-level.json")
		})

		t.Run("Seg without labels passes", func(t *testing.T) {
			seg := testhelpers.NewSeg("test", "Test", map[string]domain.L1Overrides{
				"prod": testhelpers.NewL1Override("A", "1", nil),
			})
			assertValidationPasses(t, seg, "seg-level.json")
		})
	})
}
