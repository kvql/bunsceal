package plugins

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
)

const (
	ScopeInScope    = "in-scope"
	ScopeOutOfScope = "out-of-scope"
)

type ComplianceConfig struct {
	Common                PluginsCommonSettings           `yaml:"common_settings"`
	RationaleLength       int                             `yaml:"rationale_length"`
	EnforceScopeHierarchy bool                            `yaml:"enforce_scope_hierarchy"`
	Definitions           map[string]ComplianceDefinition `yaml:"definitions"`
}

type ComplianceDefinition struct {
	DescriptiveName  string `yaml:"name"`
	Description      string `yaml:"description"`
	RequirementsLink string `yaml:"requirements_link,omitempty"`
}

type CompliancePlugin struct {
	Config    *ComplianceConfig
	Namespace string
}

func NewCompliancePlugin(config *ComplianceConfig, prefix string) *CompliancePlugin {
	return &CompliancePlugin{
		Config:    config,
		Namespace: prefix + "compliance",
	}
}

func (p CompliancePlugin) validateNamespaceLabels(labels map[string]string, ctx string, errs *[]error) int {
	foundKeys := 0
	for reqID := range p.Config.Definitions {
		scope, hasScope := labels[reqID]
		rationale, hasRat := labels[reqID+"_rationale"]

		if hasScope && !hasRat {
			foundKeys++
			*errs = append(*errs, fmt.Errorf("%s has %s but missing %s_rationale", ctx, reqID, reqID))
		}
		if hasRat && !hasScope {
			foundKeys++
			*errs = append(*errs, fmt.Errorf("%s has %s_rationale but missing %s", ctx, reqID, reqID))
		}
		if hasScope && hasRat {
			foundKeys++
			foundKeys++
			// Validate scope value
			if scope != ScopeInScope && scope != ScopeOutOfScope {
				*errs = append(*errs, fmt.Errorf("%s invalid scope value %s for %s (must be '%s' or '%s')", ctx, scope, reqID, ScopeInScope, ScopeOutOfScope))
			}
			// Validate rationale length
			if len(rationale) < p.Config.RationaleLength {
				*errs = append(*errs, fmt.Errorf("%s %s_rationale too short (min %d chars)", ctx, reqID, p.Config.RationaleLength))
			}
		}
	}
	return foundKeys
}

func (p CompliancePlugin) ValidateLabels(seg *domain.Seg) PluginValidationResult {
	result := PluginValidationResult{Valid: false, Errors: []error{}}

	segLabels := seg.LabelNamespaces[p.Namespace]
	p.validateNamespaceLabels(segLabels, "segment "+seg.ID, &result.Errors)

	// L1 segments are not required to have complete compliance labels (unlike classification)
	// Compliance is opt-in per requirement

	// Validate override labels
	for parentID, override := range seg.L1Overrides {
		if len(override.LabelNamespaces[p.Namespace]) > 0 {
			p.validateNamespaceLabels(override.LabelNamespaces[p.Namespace], fmt.Sprintf("segment %s l1_override[%s]", seg.ID, parentID), &result.Errors)
		}
	}

	result.Valid = len(result.Errors) == 0
	return result
}

// ValidateRelationship checks parent-child compliance scope hierarchy when EnforceScopeHierarchy is enabled.
// A child cannot be in-scope for a requirement that the parent doesn't have defined.
func (p CompliancePlugin) ValidateRelationship(parent, child *domain.Seg) []error {
	if !p.Config.EnforceScopeHierarchy {
		return nil
	}

	var errs []error
	override, hasOverride := child.L1Overrides[parent.ID]

	for reqID := range p.Config.Definitions {
		parentScope := parent.LabelNamespaces[p.Namespace][reqID]
		childScope := child.LabelNamespaces[p.Namespace][reqID]
		if hasOverride && len(override.LabelNamespaces[p.Namespace]) > 0 {
			if overrideScope, ok := override.LabelNamespaces[p.Namespace][reqID]; ok {
				childScope = overrideScope
			}
		}

		// If parent doesn't have the requirement defined and child is in-scope, error
		if parentScope == "" && childScope == ScopeInScope {
			errs = append(errs, fmt.Errorf("child %s is in-scope for %s but parent %s doesn't have this requirement defined",
				child.ID, reqID, parent.ID))
		}
	}
	return errs
}

func (p CompliancePlugin) GetEnabled() bool {
	return p.Config.Common.LabelInheritance
}

func (p CompliancePlugin) GetNamespace() string {
	return p.Namespace
}

func (p CompliancePlugin) GetImageData() []ImageGroupingData {
	dataList := []ImageGroupingData{}
	for key, def := range p.Config.Definitions {
		data := ImageGroupingData{
			Namespace:     p.Namespace,
			DisplayName:   def.DescriptiveName,
			OrderedValues: []string{ScopeInScope, ScopeOutOfScope},
			OrderMap:      map[string]int{ScopeInScope: 0, ScopeOutOfScope: 1},
			Key:           key,
			ValuesMap:     map[string]string{ScopeInScope: "In Scope", ScopeOutOfScope: "Out of Scope"},
		}
		dataList = append(dataList, data)
	}
	return dataList
}
