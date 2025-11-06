package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

func TestFileSegL1Repository(t *testing.T) {

	t.Run("Fails with non-existent directory", func(t *testing.T) {
		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL1Repository(validator)

		_, err := repository.LoadAll("/non/existent/path")
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Loads correct count of files", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Environment 1", ID: "env-one", Sensitivity: "A", Criticality: "1"},
			{Name: "Environment 2", ID: "env-two", Sensitivity: "B", Criticality: "2"},
			{Name: "Environment 3", ID: "env-three", Sensitivity: "C", Criticality: "3"},
		})

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL1Repository(validator)
		segL1s, err := repository.LoadAll(tmpDir)

		if err != nil {
			t.Fatalf("LoadAll: unexpected error: %v", err)
		}
		if len(segL1s) != 3 {
			t.Errorf("Expected 3 SegL1s, got %d", len(segL1s))
		}
	})
}

func TestFileSegL2Repository(t *testing.T) {
	t.Run("Fails with non-existent directory", func(t *testing.T) {
		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL2Repository(validator)

		_, err := repository.LoadAll("/non/existent/path")
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Loads correct count of files", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Domain 1", ID: "domain1"},
			{Name: "Domain 2", ID: "domain2"},
			{Name: "Domain 3", ID: "domain3"},
		})

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL2Repository(validator)
		segL2s, err := repository.LoadAll(tmpDir)

		if err != nil {
			t.Fatalf("LoadAll: unexpected error: %v", err)
		}
		if len(segL2s) != 3 {
			t.Errorf("Expected 3 SegL2s, got %d", len(segL2s))
		}
	})
}

func TestParseSegL1File(t *testing.T) {
	t.Run("Successfully parses valid SegL1 file", func(t *testing.T) {
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
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("segl1", validYAML)

		validator := mustCreateValidator(t, "../../schema")
		segL1, err := parseSegL1File(tmpFile, validator)

		if err != nil {
			t.Fatalf("parseSegL1File: unexpected error: %v", err)
		}
		if segL1.ID != "test" {
			t.Errorf("SegL1 ID: got %v, want test", segL1.ID)
		}
		if segL1.Name != "Test Environment" {
			t.Errorf("SegL1 Name: got %v, want Test Environment", segL1.Name)
		}
	})

	t.Run("Fails validation for invalid data", func(t *testing.T) {
		invalidYAML := `name: "Test"
id: "test"
# Missing description and other required fields
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("segl1", invalidYAML)

		validator := mustCreateValidator(t, "../../schema")
		_, err := parseSegL1File(tmpFile, validator)

		if err == nil {
			t.Error("Expected validation error for invalid data but got nil")
		}
	})
}

func TestParseSDFile(t *testing.T) {
	t.Run("Successfully parses valid SegL2 file", func(t *testing.T) {
		validYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain"
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    compliance_reqs:
      - pci-dss
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("segl2", validYAML)

		validator := mustCreateValidator(t, "../../schema")
		fl := NewFileSegL2Repository(validator)
		segL2, err := fl.parseSegL2File(tmpFile)

		if err != nil {
			t.Fatalf("parseSDFile: unexpected error: %v", err)
		}
		if segL2.ID != "test" {
			t.Errorf("SegL2 ID: got %v, want test", segL2.ID)
		}
		if segL2.Name != "Test Domain" {
			t.Errorf("SegL2 Name: got %v, want Test Domain", segL2.Name)
		}
		if segL2.Prominence != 1 {
			t.Errorf("SegL2 prominence got %v, expected default 1", segL2.Prominence)
		}
	})
	t.Run("Setting prominence", func(t *testing.T) {
		validYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain"
prominence: 2
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("segl2", validYAML)

		validator := mustCreateValidator(t, "../../schema")
		fl := NewFileSegL2Repository(validator)
		segL2, err := fl.parseSegL2File(tmpFile)

		if err != nil {
			t.Fatalf("parseSDFile: unexpected error: %v", err)
		}
		if segL2.Prominence != 2 {
			t.Errorf("SegL2 prominence got %v, expected default 2", segL2.Prominence)
		}
	})

	t.Run("Fails with unsupported version", func(t *testing.T) {
		invalidYAML := `version: "99.0"
name: "Test Domain"
id: "test"
description: "Test domain"
env_details:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length."
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("segl2", invalidYAML)

		validator := mustCreateValidator(t, "../../schema")
		fl := NewFileSegL2Repository(validator)
		_, err := fl.parseSegL2File(tmpFile)

		if err == nil {
			t.Error("Expected error for unsupported version but got nil")
		}
	})
}
