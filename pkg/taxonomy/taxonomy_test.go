package taxonomy

import (
	"strings"
	"testing"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
)

func TestApplyInheritance(t *testing.T) {
	t.Run("Inherits sensitivity from SegL1 when empty", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data requiring highest classification level.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue.",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"infra": {
					Name:        "Infrastructure",
					ID:          "infra",
					Description: "Infrastructure domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							// Empty sensitivity and rationale - should inherit
							Criticality:          "2",
							CriticalityRationale: "Custom criticality for infrastructure.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["infra"].L1Overrides["prod"]
		if L1Overrides.Sensitivity != "A" {
			t.Errorf("Expected sensitivity 'A', got %s", L1Overrides.Sensitivity)
		}
		if !strings.HasPrefix(L1Overrides.SensitivityRationale, "Inherited: ") {
			t.Errorf("Expected rationale to start with 'Inherited: ', got %s", L1Overrides.SensitivityRationale)
		}
		if !strings.Contains(L1Overrides.SensitivityRationale, "Production handles customer data") {
			t.Error("Expected inherited rationale to contain original text")
		}
	})

	t.Run("Inherits criticality from SegL1 when empty", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"staging": {
					ID:                   "staging",
					Name:                 "Staging",
					Sensitivity:          "D",
					SensitivityRationale: "Staging contains no production data.",
					Criticality:          "5",
					CriticalityRationale: "Staging downtime impacts development velocity only.",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"staging": {
							Sensitivity:          "C",
							SensitivityRationale: "Custom sensitivity for application staging.",
							// Empty criticality - should inherit
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["app"].L1Overrides["staging"]
		if L1Overrides.Criticality != "5" {
			t.Errorf("Expected criticality '5', got %s", L1Overrides.Criticality)
		}
		if !strings.HasPrefix(L1Overrides.CriticalityRationale, "Inherited: ") {
			t.Errorf("Expected rationale to start with 'Inherited: ', got %s", L1Overrides.CriticalityRationale)
		}
	})

	t.Run("Does not inherit when sensitivity is set", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Sensitivity:          "A",
					SensitivityRationale: "Environment default rationale.",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"sec": {
					Name:        "Security",
					ID:          "sec",
					Description: "Security domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							Sensitivity:          "B",
							SensitivityRationale: "Custom sensitivity for security domain.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["sec"].L1Overrides["prod"]
		if L1Overrides.Sensitivity != "B" {
			t.Errorf("Expected sensitivity to remain 'B', got %s", L1Overrides.Sensitivity)
		}
		if strings.HasPrefix(L1Overrides.SensitivityRationale, "Inherited: ") {
			t.Error("Expected custom rationale to not be replaced with inherited prefix")
		}
	})

	t.Run("Inherits compliance requirements when nil", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox"},
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							// nil ComplianceReqs - should inherit
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.ComplianceReqs) != 2 {
			t.Errorf("Expected 2 inherited compliance reqs, got %d", len(L1Overrides.ComplianceReqs))
		}
		if L1Overrides.ComplianceReqs[0] != "pci-dss" || L1Overrides.ComplianceReqs[1] != "sox" {
			t.Errorf("Expected inherited compliance reqs [pci-dss, sox], got %v", L1Overrides.ComplianceReqs)
		}
	})

	t.Run("Does not inherit compliance requirements when set", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox", "hipaa"},
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs:       []string{"pci-dss"}, // Custom subset
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
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

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.ComplianceReqs) != 1 {
			t.Errorf("Expected custom compliance reqs to remain (1 item), got %d", len(L1Overrides.ComplianceReqs))
		}
		if L1Overrides.ComplianceReqs[0] != "pci-dss" {
			t.Errorf("Expected compliance req [pci-dss], got %v", L1Overrides.ComplianceReqs)
		}
	})

	t.Run("Populates CompReqs map with full details", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID: "prod",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs:       []string{"pci-dss", "sox"},
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["app"].L1Overrides["prod"]
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
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID: "prod",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							ComplianceReqs:       []string{"pci-dss", "invalid-scope"},
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		ApplyInheritance(&txy)

		L1Overrides := txy.SegL2s["app"].L1Overrides["prod"]
		if len(L1Overrides.CompReqs) != 1 {
			t.Errorf("Expected only 1 valid entry in CompReqs map, got %d", len(L1Overrides.CompReqs))
		}
		if _, ok := L1Overrides.CompReqs["invalid-scope"]; ok {
			t.Error("Expected invalid-scope to be skipped in CompReqs map")
		}
	})

	t.Run("Handles multiple environments in one SegL2", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Sensitivity:          "A",
					SensitivityRationale: "Production sensitivity.",
					Criticality:          "1",
					CriticalityRationale: "Production criticality.",
				},
				"staging": {
					ID:                   "staging",
					Sensitivity:          "D",
					SensitivityRationale: "Staging sensitivity.",
					Criticality:          "5",
					CriticalityRationale: "Staging criticality.",
				},
			},
			SegL2s: map[string]domain.SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {
							// Should inherit from prod SegL1
						},
						"staging": {
							// Should inherit from staging SegL1
						},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		ApplyInheritance(&txy)

		prodDetails := txy.SegL2s["app"].L1Overrides["prod"]
		if prodDetails.Sensitivity != "A" {
			t.Errorf("Expected prod to inherit 'A', got %s", prodDetails.Sensitivity)
		}

		stagingDetails := txy.SegL2s["app"].L1Overrides["staging"]
		if stagingDetails.Sensitivity != "D" {
			t.Errorf("Expected staging to inherit 'D', got %s", stagingDetails.Sensitivity)
		}
	})
}

func TestCompleteAndValidateTaxonomy(t *testing.T) {
	t.Run("Complete valid taxonomy passes all validations", func(t *testing.T) {
		txy := domain.Taxonomy{
			ApiVersion: "v1beta1",
			SegL1s: map[string]domain.SegL1{
				"shared-service": {
					ID:                   "shared-service",
					Name:                 "Shared Service",
					Description:          "Shared service environment hosting cross-account resources and centralized services with connectivity.",
					Sensitivity:          "A",
					SensitivityRationale: "Shared services represent highest risk from lateral movement perspective and bridge between environments.",
					Criticality:          "1",
					CriticalityRationale: "All environments depend on shared services for core functionality making outages highly impactful.",
					ComplianceReqs:       []string{"pci-dss", "sox"},
				},
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Description:          "Production environment with strict security controls for customer-facing services and sensitive data.",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data including PII and financial transactions requiring highest classification.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate incident response.",
					ComplianceReqs:       []string{"pci-dss"},
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
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
			SensitivityLevels: []string{"A", "B", "C", "D"},
			CriticalityLevels: []string{"1", "2", "3", "4", "5"},
		}

		valid := CompleteAndValidateTaxonomy(&txy)
		if !valid {
			t.Error("Expected complete validation to pass for valid taxonomy")
		}
	})

	t.Run("Invalid environment compliance fails", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"shared-service": {
					ID:                   "shared-service",
					Name:                 "Shared Service",
					Description:          "Shared service environment hosting cross-account resources and centralized services with connectivity.",
					Sensitivity:          "A",
					SensitivityRationale: "Shared services represent highest risk from lateral movement perspective and bridge between environments.",
					Criticality:          "1",
					CriticalityRationale: "All environments depend on shared services for core functionality making outages highly impactful.",
					ComplianceReqs:       []string{"invalid-scope"},
				},
			},
			SegL2s:   map[string]domain.SegL2{},
			CompReqs: map[string]domain.CompReq{},
		}

		valid := CompleteAndValidateTaxonomy(&txy)
		if valid {
			t.Error("Expected validation to fail for invalid environment compliance")
		}
	})

	t.Run("Missing shared-service environment fails", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Description:          "Production environment with strict security controls for customer-facing services and sensitive data.",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data including PII and financial transactions requiring highest classification.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue streams requiring immediate incident response.",
				},
			},
			SegL2s:   map[string]domain.SegL2{},
			CompReqs: map[string]domain.CompReq{},
		}

		valid := CompleteAndValidateTaxonomy(&txy)
		if valid {
			t.Error("Expected validation to fail for missing shared-service")
		}
	})

	t.Run("Invalid SegL2 references fail", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"shared-service": {
					ID:                   "shared-service",
					Name:                 "Shared Service",
					Description:          "Shared service environment hosting cross-account resources and centralized services with connectivity.",
					Sensitivity:          "A",
					SensitivityRationale: "Shared services represent highest risk from lateral movement perspective and bridge between environments.",
					Criticality:          "1",
					CriticalityRationale: "All environments depend on shared services for core functionality making outages highly impactful.",
					ComplianceReqs:       []string{"pci-dss"},
				},
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
			CompReqs: map[string]domain.CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid := CompleteAndValidateTaxonomy(&txy)
		if valid {
			t.Error("Expected validation to fail for invalid SegL2 environment reference")
		}
	})
}

