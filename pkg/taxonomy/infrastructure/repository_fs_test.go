package infrastructure

import (
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

func TestFileSegL1Repository(t *testing.T) {

	t.Run("Fails with non-existent directory", func(t *testing.T) {
		validator := schemaValidation.MustCreateValidator(t)
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

		validator := schemaValidation.MustCreateValidator(t)
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

func TestFileSegRepository(t *testing.T) {
	t.Run("Fails with non-existent directory", func(t *testing.T) {
		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator)

		_, err := repository.LoadAll("/non/existent/path")
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Loads correct count of files", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]testhelpers.SegFixture{
			{Name: "Domain 1", ID: "domain-one"},
			{Name: "Domain 2", ID: "domain-two"},
			{Name: "Domain 3", ID: "domain-three"},
		})

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator)
		Segs, err := repository.LoadAll(tmpDir)

		if err != nil {
			t.Fatalf("LoadAll: unexpected error: %v", err)
		}
		if len(Segs) != 3 {
			t.Errorf("Expected 3 Segs, got %d", len(Segs))
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

		validator := schemaValidation.MustCreateValidator(t)
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

		validator := schemaValidation.MustCreateValidator(t)
		_, err := parseSegL1File(tmpFile, validator)

		if err == nil {
			t.Error("Expected validation error for invalid data but got nil")
		}
	})
}

func TestParseSDFile(t *testing.T) {
	t.Run("Successfully parses valid Seg file", func(t *testing.T) {
		validYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain for validating file parsing and schema validation requirements"
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
l1_parents:
  - production
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
		tmpFile := files.CreateYAMLFile("Seg", validYAML)

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		Seg, err := fl.parseSegFile(tmpFile)

		if err != nil {
			t.Fatalf("parseSDFile: unexpected error: %v", err)
		}
		if Seg.ID != "test" {
			t.Errorf("Seg ID: got %v, want test", Seg.ID)
		}
		if Seg.Name != "Test Domain" {
			t.Errorf("Seg Name: got %v, want Test Domain", Seg.Name)
		}
		if Seg.Prominence != 1 {
			t.Errorf("Seg prominence got %v, expected default 1", Seg.Prominence)
		}
	})
	t.Run("Setting prominence", func(t *testing.T) {
		validYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain for validating prominence settings and configuration options"
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
prominence: 2
l1_parents:
  - production
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length. Test rationale with sufficient length."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length. Test rationale with sufficient length."
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("Seg", validYAML)

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		Seg, err := fl.parseSegFile(tmpFile)

		if err != nil {
			t.Fatalf("parseSDFile: unexpected error: %v", err)
		}
		if Seg.Prominence != 2 {
			t.Errorf("Seg prominence got %v, expected default 2", Seg.Prominence)
		}
	})

	t.Run("Fails with unsupported version", func(t *testing.T) {
		invalidYAML := `version: "99.0"
name: "Test Domain"
id: "test"
description: "Test domain for validating unsupported version error handling and detection"
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length."
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("Seg", invalidYAML)

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		_, err := fl.parseSegFile(tmpFile)

		if err == nil {
			t.Error("Expected error for unsupported version but got nil")
		}
	})

	t.Run("Fails when l1_overrides key not in l1_parents", func(t *testing.T) {
		invalidYAML := `version: "1.0"
name: "Test Domain"
id: "test"
description: "Test security domain for validating l1_overrides key consistency with parents"
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
l1_parents:
  - production
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
  staging:
    sensitivity: "D"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
    criticality: "5"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
`
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("Seg", invalidYAML)

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		_, err := fl.parseSegFile(tmpFile)

		if err == nil {
			t.Error("Expected L1 consistency validation error but got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "L1 consistency validation failed") {
			t.Errorf("Expected L1 consistency error, got: %v", err)
		}
	})
}
