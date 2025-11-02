package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
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

	t.Run("Loads correct count of unique files", func(t *testing.T) {
		t.Skip("TODO: LoadSegL1Files has hardcoded './schema' path - requires refactoring to accept schema path parameter")
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Environment 1", ID: "env1", Sensitivity: "A", Criticality: "1"},
			{Name: "Environment 2", ID: "env2", Sensitivity: "B", Criticality: "2"},
			{Name: "Environment 3", ID: "env3", Sensitivity: "C", Criticality: "3"},
		})

		segL1s, err := LoadSegL1Files(tmpDir)
		testhelpers.AssertNoError(t, err, "LoadSegL1Files")
		testhelpers.AssertMapLength(t, segL1s, 3, "SegL1 map")
	})

	t.Run("Validates uniqueness of SegL1 IDs", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Environment 1", ID: "duplicate", Sensitivity: "A", Criticality: "1"},
			{Name: "Environment 2", ID: "duplicate", Sensitivity: "B", Criticality: "2"},
		})

		segL1s, err := LoadSegL1Files(tmpDir)
		testhelpers.AssertError(t, err, "Expected error for duplicate IDs")

		// Verify map is nil or smaller than file count (duplicates rejected)
		if segL1s != nil && len(segL1s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL1s))
		}
	})

	t.Run("Validates uniqueness of SegL1 names", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Duplicate Name", ID: "env1", Sensitivity: "A", Criticality: "1"},
			{Name: "Duplicate Name", ID: "env2", Sensitivity: "B", Criticality: "2"},
		})

		segL1s, err := LoadSegL1Files(tmpDir)
		testhelpers.AssertError(t, err, "Expected error for duplicate names")

		// Verify map is nil or smaller than file count (duplicates rejected)
		if segL1s != nil && len(segL1s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL1s))
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

		testhelpers.AssertNoError(t, err, "parseSegL1File")
		testhelpers.AssertEqual(t, segL1.ID, "test", "SegL1 ID")
		testhelpers.AssertEqual(t, segL1.Name, "Test Environment", "SegL1 Name")
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

		testhelpers.AssertError(t, err, "Expected validation error for invalid data")
	})
}
