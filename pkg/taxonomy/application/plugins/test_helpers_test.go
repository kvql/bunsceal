package plugins

import "github.com/kvql/bunsceal/pkg/domain"

// Namespace constants for test labels
const testNs = "bunsceal.plugin.classifications"
const complianceTestNs = "bunsceal.plugin.compliance"

func newTestConfig(inheritance bool, rationaleLen int) *ClassificationsConfig {
	return &ClassificationsConfig{
		Common: PluginsCommonSettings{
			LabelInheritance:  inheritance,
			RequireCompleteL1: true,
		},
		RationaleLength: rationaleLen,
		Definitions: map[string]ClassificationDefinition{
			"sensitivity": {
				DescriptiveName: "Data Sensitivity",
				EnforceOrder:    true,
				Values:          map[string]string{"high": "High sensitivity", "low": "Low sensitivity"},
				Order:           []string{"high", "low"},
			},
		},
	}
}

func newTestConfigWithOrder(inheritance bool, rationaleLen int, enforceOrder bool) *ClassificationsConfig {
	return &ClassificationsConfig{
		Common: PluginsCommonSettings{
			LabelInheritance:  inheritance,
			RequireCompleteL1: true,
		},
		RationaleLength: rationaleLen,
		Definitions: map[string]ClassificationDefinition{
			"sensitivity": {
				DescriptiveName: "Data Sensitivity",
				EnforceOrder:    enforceOrder,
				Values:          map[string]string{"high": "High", "medium": "Medium", "low": "Low"},
				Order:           []string{"high", "medium", "low"},
			},
			"criticality": {
				DescriptiveName: "Business Criticality",
				EnforceOrder:    enforceOrder,
				Values:          map[string]string{"1": "Critical", "2": "High", "3": "Medium", "4": "Low"},
				Order:           []string{"1", "2", "3", "4"},
			},
		},
	}
}

func newTestSeg(id string, labels []string) *domain.Seg {
	seg := &domain.Seg{ID: id, Name: id, Labels: labels}
	seg.ParseLabels()
	return seg
}

// Helper to build fully qualified label
func label(key, value string) string {
	return testNs + "/" + key + ":" + value
}

// Helper to build compliance label
func complianceLabel(key, value string) string {
	return complianceTestNs + "/" + key + ":" + value
}

// Helper to create test compliance config
func newComplianceTestConfig(inheritance bool, rationaleLen int, enforceHierarchy bool) *ComplianceConfig {
	return &ComplianceConfig{
		Common:                PluginsCommonSettings{LabelInheritance: inheritance},
		RationaleLength:       rationaleLen,
		EnforceScopeHierarchy: enforceHierarchy,
		Definitions: map[string]ComplianceDefinition{
			"pci-dss": {
				DescriptiveName:  "PCI DSS",
				Description:      "Payment Card Industry Data Security Standard",
				RequirementsLink: "https://www.pcisecuritystandards.org/",
			},
			"soc2": {
				DescriptiveName:  "SOC 2",
				Description:      "Service Organisation Control 2",
				RequirementsLink: "",
			},
		},
	}
}
