package validation

import (
	"fmt"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

// LogicRule defines the interface for business logic validation rules.
type LogicRule interface {
	Validate(taxonomy *domain.Taxonomy) []error
}

// ValidationResult captures the results of a single rule's validation.
type ValidationResult struct {
	RuleName string
	Errors   []error
}

// LogicRuleSet holds a collection of business logic validation rules.
type LogicRuleSet struct {
	LogicRules map[string]LogicRule
}

// NewLogicRuleSet creates a new LogicRuleSet based on the provided configuration.
// Only enabled rules are instantiated and added to the set.
func NewLogicRuleSet(config configdomain.Config, pluginMap plugins.Plugins) *LogicRuleSet {
	rs := &LogicRuleSet{LogicRules: make(map[string]LogicRule)}

	if config.Rules.SharedService.Enabled {
		classificationPlugin := pluginMap["classifications"]
		compliancePlugin := pluginMap["compliance"]
		rs.LogicRules["SharedService"] = NewLogicRuleSharedService(config.Rules.SharedService, classificationPlugin, compliancePlugin)
	}

	if config.Rules.Uniqueness.Enabled {
		rs.LogicRules["Uniqueness"] = NewLogicRuleUniqueness(config.Rules.Uniqueness)
	}

	return rs
}

// ValidateAll runs all configured rules against the taxonomy and aggregates the results.
// Returns a slice of ValidationResult, one for each rule that produced errors.
func (rs *LogicRuleSet) ValidateAll(taxonomy *domain.Taxonomy) []ValidationResult {
	var results []ValidationResult

	for name, logicRule := range rs.LogicRules {
		if errs := logicRule.Validate(taxonomy); len(errs) > 0 {
			results = append(results, ValidationResult{
				RuleName: name,
				Errors:   errs,
			})
		}
	}

	return results
}

// LogicRuleSharedService validates the shared-services environment.
// This rule ensures the shared-services environment meets the strictest requirements.
type LogicRuleSharedService struct {
	config            configdomain.GeneralBooleanConfig
	classificationPlugin plugins.Plugin
	compliancePlugin  plugins.Plugin
}

// NewLogicRuleSharedService creates a new SharedService validation rule.
func NewLogicRuleSharedService(config configdomain.GeneralBooleanConfig, classificationPlugin plugins.Plugin, compliancePlugin plugins.Plugin) *LogicRuleSharedService {
	return &LogicRuleSharedService{
		config:            config,
		classificationPlugin: classificationPlugin,
		compliancePlugin:  compliancePlugin,
	}
}

// Validate checks that the shared-service environment meets all requirements.
// Returns a slice of errors if validation fails, or an empty slice if valid.
func (r *LogicRuleSharedService) Validate(taxonomy *domain.Taxonomy) []error {
	var errs []error
	envName := "shared-service"
	ns := "bunsceal.plugin.classifications"

	if _, ok := taxonomy.SegL1s[envName]; !ok {
		err := fmt.Errorf("%s environment not found", envName)
		o11y.Log.Printf("%v", err)
		errs = append(errs, err)
		return errs
	}

	sharedSeg := taxonomy.SegL1s[envName]

	// Helper to read classification from plugin labels
	getSensitivity := func(seg domain.Seg) string {
		if val, exists := seg.LabelNamespaces[ns]["sensitivity"]; exists {
			return val
		}
		return ""
	}
	getCriticality := func(seg domain.Seg) string {
		if val, exists := seg.LabelNamespaces[ns]["criticality"]; exists {
			return val
		}
		return ""
	}

	// Get expected highest values from classifications plugin
	imageData := r.classificationPlugin.GetImageData()
	expectedSens := ""
	expectedCrit := ""
	for _, data := range imageData {
		if data.Key == "sensitivity" && len(data.OrderedValues) > 0 {
			expectedSens = data.OrderedValues[0]
		}
		if data.Key == "criticality" && len(data.OrderedValues) > 0 {
			expectedCrit = data.OrderedValues[0]
		}
	}

	if getSensitivity(sharedSeg) != expectedSens || getCriticality(sharedSeg) != expectedCrit {
		err := fmt.Errorf("%s environment does not have the highest sensitivity or criticality", envName)
		o11y.Log.Printf("%v", err)
		errs = append(errs, err)
	}

	// Check if shared-service has all compliance requirements defined via labels
	if r.compliancePlugin != nil {
		if cp, ok := r.compliancePlugin.(*plugins.CompliancePlugin); ok {
			ns := cp.GetNamespace()
			totalCompReqs := len(cp.Config.Definitions)
			inScopeCount := 0

			// Count how many compliance requirements are in-scope
			for reqID := range cp.Config.Definitions {
				if scope, exists := sharedSeg.LabelNamespaces[ns][reqID]; exists && scope == plugins.ScopeInScope {
					inScopeCount++
				}
			}

			if inScopeCount != totalCompReqs {
				err := fmt.Errorf("%s environment does not have all compliance requirements in-scope", envName)
				o11y.Log.Printf("%v", err)
				errs = append(errs, err)
			}
		}
	}

	return errs
}

// LogicRuleUniqueness validates that specified fields are unique across taxonomy entities.
type LogicRuleUniqueness struct {
	config configdomain.UniquenessConfig
}

// NewLogicRuleUniqueness creates a new Uniqueness validation rule.
func NewLogicRuleUniqueness(config configdomain.UniquenessConfig) *LogicRuleUniqueness {
	return &LogicRuleUniqueness{config: config}
}

// Validate checks that the configured keys are unique across L1 and L2 segments.
// Returns a slice of errors if validation fails, or an empty slice if valid.
func (r *LogicRuleUniqueness) Validate(taxonomy *domain.Taxonomy) []error {
	var errs []error

	// Track which keys to check
	for _, key := range r.config.CheckKeys {
		// Check L1 names
		l1Values := make(map[string]bool)
		for _, seg := range taxonomy.SegL1s {
			val, err := seg.GetKeyString(key)
			if err != nil {
				return []error{err}
			}
			if l1Values[val] {
				err := fmt.Errorf("duplicate L1 name found: %s", val)
				o11y.Log.Printf("%v", err)
				errs = append(errs, err)
			}
			l1Values[val] = true
		}

		// Check L2 names
		l2Values := make(map[string]bool)
		for _, seg := range taxonomy.SegsL2s {
			val, err := seg.GetKeyString(key)
			if err != nil {
				return []error{err}
			}
			if l2Values[val] {
				err := fmt.Errorf("duplicate L2 name found: %s", val)
				o11y.Log.Printf("%v", err)
				errs = append(errs, err)
			}
			l2Values[val] = true
		}
	}

	return errs
}
