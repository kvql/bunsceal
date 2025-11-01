package taxonomy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadSegL1Files(t *testing.T) {
	t.Run("Successfully loads valid SegL1 files", func(t *testing.T) {
		t.Skip("TODO: LoadSegL1Files has hardcoded './schema' path - requires refactoring to accept schema path parameter")
		dir := "../../example/taxonomy/security-environments"
		segL1s, err := LoadSegL1Files(dir)
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(segL1s) == 0 {
			t.Error("Expected at least one SegL1 to be loaded")
		}

		// Check that shared-service was loaded (required by validation)
		if _, ok := segL1s["shared-service"]; !ok {
			t.Error("Expected shared-service environment to be loaded")
		}
	})

	t.Run("Fails with non-existent directory", func(t *testing.T) {
		dir := "/non/existent/path"
		_, err := LoadSegL1Files(dir)
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Validates uniqueness of SegL1 IDs", func(t *testing.T) {
		// Create temporary directory with duplicate IDs
		tmpDir, err := os.MkdirTemp("", "segl1-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create two files with the same ID
		file1 := `name: "Environment 1"
id: "duplicate"
description: "This is a test environment with sufficient description length to meet minimum requirements."
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
`
		file2 := `name: "Environment 2"
id: "duplicate"
description: "This is another test environment with sufficient description length to meet minimum requirements."
sensitivity: "B"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
criticality: "2"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
`
		if err := os.WriteFile(filepath.Join(tmpDir, "env1.yaml"), []byte(file1), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "env2.yaml"), []byte(file2), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		_, err = LoadSegL1Files(tmpDir)
		if err == nil {
			t.Error("Expected error for duplicate IDs")
		}
	})
}

func TestLoadSegL2Files(t *testing.T) {
	t.Run("Successfully loads valid SegL2 files", func(t *testing.T) {
		t.Skip("TODO: LoadSegL2Files has hardcoded './schema' path - requires refactoring to accept schema path parameter")
		dir := "../../example/taxonomy/security-domains"
		segL2s, err := LoadSegL2Files(dir)
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(segL2s) == 0 {
			t.Error("Expected at least one SegL2 to be loaded")
		}
	})

	t.Run("Fails with non-existent directory", func(t *testing.T) {
		dir := "/non/existent/path"
		_, err := LoadSegL2Files(dir)
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Validates uniqueness of SegL2 IDs", func(t *testing.T) {
		// Create temporary directory with duplicate IDs
		tmpDir, err := os.MkdirTemp("", "segl2-test")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		// Create two files with the same ID
		file1 := `version: "1.0"
name: "Domain 1"
id: "duplicate"
description: "Test domain"
env_details:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
`
		file2 := `version: "1.0"
name: "Domain 2"
id: "duplicate"
description: "Another test domain"
env_details:
  staging:
    sensitivity: "D"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
    criticality: "5"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
`
		if err := os.WriteFile(filepath.Join(tmpDir, "domain1.yaml"), []byte(file1), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		if err := os.WriteFile(filepath.Join(tmpDir, "domain2.yaml"), []byte(file2), 0644); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}

		_, err = LoadSegL2Files(tmpDir)
		if err == nil {
			t.Error("Expected error for duplicate IDs")
		}
	})
}

func TestLoadCompScope(t *testing.T) {
	t.Run("Successfully loads valid compliance requirements", func(t *testing.T) {
		t.Skip("TODO: LoadCompScope has hardcoded './schema' path - requires refactoring to accept schema path parameter")
		file := "../../example/taxonomy/compliance_requirements.yaml"
		compReqs, err := LoadCompScope(file)
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(compReqs) == 0 {
			t.Error("Expected at least one compliance requirement to be loaded")
		}
	})

	t.Run("Fails with non-existent file", func(t *testing.T) {
		file := "/non/existent/file.yaml"
		_, err := LoadCompScope(file)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("Fails with invalid YAML", func(t *testing.T) {
		// Create temporary invalid file
		tmpFile, err := os.CreateTemp("", "comp-req-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		invalidYAML := "this is not valid yaml: {["
		if _, err := tmpFile.WriteString(invalidYAML); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		tmpFile.Close()

		_, err = LoadCompScope(tmpFile.Name())
		if err == nil {
			t.Error("Expected error for invalid YAML")
		}
	})
}

func TestParseSegL1File(t *testing.T) {
	t.Run("Successfully parses valid SegL1 file", func(t *testing.T) {
		// Create temporary valid file
		tmpFile, err := os.CreateTemp("", "segl1-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		validYAML := `name: "Test Environment"
id: "test"
description: "This is a test environment with sufficient description length to meet minimum requirements for validation."
sensitivity: "A"
sensitivity_rationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions."
criticality: "1"
criticality_rationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions."
compliance_reqs:
  - pci-dss
`
		if _, err := tmpFile.WriteString(validYAML); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		tmpFile.Close()

		schemaValidator, err := NewSchemaValidator("../../schema")
		if err != nil {
			t.Fatalf("Failed to create schema validator: %v", err)
		}

		segL1, err := parseSegL1File(tmpFile.Name(), schemaValidator)
		if err != nil {
			t.Errorf("Expected successful parse, got error: %v", err)
		}
		if segL1.ID != "test" {
			t.Errorf("Expected ID 'test', got %s", segL1.ID)
		}
		if segL1.Name != "Test Environment" {
			t.Errorf("Expected name 'Test Environment', got %s", segL1.Name)
		}
	})

	t.Run("Fails validation for invalid data", func(t *testing.T) {
		// Create temporary invalid file (missing required fields)
		tmpFile, err := os.CreateTemp("", "segl1-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		invalidYAML := `name: "Test"
id: "test"
# Missing description and other required fields
`
		if _, err := tmpFile.WriteString(invalidYAML); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		tmpFile.Close()

		schemaValidator, err := NewSchemaValidator("../../schema")
		if err != nil {
			t.Fatalf("Failed to create schema validator: %v", err)
		}

		_, err = parseSegL1File(tmpFile.Name(), schemaValidator)
		if err == nil {
			t.Error("Expected validation error for invalid data")
		}
	})
}

func TestParseSDFile(t *testing.T) {
	t.Run("Successfully parses valid SegL2 file", func(t *testing.T) {
		// Create temporary valid file
		tmpFile, err := os.CreateTemp("", "segl2-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		validYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain"
env_details:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    compliance_reqs:
      - pci-dss
`
		if _, err := tmpFile.WriteString(validYAML); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		tmpFile.Close()

		schemaValidator, err := NewSchemaValidator("../../schema")
		if err != nil {
			t.Fatalf("Failed to create schema validator: %v", err)
		}

		segL2, err := parseSDFile(tmpFile.Name(), schemaValidator)
		if err != nil {
			t.Errorf("Expected successful parse, got error: %v", err)
		}
		if segL2.ID != "test" {
			t.Errorf("Expected ID 'test', got %s", segL2.ID)
		}
		if segL2.Name != "Test Domain" {
			t.Errorf("Expected name 'Test Domain', got %s", segL2.Name)
		}
	})

	t.Run("Fails with unsupported version", func(t *testing.T) {
		// Create temporary file with unsupported version
		tmpFile, err := os.CreateTemp("", "segl2-*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		invalidYAML := `version: "99.0"
name: "Test Domain"
id: "test"
description: "Test domain"
env_details:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length."
`
		if _, err := tmpFile.WriteString(invalidYAML); err != nil {
			t.Fatalf("Failed to write test file: %v", err)
		}
		tmpFile.Close()

		schemaValidator, err := NewSchemaValidator("../../schema")
		if err != nil {
			t.Fatalf("Failed to create schema validator: %v", err)
		}

		_, err = parseSDFile(tmpFile.Name(), schemaValidator)
		if err == nil {
			t.Error("Expected error for unsupported version")
		}
	})
}
