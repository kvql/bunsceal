package validation

import (
	"testing"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
)

const testNs = "bunsceal.plugin.classifications"

// newSegWithClassification creates a Seg with classification labels
func newSegWithClassification(id, name, sensitivity, criticality string, compReqs []string) domain.Seg {
	labels := []string{
		testNs + "/sensitivity:" + sensitivity,
		testNs + "/criticality:" + criticality,
	}
	seg := domain.Seg{
		ID:             id,
		Name:           name,
		Labels:         labels,
		ComplianceReqs: compReqs,
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

		ruleSet := NewLogicRuleSet(config)

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

		ruleSet := NewLogicRuleSet(config)

		if len(ruleSet.LogicRules) != 1 {
			t.Errorf("Expected 1 rule, got %d", len(ruleSet.LogicRules))
		}

		if _, ok := ruleSet.LogicRules["Uniqueness"]; !ok {
			t.Error("Expected Uniqueness rule to be present")
		}
	})

	t.Run("Uses default configuration", func(t *testing.T) {
		config := configdomain.DefaultConfig()

		ruleSet := NewLogicRuleSet(config)

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
				"shared-service": newSegWithClassification("shared-service", "Shared Service", "A", "1", []string{"req1"}),
				"prod":           newSegWithClassification("prod", "Production", "B", "2", []string{"req1"}),
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					ID:   "app",
					Name: "Application",
				},
			},
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
			},
		}

		ruleSet := NewLogicRuleSet(config)
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
			},
		}

		ruleSet := NewLogicRuleSet(config)
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
			CompReqs: map[string]domain.CompReq{},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
			},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
			},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
				"req2": {Name: "Requirement 2"},
			},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
				"req2": {Name: "Requirement 2"},
			},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			CompReqs: map[string]domain.CompReq{
				"req1": {Name: "Requirement 1"},
				"req2": {Name: "Requirement 2"},
			},
		}

		rule := NewLogicRuleSharedService(configdomain.GeneralBooleanConfig{Enabled: true})
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
			SegsL2s:  map[string]domain.Seg{},
			CompReqs: map[string]domain.CompReq{},
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
			CompReqs: map[string]domain.CompReq{},
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
			CompReqs: map[string]domain.CompReq{},
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
			CompReqs: map[string]domain.CompReq{},
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
