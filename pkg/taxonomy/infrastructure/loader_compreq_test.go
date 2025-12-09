package infrastructure

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

func TestLoadCompScope(t *testing.T) {
	validator := schemaValidation.MustCreateValidator(t)
	t.Run("Successfully loads valid compliance requirements", func(t *testing.T) {
		file := "../../../example/taxonomy/compliance_requirements.yaml"
		compReqs, err := LoadCompScope(file, validator)
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if len(compReqs) == 0 {
			t.Error("Expected at least one compliance requirement to be loaded")
		}

		// Verify at least one common compliance requirement exists
		if _, ok := compReqs["pci-dss"]; !ok {
			if _, ok := compReqs["sox"]; !ok {
				t.Error("Expected at least one standard compliance requirement (pci-dss or sox)")
			}
		}
	})

	t.Run("Fails with non-existent file", func(t *testing.T) {
		_, err := LoadCompScope("/non/existent/file.yaml", validator)
		if err == nil {
			t.Error("Expected error for non-existent file")
		}
	})

	t.Run("Fails with invalid YAML", func(t *testing.T) {
		files := testhelpers.NewTestFiles(t)
		tmpFile := files.CreateYAMLFile("compliance-req", "this is not valid yaml: {[")

		_, err := LoadCompScope(tmpFile, validator)
		if err == nil {
			t.Error("Expected error for invalid YAML but got nil")
		}
	})
}
