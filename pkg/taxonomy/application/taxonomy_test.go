package application

import (
	"os"
	"testing"

	"github.com/kvql/bunsceal/pkg/config"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

func TestApplyInheritance(t *testing.T) {
	t.Run("Inherits compliance requirements when nil", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox"},
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							// nil ComplianceReqs - should inherit
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		ApplyInheritance(&txy, nil)

		L1Overrides := txy.SegsL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.ComplianceReqs) != 2 {
			t.Errorf("Expected 2 inherited compliance reqs, got %d", len(L1Overrides.ComplianceReqs))
		}
		if L1Overrides.ComplianceReqs[0] != "pci-dss" || L1Overrides.ComplianceReqs[1] != "sox" {
			t.Errorf("Expected inherited compliance reqs [pci-dss, sox], got %v", L1Overrides.ComplianceReqs)
		}
	})

	t.Run("Does not inherit compliance requirements when set", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox", "hipaa"},
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs: []string{"pci-dss"}, // Custom subset
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
				"hipaa":   {Name: "HIPAA", Description: "Health Insurance Portability and Accountability Act", ReqsLink: "https://www.hhs.gov/hipaa/"},
			},
		}

		ApplyInheritance(&txy, nil)

		L1Overrides := txy.SegsL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.ComplianceReqs) != 1 {
			t.Errorf("Expected custom compliance reqs to remain (1 item), got %d", len(L1Overrides.ComplianceReqs))
		}
		if L1Overrides.ComplianceReqs[0] != "pci-dss" {
			t.Errorf("Expected compliance req [pci-dss], got %v", L1Overrides.ComplianceReqs)
		}
	})

	t.Run("Populates CompReqs map with full details", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID: "prod",
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs: []string{"pci-dss", "sox"},
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		ApplyInheritance(&txy, nil)

		L1Overrides := txy.SegsL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.CompReqs) != 2 {
			t.Errorf("Expected 2 entries in CompReqs map, got %d", len(L1Overrides.CompReqs))
		}
		if compReq, ok := L1Overrides.CompReqs["pci-dss"]; !ok {
			t.Error("Expected pci-dss in CompReqs map")
		} else {
			if compReq.Name != "PCI DSS" {
				t.Errorf("Expected PCI DSS name, got %s", compReq.Name)
			}
		}
		if compReq, ok := L1Overrides.CompReqs["sox"]; !ok {
			t.Error("Expected sox in CompReqs map")
		} else {
			if compReq.Name != "SOX" {
				t.Errorf("Expected SOX name, got %s", compReq.Name)
			}
		}
	})

	t.Run("Skips invalid compliance requirements in CompReqs map", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID: "prod",
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs: []string{"pci-dss", "invalid-scope"},
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		ApplyInheritance(&txy, nil)

		L1Overrides := txy.SegsL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.CompReqs) != 1 {
			t.Errorf("Expected only 1 valid entry in CompReqs map, got %d", len(L1Overrides.CompReqs))
		}
		if _, ok := L1Overrides.CompReqs["invalid-scope"]; ok {
			t.Error("Expected invalid-scope to be skipped in CompReqs map")
		}
	})

	t.Run("Handles L2 with multiple L1 parents", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss"},
				},
				"staging": {
					ID:             "staging",
					ComplianceReqs: []string{"sox"},
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod", "staging"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod":    {},
						"staging": {},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS"},
				"sox":     {Name: "SOX"},
			},
		}

		ApplyInheritance(&txy, nil)

		prodOverride := txy.SegsL2s["app"].L1Overrides["prod"]
		if len(prodOverride.ComplianceReqs) != 1 || prodOverride.ComplianceReqs[0] != "pci-dss" {
			t.Errorf("Expected prod to inherit [pci-dss], got %v", prodOverride.ComplianceReqs)
		}

		stagingOverride := txy.SegsL2s["app"].L1Overrides["staging"]
		if len(stagingOverride.ComplianceReqs) != 1 || stagingOverride.ComplianceReqs[0] != "sox" {
			t.Errorf("Expected staging to inherit [sox], got %v", stagingOverride.ComplianceReqs)
		}
	})
}

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
