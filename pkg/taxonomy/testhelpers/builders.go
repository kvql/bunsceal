package testhelpers

import (
	"sort"

	"github.com/kvql/bunsceal/pkg/domain"
)

// Taxonomy Builders

// NewTestTaxonomy creates a minimal valid taxonomy for testing
func NewTestTaxonomy() *domain.Taxonomy {
	return &domain.Taxonomy{
		ApiVersion:        "v1beta1",
		SegL1s:            make(map[string]domain.Seg),
		Segs:              make(map[string]domain.Seg),
		CompReqs:          make(map[string]domain.CompReq),
		SensitivityLevels: []string{"A", "B", "C", "D"},
		CriticalityLevels: []string{"1", "2", "3", "4", "5"},
	}
}

// NewCompleteTaxonomy creates a complete valid taxonomy with all standard components
func NewCompleteTaxonomy() *domain.Taxonomy {
	txy := NewTestTaxonomy()
	txy.SegL1s["shared-service"] = NewSharedServiceSegL1()
	txy.SegL1s["prod"] = NewProdSegL1()
	txy.CompReqs = NewStandardCompReqs()
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
	if txy.Segs == nil {
		txy.Segs = make(map[string]domain.Seg)
	}
	txy.Segs[id] = seg
	return txy
}

// WithCompReq adds a compliance requirement to the taxonomy
func WithCompReq(txy *domain.Taxonomy, id string, req domain.CompReq) *domain.Taxonomy {
	if txy.CompReqs == nil {
		txy.CompReqs = make(map[string]domain.CompReq)
	}
	txy.CompReqs[id] = req
	return txy
}

// WithStandardCompReqs adds standard compliance requirements (pci-dss, sox)
func WithStandardCompReqs(txy *domain.Taxonomy) *domain.Taxonomy {
	for id, req := range NewStandardCompReqs() {
		txy = WithCompReq(txy, id, req)
	}
	return txy
}

// SegL1 Builders
// --------------

// NewProdSegL1 creates a standard production SegL1
func NewProdSegL1() domain.Seg {
	return domain.Seg{
		ID:                   "prod",
		Name:                 "Production",
		Description:          "Production environment with strict security controls for customer-facing services and data.",
		Sensitivity:          "A",
		SensitivityRationale: "Production handles customer data requiring highest classification level and protection.",
		Criticality:          "1",
		CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate response.",
		ComplianceReqs:       []string{"pci-dss", "sox"},
	}
}

// NewStagingSegL1 creates a standard staging SegL1
func NewStagingSegL1() domain.Seg {
	return domain.Seg{
		ID:                   "staging",
		Name:                 "Staging",
		Description:          "Pre-production staging environment for final testing and validation before deployment cycles.",
		Sensitivity:          "D",
		SensitivityRationale: "Staging contains no production or customer data, only synthetic test data generated for validation.",
		Criticality:          "5",
		CriticalityRationale: "Staging downtime impacts development velocity but has no direct customer or revenue impact.",
		ComplianceReqs:       []string{},
	}
}

// NewSharedServiceSegL1 creates a standard shared-service SegL1
func NewSharedServiceSegL1() domain.Seg {
	return domain.Seg{
		ID:                   "shared-service",
		Name:                 "Shared Service",
		Description:          "Shared service environment hosting cross-account resources and centralised services with connectivity.",
		Sensitivity:          "A",
		SensitivityRationale: "Shared services represent highest risk from lateral movement perspective and bridge between environments.",
		Criticality:          "1",
		CriticalityRationale: "All environments depend on shared services for core functionality making outages highly impactful.",
		ComplianceReqs:       []string{"pci-dss", "sox"},
	}
}

// NewSegL1 creates a SegL1 with the given parameters
func NewSegL1(id, name, sensitivity, criticality string, compReqs []string) domain.Seg {
	return domain.Seg{
		ID:                   id,
		Name:                 name,
		Description:          "This is a valid description with sufficient length to meet minimum requirements for validation purposes.",
		Sensitivity:          sensitivity,
		SensitivityRationale: "Valid sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
		Criticality:          criticality,
		CriticalityRationale: "Valid criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
		ComplianceReqs:       compReqs,
	}
}

// Seg Builders
// --------------

// NewAppSeg creates a standard application Seg
func NewAppSeg() domain.Seg {
	return domain.Seg{
		Name:        "Application",
		ID:          "app",
		Description: "Application domain for core business services",
		L1Parents:   []string{"prod"},
		L1Overrides: map[string]domain.L1Overrides{
			"prod": {
				Sensitivity:          "A",
				SensitivityRationale: "Applications handle customer PII and payment information requiring highest protection level.",
				Criticality:          "1",
				CriticalityRationale: "Application services are customer-facing and directly generate revenue requiring maximum uptime.",
				ComplianceReqs:       []string{"pci-dss"},
			},
		},
	}
}

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
		Description: "This is a valid description with sufficient length to meet minimum requirements for validation purposes.",
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
		Description: "This is a valid description with sufficient length to meet minimum requirements for validation purposes.",
		L1Parents:   l1Parents,
		L1Overrides: overrides,
	}
}

// NewL1Override creates a L1Override with the given parameters
func NewL1Override(sensitivity, criticality string, compReqs []string) domain.L1Overrides {
	return domain.L1Overrides{
		Sensitivity:          sensitivity,
		SensitivityRationale: "Valid sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
		Criticality:          criticality,
		CriticalityRationale: "Valid criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
		ComplianceReqs:       compReqs,
	}
}

// CompReq Builders
// ----------------

// NewStandardCompReqs returns the standard compliance requirements (pci-dss, sox)
func NewStandardCompReqs() map[string]domain.CompReq {
	return map[string]domain.CompReq{
		"pci-dss": {
			Name:        "PCI DSS",
			Description: "Payment Card Industry Data Security Standard",
			ReqsLink:    "https://www.pcisecuritystandards.org/",
		},
		"sox": {
			Name:        "SOX",
			Description: "Sarbanes-Oxley Act",
			ReqsLink:    "https://www.sox-online.com/",
		},
	}
}

// NewCompReq creates a compliance requirement
func NewCompReq(name, description, link string) domain.CompReq {
	return domain.CompReq{
		Name:        name,
		Description: description,
		ReqsLink:    link,
	}
}
