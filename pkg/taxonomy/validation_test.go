package taxonomy

import (
	"testing"
)

func TestValidateEnv(t *testing.T) {
	t.Run("Valid compliance requirements pass", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		valid := txy.validateEnv()
		if !valid {
			t.Error("Expected validation to pass for valid compliance requirements")
		}
	})

	t.Run("Invalid compliance requirement fails", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid := txy.validateEnv()
		if valid {
			t.Error("Expected validation to fail for invalid compliance requirement")
		}
	})

	t.Run("Empty compliance requirements pass", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			CompReqs: map[string]CompReq{},
		}

		valid := txy.validateEnv()
		if !valid {
			t.Error("Expected validation to pass for empty compliance requirements")
		}
	})
}

func TestValidateSecurityDomains(t *testing.T) {
	t.Run("Valid security domains pass", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain for core business services",
					L1Overrides: map[string]EnvDetails{
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
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid, failures := txy.validateSecurityDomains()
		if !valid {
			t.Errorf("Expected validation to pass, got %d failures", failures)
		}
		if failures != 0 {
			t.Errorf("Expected 0 failures, got %d", failures)
		}
	})

	t.Run("Invalid compliance requirement in SegL2 fails", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							ComplianceReqs: []string{"invalid-scope"},
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		valid, failures := txy.validateSecurityDomains()
		if valid {
			t.Error("Expected validation to fail for invalid compliance requirement")
		}
		if failures == 0 {
			t.Error("Expected at least 1 failure")
		}
	})

	t.Run("Invalid environment ID in SegL2 fails", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"invalid-env": {
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
							Criticality:          "1",
							CriticalityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		valid, failures := txy.validateSecurityDomains()
		if valid {
			t.Error("Expected validation to fail for invalid environment ID")
		}
		if failures == 0 {
			t.Error("Expected at least 1 failure")
		}
	})

	t.Run("Multiple validation failures counted", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {ID: "prod", Name: "Production"},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							ComplianceReqs: []string{"invalid1", "invalid2"},
						},
						"invalid-env": {
							ComplianceReqs: []string{"invalid3"},
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		valid, failures := txy.validateSecurityDomains()
		if valid {
			t.Error("Expected validation to fail")
		}
		// Should have failures for: 3 invalid comp reqs + 1 invalid env = 4 failures
		if failures < 3 {
			t.Errorf("Expected at least 3 failures, got %d", failures)
		}
	})
}

func TestCompleteAndValidateTaxonomy(t *testing.T) {
	t.Run("Complete valid taxonomy passes all validations", func(t *testing.T) {
		txy := Taxonomy{
			ApiVersion: "v1beta1",
			SegL1s: map[string]SegL1{
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
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain for core business services",
					L1Overrides: map[string]EnvDetails{
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
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
			SensitivityLevels: []string{"A", "B", "C", "D"},
			CriticalityLevels: []string{"1", "2", "3", "4", "5"},
		}

		valid := txy.CompleteAndValidateTaxonomy()
		if !valid {
			t.Error("Expected complete validation to pass for valid taxonomy")
		}
	})

	t.Run("Invalid environment compliance fails", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			SegL2s:   map[string]SegL2{},
			CompReqs: map[string]CompReq{},
		}

		valid := txy.CompleteAndValidateTaxonomy()
		if valid {
			t.Error("Expected validation to fail for invalid environment compliance")
		}
	})

	t.Run("Missing shared-service environment fails", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			SegL2s:   map[string]SegL2{},
			CompReqs: map[string]CompReq{},
		}

		valid := txy.CompleteAndValidateTaxonomy()
		if valid {
			t.Error("Expected validation to fail for missing shared-service")
		}
	})

	t.Run("Invalid SegL2 references fail", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"invalid-env": {
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
							Criticality:          "1",
							CriticalityRationale: "Test rationale with sufficient length to meet minimum requirements for validation.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		valid := txy.CompleteAndValidateTaxonomy()
		if valid {
			t.Error("Expected validation to fail for invalid SegL2 environment reference")
		}
	})
}
