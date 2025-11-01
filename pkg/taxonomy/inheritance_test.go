package taxonomy

import (
	"strings"
	"testing"
)

func TestApplyInheritance(t *testing.T) {
	t.Run("Inherits sensitivity from SegL1 when empty", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID:                   "prod",
					Name:                 "Production",
					Sensitivity:          "A",
					SensitivityRationale: "Production handles customer data requiring highest classification level.",
					Criticality:          "1",
					CriticalityRationale: "Production outages directly impact customers and revenue.",
				},
			},
			SegL2s: map[string]SegL2{
				"infra": {
					Name:        "Infrastructure",
					ID:          "infra",
					Description: "Infrastructure domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							// Empty sensitivity and rationale - should inherit
							Criticality:          "2",
							CriticalityRationale: "Custom criticality for infrastructure.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["infra"].L1Overrides["prod"]
		if envDetails.Sensitivity != "A" {
			t.Errorf("Expected sensitivity 'A', got %s", envDetails.Sensitivity)
		}
		if !strings.HasPrefix(envDetails.SensitivityRationale, "Inherited: ") {
			t.Errorf("Expected rationale to start with 'Inherited: ', got %s", envDetails.SensitivityRationale)
		}
		if !strings.Contains(envDetails.SensitivityRationale, "Production handles customer data") {
			t.Error("Expected inherited rationale to contain original text")
		}
	})

	t.Run("Inherits criticality from SegL1 when empty", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"staging": {
					ID:                   "staging",
					Name:                 "Staging",
					Sensitivity:          "D",
					SensitivityRationale: "Staging contains no production data.",
					Criticality:          "5",
					CriticalityRationale: "Staging downtime impacts development velocity only.",
				},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"staging": {
							Sensitivity:          "C",
							SensitivityRationale: "Custom sensitivity for application staging.",
							// Empty criticality - should inherit
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["app"].L1Overrides["staging"]
		if envDetails.Criticality != "5" {
			t.Errorf("Expected criticality '5', got %s", envDetails.Criticality)
		}
		if !strings.HasPrefix(envDetails.CriticalityRationale, "Inherited: ") {
			t.Errorf("Expected rationale to start with 'Inherited: ', got %s", envDetails.CriticalityRationale)
		}
	})

	t.Run("Does not inherit when sensitivity is set", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID:                   "prod",
					Sensitivity:          "A",
					SensitivityRationale: "Environment default rationale.",
				},
			},
			SegL2s: map[string]SegL2{
				"sec": {
					Name:        "Security",
					ID:          "sec",
					Description: "Security domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							Sensitivity:          "B",
							SensitivityRationale: "Custom sensitivity for security domain.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["sec"].L1Overrides["prod"]
		if envDetails.Sensitivity != "B" {
			t.Errorf("Expected sensitivity to remain 'B', got %s", envDetails.Sensitivity)
		}
		if strings.HasPrefix(envDetails.SensitivityRationale, "Inherited: ") {
			t.Error("Expected custom rationale to not be replaced with inherited prefix")
		}
	})

	t.Run("Inherits compliance requirements when nil", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox"},
				},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							// nil ComplianceReqs - should inherit
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["app"].L1Overrides["prod"]
		if len(envDetails.ComplianceReqs) != 2 {
			t.Errorf("Expected 2 inherited compliance reqs, got %d", len(envDetails.ComplianceReqs))
		}
		if envDetails.ComplianceReqs[0] != "pci-dss" || envDetails.ComplianceReqs[1] != "sox" {
			t.Errorf("Expected inherited compliance reqs [pci-dss, sox], got %v", envDetails.ComplianceReqs)
		}
	})

	t.Run("Does not inherit compliance requirements when set", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID:             "prod",
					ComplianceReqs: []string{"pci-dss", "sox", "hipaa"},
				},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							ComplianceReqs:       []string{"pci-dss"}, // Custom subset
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
				"hipaa":   {Name: "HIPAA", Description: "Health Insurance Portability and Accountability Act", ReqsLink: "https://www.hhs.gov/hipaa/"},
			},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["app"].L1Overrides["prod"]
		if len(envDetails.ComplianceReqs) != 1 {
			t.Errorf("Expected custom compliance reqs to remain (1 item), got %d", len(envDetails.ComplianceReqs))
		}
		if envDetails.ComplianceReqs[0] != "pci-dss" {
			t.Errorf("Expected compliance req [pci-dss], got %v", envDetails.ComplianceReqs)
		}
	})

	t.Run("Populates CompReqs map with full details", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID: "prod",
				},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							ComplianceReqs:       []string{"pci-dss", "sox"},
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
				"sox":     {Name: "SOX", Description: "Sarbanes-Oxley Act", ReqsLink: "https://www.sox-online.com/"},
			},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["app"].L1Overrides["prod"]
		if len(envDetails.CompReqs) != 2 {
			t.Errorf("Expected 2 entries in CompReqs map, got %d", len(envDetails.CompReqs))
		}
		if compReq, ok := envDetails.CompReqs["pci-dss"]; !ok {
			t.Error("Expected pci-dss in CompReqs map")
		} else {
			if compReq.Name != "PCI DSS" {
				t.Errorf("Expected PCI DSS name, got %s", compReq.Name)
			}
		}
		if compReq, ok := envDetails.CompReqs["sox"]; !ok {
			t.Error("Expected sox in CompReqs map")
		} else {
			if compReq.Name != "SOX" {
				t.Errorf("Expected SOX name, got %s", compReq.Name)
			}
		}
	})

	t.Run("Skips invalid compliance requirements in CompReqs map", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
				"prod": {
					ID: "prod",
				},
			},
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							ComplianceReqs:       []string{"pci-dss", "invalid-scope"},
							Sensitivity:          "A",
							SensitivityRationale: "Test rationale.",
						},
					},
				},
			},
			CompReqs: map[string]CompReq{
				"pci-dss": {Name: "PCI DSS", Description: "Payment Card Industry Data Security Standard", ReqsLink: "https://www.pcisecuritystandards.org/"},
			},
		}

		txy.ApplyInheritance()

		envDetails := txy.SegL2s["app"].L1Overrides["prod"]
		if len(envDetails.CompReqs) != 1 {
			t.Errorf("Expected only 1 valid entry in CompReqs map, got %d", len(envDetails.CompReqs))
		}
		if _, ok := envDetails.CompReqs["invalid-scope"]; ok {
			t.Error("Expected invalid-scope to be skipped in CompReqs map")
		}
	})

	t.Run("Handles multiple environments in one SegL2", func(t *testing.T) {
		txy := Taxonomy{
			SegL1s: map[string]SegL1{
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
			SegL2s: map[string]SegL2{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Overrides: map[string]EnvDetails{
						"prod": {
							// Should inherit from prod SegL1
						},
						"staging": {
							// Should inherit from staging SegL1
						},
					},
				},
			},
			CompReqs: map[string]CompReq{},
		}

		txy.ApplyInheritance()

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
