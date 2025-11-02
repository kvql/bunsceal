package validation

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
)

func TestValidateL1Definitions(t *testing.T) {
	t.Run("Valid compliance requirements pass", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Description:          "Production environment with strict security controls for customer-facing services and data.",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data requiring highest classification level and protection.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate response.",
					ComplianceReqs:       []string{"pci-dss", "sox"},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		valid := ValidateL1Definitions(&txy)
		if !valid {
			t.Error("Expected validation to pass for valid compliance requirements")
		}
	})

	t.Run("Invalid compliance requirement fails", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Description:          "Production environment with strict security controls for customer-facing services and data.",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data requiring highest classification level and protection.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate response.",
					ComplianceReqs:       []string{"invalid-scope"},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid := ValidateL1Definitions(&txy)
		if valid {
			t.Error("Expected validation to fail for invalid compliance requirement")
		}
	})

	t.Run("Empty compliance requirements pass", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"staging": {
					ID:                   "staging",
					Name:                 "Staging",
					Description:          "Pre-production staging environment for final testing and validation before deployment cycles.",
					Sensitivity:          "D",
					SensitivityRationale: "Staging contains no production or customer data, only synthetic test data generated for validation.",
					Criticality:          "5",
					CriticalityRationale: "Staging downtime impacts development velocity but has no direct customer or revenue impact.",
					ComplianceReqs:       []string{},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		valid := ValidateL1Definitions(&txy)
		if !valid {
			t.Error("Expected validation to pass for empty compliance requirements")
		}
	})
}

func TestValidateL2Definition(t *testing.T) {
	t.Run("Valid security domains pass", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Description:          "Production environment with strict security controls for customer-facing services and data.",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data requiring highest classification level and protection.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate response.",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain for core business services",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							Sensitivity:          "A",
							SensitivityRationale: "Applications handle customer PII and payment information requiring highest protection level.",
							Criticality:          "1",
							CriticalityRationale: "Application services are customer-facing and directly generate revenue requiring maximum uptime.",
							ComplianceReqs:       []string{"pci-dss"},
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid, failures := ValidateL2Definition(&txy)
		if !valid {
			t.Errorf("Expected validation to pass, got %d failures", failures)
		}
		if failures != 0 {
			t.Errorf("Expected 0 failures, got %d", failures)
		}
	})

	t.Run("Invalid compliance requirement in SegL2 fails", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs: []string{"invalid-scope"},
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		valid, failures := ValidateL2Definition(&txy)
		if valid {
			t.Error("Expected validation to fail for invalid compliance requirement")
		}
		if failures == 0 {
			t.Error("Expected at least 1 failure")
		}
	})

	t.Run("Invalid environment ID in SegL2 fails", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"invalid-env": {
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
							Criticality:          "1",
							CriticalityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		valid, failures := ValidateL2Definition(&txy)
		if valid {
			t.Error("Expected validation to fail for invalid environment ID")
		}
		if failures == 0 {
			t.Error("Expected at least 1 failure")
		}
	})

	t.Run("Multiple validation failures counted", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs: []string{"invalid1", "invalid2"},
						},
						"invalid-env": {
							ComplianceReqs: []string{"invalid3"},
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		valid, failures := ValidateL2Definition(&txy)
		if valid {
			t.Error("Expected validation to fail")
		}
		// Should have failures for: 3 invalid comp reqs + 1 invalid env = 4 failures
		if failures < 3 {
			t.Errorf("Expected at least 3 failures, got %d", failures)
		}
	})
}

func TestValidateSharedServices(t *testing.T) {
	t.Run("Shared service environment not found", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: make(map[string]domain.SegL1),
		}
		valid, _ := ValidateSharedServices(txy)
		if valid {
			t.Error("Expected valid to be false")
		}
	})

	t.Run("Shared service environment with incorrect sensitivity, criticality and comp requirements", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"shared-service": {
					Sensitivity:    "b",
					Criticality:    "2",
					ComplianceReqs: []string{"req1", "req2"},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"req1": {},
				"req2": {},
				"req3": {},
			},
		}
		valid, failures := ValidateSharedServices(txy)
		if valid {
			t.Error("Expected valid to be false")
		}
		if failures != 2 {
			t.Errorf("Expected failures to be 2, got %d", failures)
		}
	})

	t.Run("Valid shared service environment", func(t *testing.T) {
		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"shared-service": {
					Sensitivity:    domain.SenseOrder[0],
					Criticality:    domain.CritOrder[0],
					ComplianceReqs: []string{"req1", "req2"},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"req1": {},
				"req2": {},
			},
		}
		valid, failures := ValidateSharedServices(txy)
		if !valid {
			t.Error("Expected valid to be true")
		}
		if failures != 0 {
			t.Errorf("Expected failures to be 0, got %d", failures)
		}
	})
}
