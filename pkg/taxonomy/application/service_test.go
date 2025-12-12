package application

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/taxonomy/infrastructure"
	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
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
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Environment 1", ID: "env-one", Sensitivity: "A", Criticality: "1"},
			{Name: "Environment 2", ID: "env-two", Sensitivity: "B", Criticality: "2"},
			{Name: "Environment 3", ID: "env-three", Sensitivity: "C", Criticality: "3"},
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
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL1Files([]testhelpers.SegL1Fixture{
			{Name: "Environment 1", ID: "duplicate", Sensitivity: "A", Criticality: "1"},
			{Name: "Environment 2", ID: "duplicate", Sensitivity: "B", Criticality: "2"},
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
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]testhelpers.SegFixture{
			{Name: "Domain 1", ID: "domain-one"},
			{Name: "Domain 2", ID: "domain-two"},
			{Name: "Domain 3", ID: "domain-three"},
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
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegFiles([]testhelpers.SegFixture{
			{Name: "Domain 1", ID: "duplicate"},
			{Name: "Domain 2", ID: "duplicate"},
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
