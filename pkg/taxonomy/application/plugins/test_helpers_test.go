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
				Values:          map[string]string{"high": "High sensitivity", "low": "Low sensitivity"},
				Order:           []string{"high", "low"},
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
