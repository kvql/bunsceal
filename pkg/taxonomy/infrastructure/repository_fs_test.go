package infrastructure

import (
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/domain/testhelpers"
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
		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSegL1("env-one", "Environment 1", "A", "1", []string{}),
			testhelpers.NewSegL1("env-two", "Environment 2", "B", "2", []string{}),
			testhelpers.NewSegL1("env-three", "Environment 3", "C", "3", []string{}),
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
		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSeg("domain-one", "Domain 1", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
			}),
			testhelpers.NewSeg("domain-two", "Domain 2", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
			}),
			testhelpers.NewSeg("domain-three", "Domain 3", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
			}),
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
		// Use domain builder to create valid L1 segment
		seg := testhelpers.NewSegL1("test", "Test Environment", "A", "1", []string{"pci-dss"})

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

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
		files := NewTestFiles(t)
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
		// Use domain builder to create valid L2 segment (NO top-level sensitivity/criticality)
		seg := testhelpers.NewSeg("test", "Test Domain", map[string]domain.L1Overrides{
			"production": testhelpers.NewL1Override("A", "1", []string{"pci-dss"}),
		})

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

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
		// Use domain builder and set prominence
		seg := testhelpers.NewSeg("test", "Test Domain", map[string]domain.L1Overrides{
			"production": testhelpers.NewL1Override("A", "1", []string{}),
		})
		seg.Prominence = 2

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

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
		// Testing unsupported version - L2 should NOT have top-level sensitivity/criticality
		invalidYAML := `version: "99.0"
name: "Test Domain"
id: "test"
description: "Test domain for validating unsupported version error handling and detection"
l1_parents:
  - production
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
`
		files := NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("Seg", invalidYAML)

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		_, err := fl.parseSegFile(tmpFile)

		if err == nil {
			t.Error("Expected error for unsupported version but got nil")
		}
	})

	t.Run("Fails when l1_overrides key not in l1_parents", func(t *testing.T) {
		// Create L2 with override for 'staging' that's not in l1_parents
		seg := testhelpers.NewSegWithParents("test", "Test Domain",
			[]string{"production"}, // parents list
			map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
				"staging":    testhelpers.NewL1Override("D", "5", []string{}), // NOT in parents!
			},
		)

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

		validator := schemaValidation.MustCreateValidator(t)
		fl := NewFileSegRepository(validator)
		_, err := fl.parseSegFile(tmpFile)

		if err == nil {
			t.Error("Expected L1 consistency validation error but got nil")
		}
		if err != nil && !strings.Contains(err.Error(), "PostLoad validation failed") {
			t.Errorf("Expected PostLoad validation error, got: %v", err)
		}
		if err != nil && !strings.Contains(err.Error(), "l1_overrides") {
			t.Errorf("Expected l1_overrides consistency error, got: %v", err)
		}
	})
}
