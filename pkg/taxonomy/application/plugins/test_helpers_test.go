package plugins

import "github.com/kvql/bunsceal/pkg/domain"

// Namespace constant for test labels
const testNs = "bunsceal.plugin.classifications"

func newTestConfig(inheritance bool, rationaleLen int) *ClassificationsConfig {
	return &ClassificationsConfig{
		Common:          PluginsCommonSettings{LabelInheritance: inheritance},
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
		Common:          PluginsCommonSettings{LabelInheritance: inheritance},
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
