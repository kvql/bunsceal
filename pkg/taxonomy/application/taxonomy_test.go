package application

import (
	"os"
	"testing"

	"github.com/kvql/bunsceal/pkg/config"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

func TestValidatePluginLabels(t *testing.T) {
	t.Run("Returns nil when no plugins configured", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s:  map[string]domain.Seg{},
			SegsL2s: map[string]domain.Seg{},
		}

		err := ValidatePluginLabels(&txy, nil)

		if err != nil {
			t.Errorf("Expected nil error for nil plugins, got %v", err)
		}
	})

	t.Run("Fails when segment has invalid plugin label value", func(t *testing.T) {
		config := &plugins.ClassificationsConfig{
			Common:          plugins.PluginsCommonSettings{LabelInheritance: true},
			RationaleLength: 10,
			Definitions: map[string]plugins.ClassificationDefinition{
				"sensitivity": {
					DescriptiveName: "Sensitivity",
					Values:          map[string]string{"high": "High", "low": "Low"},
				},
			},
		}
		pluginsList := make(plugins.Plugins)
		pluginsList["classifications"] = plugins.NewClassificationPlugin(config, plugins.NsPrefix)

		seg := domain.Seg{
			ID:   "prod",
			Name: "Production",
			Labels: []string{
				"bunsceal.plugin.classifications/sensitivity:invalid-value",
				"bunsceal.plugin.classifications/sensitivity_rationale:Test rationale",
			},
		}
		seg.ParseLabels()

		txy := domain.Taxonomy{
			SegL1s:  map[string]domain.Seg{"prod": seg},
			SegsL2s: map[string]domain.Seg{},
		}

		err := ValidatePluginLabels(&txy, pluginsList)

		if err == nil {
			t.Error("Expected validation error for invalid classification value")
		}
	})

	t.Run("Passes when segment has valid plugin labels", func(t *testing.T) {
		config := &plugins.ClassificationsConfig{
			Common:          plugins.PluginsCommonSettings{LabelInheritance: true},
			RationaleLength: 10,
			Definitions: map[string]plugins.ClassificationDefinition{
				"sensitivity": {
					DescriptiveName: "Sensitivity",
					Values:          map[string]string{"high": "High", "low": "Low"},
				},
			},
		}
		pluginsList := make(plugins.Plugins)
		pluginsList["classifications"] = plugins.NewClassificationPlugin(config, plugins.NsPrefix)

		seg := domain.Seg{
			ID:   "prod",
			Name: "Production",
			Labels: []string{
				"bunsceal.plugin.classifications/sensitivity:high",
				"bunsceal.plugin.classifications/sensitivity_rationale:Valid rationale here",
			},
		}
		seg.ParseLabels()

		txy := domain.Taxonomy{
			SegL1s:  map[string]domain.Seg{"prod": seg},
			SegsL2s: map[string]domain.Seg{},
		}

		err := ValidatePluginLabels(&txy, pluginsList)

		if err != nil {
			t.Errorf("Expected no error for valid labels, got %v", err)
		}
	})
}

func TestLoadExampleTaxonomy(t *testing.T) {
	t.Run("Example taxonomy files are valid and load successfully", func(t *testing.T) {
		// Change to project root so relative paths in config work
		originalWd, err := os.Getwd()
		if err != nil {
			t.Fatalf("Failed to get working directory: %v", err)
		}
		if err := os.Chdir("../../.."); err != nil {
			t.Fatalf("Failed to change to project root: %v", err)
		}
		defer os.Chdir(originalWd)

		cfg, err := config.LoadConfig("example/config.yaml", "pkg/config/schemas")
		if err != nil {
			t.Fatalf("Failed to load example config: %v", err)
		}

		_, err = LoadTaxonomy(cfg)
		if err != nil {
			t.Errorf("Example taxonomy should load without error: %v", err)
		}
	})
}
