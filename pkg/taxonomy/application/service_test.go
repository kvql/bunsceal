package application

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/domain/testhelpers"
	"github.com/kvql/bunsceal/pkg/taxonomy/infrastructure"
)

func TestSegL1Service(t *testing.T) {
	validator := schemaValidation.MustCreateValidator(t)
	t.Run("Successfully loads and validates SegL1 files", func(t *testing.T) {
		repository := infrastructure.NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)

		segL1s, err := service.Load("../../../example/taxonomy/environments")
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

	t.Run("Returns map indexed by ID", func(t *testing.T) {
		files := infrastructure.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSegL1("env-one", "Environment 1", "A", "1", []string{}),
			testhelpers.NewSegL1("env-two", "Environment 2", "B", "2", []string{}),
			testhelpers.NewSegL1("env-three", "Environment 3", "C", "3", []string{}),
		})

		repository := infrastructure.NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)
		segL1s, err := service.Load(tmpDir)

		if err != nil {
			t.Fatalf("LoadAndValidate: unexpected error: %v", err)
		}
		if len(segL1s) != 3 {
			t.Errorf("Expected map length 3, got %d", len(segL1s))
		}

		// Verify map is indexed by ID
		if _, ok := segL1s["env-one"]; !ok {
			t.Error("Expected map to be indexed by ID 'env-one'")
		}
		if _, ok := segL1s["env-two"]; !ok {
			t.Error("Expected map to be indexed by ID 'env-two'")
		}
	})

	t.Run("Validates uniqueness of SegL1 IDs", func(t *testing.T) {
		files := infrastructure.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSegL1("duplicate", "Environment 1", "A", "1", []string{}),
			testhelpers.NewSegL1("duplicate", "Environment 2", "B", "2", []string{}),
		})

		repository := infrastructure.NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)
		segL1s, err := service.Load(tmpDir)

		if err == nil {
			t.Error("Expected error for duplicate IDs but got nil")
		}

		// Verify map is nil or empty (duplicates rejected)
		if len(segL1s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL1s))
		}
	})
}

func TestSegService(t *testing.T) {
	validator := schemaValidation.MustCreateValidator(t)
	t.Run("Successfully loads and validates Seg files", func(t *testing.T) {
		repository := infrastructure.NewFileSegRepository(validator)
		service := NewSegService(repository)

		Segs, err := service.Load("../../../example/taxonomy/segments")
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(Segs) == 0 {
			t.Error("Expected at least one Seg to be loaded")
		}
	})

	t.Run("Returns map indexed by ID", func(t *testing.T) {
		files := infrastructure.NewTestFiles(t)
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

		repository := infrastructure.NewFileSegRepository(validator)
		service := NewSegService(repository)
		Segs, err := service.Load(tmpDir)

		if err != nil {
			t.Fatalf("LoadAndValidate: unexpected error: %v", err)
		}
		if len(Segs) != 3 {
			t.Errorf("Expected map length 3, got %d", len(Segs))
		}

		// Verify map is indexed by ID
		if _, ok := Segs["domain-one"]; !ok {
			t.Error("Expected map to be indexed by ID 'domain-one'")
		}
	})

	t.Run("Validates uniqueness of Seg IDs", func(t *testing.T) {
		files := infrastructure.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]domain.Seg{
			testhelpers.NewSeg("duplicate", "Domain 1", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
			}),
			testhelpers.NewSeg("duplicate", "Domain 2", map[string]domain.L1Overrides{
				"production": testhelpers.NewL1Override("A", "1", []string{}),
			}),
		})

		repository := infrastructure.NewFileSegRepository(validator)
		service := NewSegService(repository)
		Segs, err := service.Load(tmpDir)

		if err == nil {
			t.Error("Expected error for duplicate IDs but got nil")
		}

		// Verify map is nil or empty (duplicates rejected)
		if len(Segs) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(Segs))
		}
	})

}
