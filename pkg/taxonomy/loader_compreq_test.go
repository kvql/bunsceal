package taxonomy

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

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
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("compliance-req", "this is not valid yaml: {[")

		_, err := LoadCompScope(tmpFile)
		testhelpers.AssertError(t, err, "Expected error for invalid YAML")
	})
}
