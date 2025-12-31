package validation

import (
	"testing"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

const testNs = "bunsceal.plugin.classifications"

// createMockClassificationPlugin creates a mock classification plugin for testing
func createMockClassificationPlugin() plugins.Plugin {
	config := &plugins.ClassificationsConfig{
		Definitions: map[string]plugins.ClassificationDefinition{
			"sensitivity": {
				DescriptiveName: "Sensitivity",
				Order:           []string{"A", "B", "C", "D"},
				Values: map[string]string{
					"A": "High",
					"B": "Medium",
					"C": "Low",
					"D": "N/A",
				},
			},
			"criticality": {
				DescriptiveName: "Criticality",
				Order:           []string{"1", "2", "3", "4", "5"},
				Values: map[string]string{
					"1": "Critical",
					"2": "High",
					"3": "Medium",
					"4": "Low",
					"5": "N/A",
				},
			},
		},
	}
	return plugins.NewClassificationPlugin(config, "bunsceal.plugin.")
}

// createMockCompliancePluginForLogicRules creates a mock compliance plugin for testing
func createMockCompliancePluginForLogicRules() plugins.Plugin {
	config := &plugins.ComplianceConfig{
		Definitions: map[string]plugins.ComplianceDefinition{
			"req1": {
				DescriptiveName:  "Requirement 1",
				Description:      "Test requirement 1",
				RequirementsLink: "https://example.com/req1",
			},
			"req2": {
				DescriptiveName:  "Requirement 2",
				Description:      "Test requirement 2",
				RequirementsLink: "https://example.com/req2",
			},
		},
	}
	return plugins.NewCompliancePlugin(config, "bunsceal.plugin.")
}

// createMockPluginMap creates a plugin map with mock classification and compliance plugins
func createMockPluginMap() plugins.Plugins {
	pluginMap := make(plugins.Plugins)
	pluginMap["classifications"] = createMockClassificationPlugin()
	pluginMap["compliance"] = createMockCompliancePluginForLogicRules()
	return pluginMap
}

// complianceLabel creates compliance labels for a requirement ID
func complianceLabel(reqID string) []string {
	compNs := "bunsceal.plugin.compliance"
	return []string{
		compNs + "/" + reqID + ":" + plugins.ScopeInScope,
		compNs + "/" + reqID + "_rationale:Test rationale for " + reqID,
	}
}

// newSegWithClassification creates a Seg with classification and compliance labels
func newSegWithClassification(id, name, sensitivity, criticality string, compReqIDs []string) domain.Seg {
	labels := []string{
		testNs + "/sensitivity:" + sensitivity,
		testNs + "/criticality:" + criticality,
	}

	// Add compliance labels for each requirement ID
	for _, reqID := range compReqIDs {
		labels = append(labels, complianceLabel(reqID)...)
	}

	seg := domain.Seg{
		ID:     id,
		Name:   name,
		Labels: labels,
	}
	seg.ParseLabels()
	return seg
}

func TestNewLogicRuleSet(t *testing.T) {
	t.Run("Creates empty ruleset when all rules disabled", func(t *testing.T) {
		config := configdomain.Config{
			Rules: configdomain.LogicRulesConfig{
				SharedService: configdomain.GeneralBooleanConfig{Enabled: false},
				Uniqueness:    configdomain.UniquenessConfig{Enabled: false},
			},
		}

		pluginMap := createMockPluginMap()
		ruleSet := NewLogicRuleSet(config, pluginMap)

		if len(ruleSet.LogicRules) != 0 {
			t.Errorf("Expected 0 rules, got %d", len(ruleSet.LogicRules))
		}
	})

	t.Run("Creates ruleset with Uniqueness rule enabled", func(t *testing.T) {
		config := configdomain.Config{
			Rules: configdomain.LogicRulesConfig{
				SharedService: configdomain.GeneralBooleanConfig{Enabled: false},
				Uniqueness:    configdomain.UniquenessConfig{Enabled: true, CheckKeys: []string{"name"}},
			},
		}

		pluginMap := createMockPluginMap()
		ruleSet := NewLogicRuleSet(config, pluginMap)

		if len(ruleSet.LogicRules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(ruleSet.LogicRules))
		}

		if _, ok := ruleSet.LogicRules["Uniqueness"]; !ok {
			t.Error("Expected Uniqueness rule to be present")
		}
	})

	t.Run("Uses default configuration", func(t *testing.T) {
		config := configdomain.DefaultConfig()

		pluginMap := createMockPluginMap()
		ruleSet := NewLogicRuleSet(config, pluginMap)

		if len(ruleSet.LogicRules) != 2 {
			t.Errorf("Expected 2 rules by default, got %d", len(ruleSet.LogicRules))
		}

		if _, ok := ruleSet.LogicRules["SharedService"]; !ok {
			t.Error("Expected SharedService rule to be enabled by default")
		}

		if _, ok := ruleSet.LogicRules["Uniqueness"]; !ok {
			t.Error("Expected Uniqueness rule to be enabled by default")
		}
	})
}

func TestLogicRuleSet_ValidateAll(t *testing.T) {
	t.Run("Returns empty results when all validations pass", func(t *testing.T) {
		config := configdomain.DefaultConfig()
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "A", "1", []string{"req1", "req2"}),
				"prod":           newSegWithClassification("prod", "Production", "B", "2", []string{"req1"}),
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					ID:   "app",
					Name: "Application",
				},
			},
		}

		pluginMap := createMockPluginMap()
		ruleSet := NewLogicRuleSet(config, pluginMap)
		results := ruleSet.ValidateAll(txy)

		if len(results) != 0 {
			t.Errorf("Expected no validation errors, got %d results", len(results))
			for _, result := range results {
				t.Logf("Rule '%s' failed with %d errors", result.RuleName, len(result.Errors))
				for _, err := range result.Errors {
					t.Logf("  - %v", err)
				}
			}
		}
	})

	t.Run("Returns results for failing rules", func(t *testing.T) {
		config := configdomain.DefaultConfig()
		// Missing shared-service environment
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID:   "prod",
					Name: "Production",
				},
			},
			SegsL2s: map[string]domain.Seg{},
		}

		pluginMap := createMockPluginMap()
		ruleSet := NewLogicRuleSet(config, pluginMap)
		results := ruleSet.ValidateAll(txy)

		if len(results) == 0 {
			t.Error("Expected validation errors, got none")
		}

		foundSharedServiceError := false
		for _, result := range results {
			if result.RuleName == "SharedService" {
				foundSharedServiceError = true
				if len(result.Errors) == 0 {
					t.Error("Expected SharedService rule to have errors")
				}
			}
		}

		if !foundSharedServiceError {
			t.Error("Expected SharedService rule to fail")
		}
	})
}

func TestLogicRuleSharedService_Validate(t *testing.T) {
	t.Run("Fails when shared-service environment not found", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {ID: "prod", Name: "Production"},
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors when shared-service environment not found")
		}
	})

	t.Run("Fails when shared-service has incorrect sensitivity", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "B", "1", []string{"req1"}),
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors when shared-service has incorrect sensitivity")
		}
	})

	t.Run("Fails when shared-service has incorrect criticality", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "A", "2", []string{"req1"}),
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors when shared-service has incorrect criticality")
		}
	})

	t.Run("Fails when shared-service missing compliance requirements", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "A", "1", []string{"req1"}),
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors when shared-service missing compliance requirements")
		}
	})

	t.Run("Passes when shared-service is correctly configured", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "A", "1", []string{"req1", "req2"}),
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) != 0 {
			t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
		}
	})

	t.Run("Reports multiple errors for multiple issues", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "B", "2", []string{"req1"}),
			},
		}

		mockPluginMap := createMockPluginMap()
		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true}, mockPluginMap["classifications"], mockPluginMap["compliance"])
		errs := rule.Validate(txy)

		if len(errs) != 2 {
			t.Errorf("Expected 2 errors, got %d", len(errs))
		}
	})
}

func TestLogicRuleUniqueness_Validate(t *testing.T) {
	t.Run("Fails when duplicate L1 names found", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod":    {ID: "prod", Name: "Production"},
				"staging": {ID: "staging", Name: "Production"}, // Duplicate name
			},
			SegsL2s: map[string]domain.Seg{},
		}

		rule := NewLogicRuleUniqueness(configdomain.UniquenessConfig{Enabled: true, CheckKeys: []string{"name"}})
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors for duplicate L1 names")
		}
	})

	t.Run("Fails when duplicate L2 names found", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{},
			SegsL2s: map[string]domain.Seg{
				"app1": {ID: "app1", Name: "Application"},
				"app2": {ID: "app2", Name: "Application"}, // Duplicate name
			},
		}

		rule := NewLogicRuleUniqueness(configdomain.UniquenessConfig{Enabled: true, CheckKeys: []string{"name"}})
		errs := rule.Validate(txy)

		if len(errs) == 0 {
			t.Error("Expected errors for duplicate L2 names")
		}
	})

	t.Run("Passes when all names unique", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod":    {ID: "prod", Name: "Production"},
				"staging": {ID: "staging", Name: "Staging"},
			},
			SegsL2s: map[string]domain.Seg{
				"app1": {ID: "app1", Name: "Application"},
				"app2": {ID: "app2", Name: "Database"},
			},
		}

		rule := NewLogicRuleUniqueness(configdomain.UniquenessConfig{Enabled: true, CheckKeys: []string{"name"}})
		errs := rule.Validate(txy)

		if len(errs) != 0 {
			t.Errorf("Expected no errors, got %d: %v", len(errs), errs)
		}
	})

	t.Run("Allows same name in L1 and L2", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {ID: "app", Name: "Production"}, // Same name as L1 is OK
			},
		}

		rule := NewLogicRuleUniqueness(configdomain.UniquenessConfig{Enabled: true, CheckKeys: []string{"name"}})
		errs := rule.Validate(txy)

		if len(errs) != 0 {
			t.Errorf("Expected no errors (L1 and L2 can share names), got %d: %v", len(errs), errs)
		}
	})
}

func TestDefaultConfig_Rules(t *testing.T) {
	t.Run("Default config has SharedService rule enabled", func(t *testing.T) {
		config := configdomain.DefaultConfig()

		if !config.Rules.SharedService.Enabled {
			t.Error("Expected SharedService rule to be enabled by default")
		}
	})

	t.Run("Default config has Uniqueness rule enabled", func(t *testing.T) {
		config := configdomain.DefaultConfig()

		if !config.Rules.Uniqueness.Enabled {
			t.Error("Expected Uniqueness rule to be enabled by default")
		}
	})

	t.Run("Default config checks 'name' key for uniqueness", func(t *testing.T) {
		config := configdomain.DefaultConfig()

		if len(config.Rules.Uniqueness.CheckKeys) != 1 {
			t.Errorf("Expected 1 check key, got %d", len(config.Rules.Uniqueness.CheckKeys))
		}

		if config.Rules.Uniqueness.CheckKeys[0] != "name" {
			t.Errorf("Expected check key 'name', got '%s'", config.Rules.Uniqueness.CheckKeys[0])
		}
	})
}
