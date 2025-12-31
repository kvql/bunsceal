package testhelpers

import (
	"sort"

	"github.com/kvql/bunsceal/pkg/domain"
)

// Constants for minimal valid test data
const (
	ValidDescription = "Valid description meeting minimum length requirements for schema validation"
	ValidRationale   = "Valid rationale meeting minimum character requirement for descriptions"
)

// Taxonomy Builders

// NewTestTaxonomy creates a minimal valid taxonomy for testing
func NewTestTaxonomy() *domain.Taxonomy {
	return &domain.Taxonomy{
		ApiVersion: "v1beta1",
		SegL1s:     make(map[string]domain.Seg),
		SegsL2s:    make(map[string]domain.Seg),
	}
}

// NewCompleteTaxonomy creates a complete valid taxonomy with all standard components
func NewCompleteTaxonomy() *domain.Taxonomy {
	txy := NewTestTaxonomy()
	txy.SegL1s["shared-service"] = NewSegL1("shared-service", "Shared Service", "A", "1", []string{"pci-dss", "sox"})
	txy.SegL1s["prod"] = NewSegL1("prod", "Production", "A", "1", []string{"pci-dss", "sox"})
	return txy
}

// WithSegL1 adds a SegL1 to the taxonomy
func WithSegL1(txy *domain.Taxonomy, id string, seg domain.Seg) *domain.Taxonomy {
	if txy.SegL1s == nil {
		txy.SegL1s = make(map[string]domain.Seg)
	}
	txy.SegL1s[id] = seg
	return txy
}

// WithSeg adds a Seg to the taxonomy
func WithSeg(txy *domain.Taxonomy, id string, seg domain.Seg) *domain.Taxonomy {
	if txy.SegsL2s == nil {
		txy.SegsL2s = make(map[string]domain.Seg)
	}
	txy.SegsL2s[id] = seg
	return txy
}

// SegL1 Builders
// --------------

// NewSegL1 creates a SegL1 with the given parameters
func NewSegL1(id, name, sensitivity, criticality string, compLabels []string) domain.Seg {
	labels := []string{
		"bunsceal.plugin.classifications/sensitivity:" + sensitivity,
		"bunsceal.plugin.classifications/sensitivity_rationale:" + ValidRationale,
		"bunsceal.plugin.classifications/criticality:" + criticality,
		"bunsceal.plugin.classifications/criticality_rationale:" + ValidRationale,
	}
	labels = append(labels, compLabels...)

	seg := domain.Seg{
		ID:          id,
		Name:        name,
		Description: ValidDescription,
		Labels:      labels,
	}
	seg.ParseLabels()
	return seg
}

// Seg Builders
// --------------

// NewSeg creates a Seg with the given parameters
// L1Parents is auto-populated from overrides keys for backward compatibility
func NewSeg(id, name string, overrides map[string]domain.L1Overrides) domain.Seg {
	// Extract L1Parents from overrides keys for backward compatibility
	l1Parents := make([]string, 0, len(overrides))
	for l1ID := range overrides {
		l1Parents = append(l1Parents, l1ID)
	}
	sort.Strings(l1Parents)

	return domain.Seg{
		Name:        name,
		ID:          id,
		Description: ValidDescription,
		L1Parents:   l1Parents,
		L1Overrides: overrides,
	}
}

// NewSegWithParents creates a Seg with explicit parents and overrides
// Allows testing parent without override scenarios
func NewSegWithParents(id, name string, l1Parents []string, overrides map[string]domain.L1Overrides) domain.Seg {
	return domain.Seg{
		Name:        name,
		ID:          id,
		Description: ValidDescription,
		L1Parents:   l1Parents,
		L1Overrides: overrides,
	}
}

// NewL1Override creates a L1Override with the given parameters
func NewL1Override(sensitivity, criticality string, compLabels []string) domain.L1Overrides {
	labels := []string{
		"bunsceal.plugin.classifications/sensitivity:" + sensitivity,
		"bunsceal.plugin.classifications/sensitivity_rationale:" + ValidRationale,
		"bunsceal.plugin.classifications/criticality:" + criticality,
		"bunsceal.plugin.classifications/criticality_rationale:" + ValidRationale,
	}
	labels = append(labels, compLabels...)

	override := domain.L1Overrides{
		Labels: labels,
	}
	// Parse the labels into ParsedLabels and LabelNamespaces maps
	override.ParseLabels()
	return override
}
