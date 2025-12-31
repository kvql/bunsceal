package infrastructure

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/domain/testhelpers"
)

// helper to create a config pointing to tmpDir for a specific level
func testConfig(tmpDir string) ConfigFsReposistory {
	return ConfigFsReposistory{
		TaxonomyDir: tmpDir,
		L1Dir:       "",
		L2Dir:       "",
	}
}

func TestFileSegRepository_LoadLevel(t *testing.T) {
	t.Run("Fails with non-existent directory", func(t *testing.T) {
		cfg := ConfigFsReposistory{
			TaxonomyDir: "/non/existent/path",
			L1Dir:       "l1",
			L2Dir:       "l2",
		}
		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, cfg)

		_, err := repository.LoadLevel("1")
		if err == nil {
			t.Error("Expected error for non-existent directory")
		}
	})

	t.Run("Loads correct count of L1 files", func(t *testing.T) {
		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSegL1("env-one", "Environment 1", "A", "1", nil),
			testhelpers.NewSegL1("env-two", "Environment 2", "B", "2", nil),
			testhelpers.NewSegL1("env-three", "Environment 3", "C", "3", nil),
		})

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		segs, err := repository.LoadLevel("1")

		if err != nil {
			t.Fatalf("LoadLevel: unexpected error: %v", err)
		}
		if len(segs) != 3 {
			t.Errorf("Expected 3 segments, got %d", len(segs))
		}
	})
	t.Run("Fail if level field doesn't match defined level", func(t *testing.T) {
		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSegL1("env-one", "Environment 1", "A", "1", nil),
		})

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		_, err := repository.LoadLevel("2")

		if err == nil {
			t.Error("Expected error for L1 file loaded as L2")
		}
	})

	t.Run("Loads correct count of L2 files", func(t *testing.T) {
		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSeg("domain-one", "Domain 1", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", nil),
			}),
			testhelpers.NewSeg("domain-two", "Domain 2", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", nil),
			}),
			testhelpers.NewSeg("domain-three", "Domain 3", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", nil),
			}),
		})

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		segs, err := repository.LoadLevel("2")

		if err != nil {
			t.Fatalf("LoadLevel: unexpected error: %v", err)
		}
		if len(segs) != 3 {
			t.Errorf("Expected 3 segments, got %d", len(segs))
		}
	})
}

func TestFileSegRepository_parseSegFile(t *testing.T) {
	t.Run("Successfully parses valid L1 file", func(t *testing.T) {
		seg := testhelpers.NewSegL1("test", "Test Environment", "A", "1", nil)

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		result, err := repository.parseSegFile(tmpFile, "1")

		if err != nil {
			t.Fatalf("parseSegFile: unexpected error: %v", err)
		}
		if result.ID != "test" {
			t.Errorf("Seg ID: got %v, want test", result.ID)
		}
		if result.Name != "Test Environment" {
			t.Errorf("Seg Name: got %v, want Test Environment", result.Name)
		}
		if result.Level != "1" {
			t.Errorf("Seg Level: got %v, want 1", result.Level)
		}
	})

	t.Run("Successfully parses valid L2 file", func(t *testing.T) {
		seg := testhelpers.NewSeg("test", "Test Domain", map[string]domain.L1Overrides{
			"production": testhelpers.NewL1Override("A", "1", nil),
		})

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		result, err := repository.parseSegFile(tmpFile, "2")

		if err != nil {
			t.Fatalf("parseSegFile: unexpected error: %v", err)
		}
		if result.ID != "test" {
			t.Errorf("Seg ID: got %v, want test", result.ID)
		}
		if result.Name != "Test Domain" {
			t.Errorf("Seg Name: got %v, want Test Domain", result.Name)
		}
		if result.Level != "2" {
			t.Errorf("Seg Level: got %v, want 2", result.Level)
		}
		if result.Prominence != 1 {
			t.Errorf("Seg prominence: got %v, want 1", result.Prominence)
		}
	})

	t.Run("Preserves prominence when set", func(t *testing.T) {
		seg := testhelpers.NewSeg("test", "Test Domain", map[string]domain.L1Overrides{
			"production": testhelpers.NewL1Override("A", "1", nil),
		})
		seg.Prominence = 2

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		result, err := repository.parseSegFile(tmpFile, "2")

		if err != nil {
			t.Fatalf("parseSegFile: unexpected error: %v", err)
		}
		if result.Prominence != 2 {
			t.Errorf("Seg prominence: got %v, want 2", result.Prominence)
		}
	})

	t.Run("Fails validation for invalid data", func(t *testing.T) {
		invalidYAML := `name: "Test"
id: "test"
# Missing description and other required fields
`
		files := NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("seg", invalidYAML)

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, ConfigFsReposistory{})
		_, err := repository.parseSegFile(tmpFile, "1")

		if err == nil {
			t.Error("Expected validation error but got nil")
		}
	})

	t.Run("Fails when l1_overrides key not in l1_parents", func(t *testing.T) {
		seg := testhelpers.NewSegWithParents("test", "Test Domain",
			[]string{"production"},
			map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", nil),
				"staging":    testhelpers.NewL1Override("D", "5", nil), // NOT in parents
			},
		)

		files := NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{seg})
		tmpFile := tmpDir + "/seg-0.yaml"

		validator := schemaValidation.MustCreateValidator(t)
		repository := NewFileSegRepository(validator, testConfig(tmpDir))
		_, err := repository.parseSegFile(tmpFile, "2")

		if err == nil {
			t.Error("Expected validation error but got nil")
		}
	})
}
