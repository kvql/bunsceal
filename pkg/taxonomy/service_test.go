package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

func TestSegL1Service(t *testing.T) {
	t.Run("Successfully loads and validates SegL1 files", func(t *testing.T) {
		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)

		segL1s, err := service.LoadAndValidate("../../example/taxonomy/environments")
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

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)
		segL1s, err := service.LoadAndValidate(tmpDir)

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

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL1Repository(validator)
		service := NewSegL1Service(repository)
		segL1s, err := service.LoadAndValidate(tmpDir)

		if err == nil {
			t.Error("Expected error for duplicate IDs but got nil")
		}

		// Verify map is nil or empty (duplicates rejected)
		if len(segL1s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL1s))
		}
	})
}

func TestSegL2Service(t *testing.T) {
	t.Run("Successfully loads and validates SegL2 files", func(t *testing.T) {
		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL2Repository(validator)
		service := NewSegL2Service(repository)

		segL2s, err := service.LoadAndValidate("../../example/taxonomy/segments")
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(segL2s) == 0 {
			t.Error("Expected at least one SegL2 to be loaded")
		}
	})

	t.Run("Returns map indexed by ID", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Domain 1", ID: "domain1"},
			{Name: "Domain 2", ID: "domain2"},
			{Name: "Domain 3", ID: "domain3"},
		})

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL2Repository(validator)
		service := NewSegL2Service(repository)
		segL2s, err := service.LoadAndValidate(tmpDir)

		if err != nil {
			t.Fatalf("LoadAndValidate: unexpected error: %v", err)
		}
		if len(segL2s) != 3 {
			t.Errorf("Expected map length 3, got %d", len(segL2s))
		}

		// Verify map is indexed by ID
		if _, ok := segL2s["domain1"]; !ok {
			t.Error("Expected map to be indexed by ID 'domain1'")
		}
	})

	t.Run("Validates uniqueness of SegL2 IDs", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpDir := files.CreateSegL2Files([]testhelpers.SegL2Fixture{
			{Name: "Domain 1", ID: "duplicate"},
			{Name: "Domain 2", ID: "duplicate"},
		})

		validator := mustCreateValidator(t, "../../schema")
		repository := NewFileSegL2Repository(validator)
		service := NewSegL2Service(repository)
		segL2s, err := service.LoadAndValidate(tmpDir)

		if err == nil {
			t.Error("Expected error for duplicate IDs but got nil")
		}

		// Verify map is nil or empty (duplicates rejected)
		if len(segL2s) == 2 {
			t.Errorf("Map size %d equals file count 2 - duplicates may have been silently overwritten", len(segL2s))
		}
	})

}
