package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

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

	t.Run("Loads correct count of unique files", func(t *testing.T) {
		t.Skip("TODO: LoadSegL2Files has hardcoded './schema' path - requires refactoring to accept schema path parameter")
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Domain 1", ID: "domain1"},
			{Name: "Domain 2", ID: "domain2"},
			{Name: "Domain 3", ID: "domain3"},
		})

		segL2s, err := LoadSegL2Files(tmpDir)
		testhelpers.AssertNoError(t, err, "LoadSegL2Files")
		testhelpers.AssertMapLength(t, segL2s, 3, "SegL2 map")
	})

	t.Run("Validates uniqueness of SegL2 IDs", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Domain 1", ID: "duplicate"},
			{Name: "Domain 2", ID: "duplicate"},
		})

		segL2s, err := LoadSegL2Files(tmpDir)
		testhelpers.AssertError(t, err, "Expected error for duplicate IDs")

		// Verify map is nil or smaller than file count (duplicates rejected)
		if segL2s != nil && len(segL2s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL2s))
		}
	})

	t.Run("Validates uniqueness of SegL2 names", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Duplicate Name", ID: "domain1"},
			{Name: "Duplicate Name", ID: "domain2"},
		})

		segL2s, err := LoadSegL2Files(tmpDir)
		testhelpers.AssertError(t, err, "Expected error for duplicate names")

		// Verify map is nil or smaller than file count (duplicates rejected)
		if segL2s != nil && len(segL2s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL2s))
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
		segL2, err := parseSDFile(tmpFile, validator)

		testhelpers.AssertNoError(t, err, "parseSDFile")
		testhelpers.AssertEqual(t, segL2.ID, "test", "SegL2 ID")
		testhelpers.AssertEqual(t, segL2.Name, "Test Domain", "SegL2 Name")
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
		_, err := parseSDFile(tmpFile, validator)

		testhelpers.AssertError(t, err, "Expected error for unsupported version")
	})
}
