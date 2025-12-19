package application

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

const testNs = "bunsceal.plugin.classifications"

func newTestPlugins(inheritance bool) *plugins.Plugins {
	config := &plugins.ClassificationsConfig{
		Common:          plugins.PluginsCommonSettings{LabelInheritance: inheritance},
		RationaleLength: 10,
		Definitions: map[string]plugins.ClassificationDefinition{
			"sensitivity": {
				DescriptiveName: "Data Sensitivity",
				Values:          map[string]string{"high": "High sensitivity", "low": "Low sensitivity"},
				Order:           []string{"high", "low"},
			},
		},
	}
	p := &plugins.Plugins{Plugins: make(map[string]plugins.Plugin)}
	p.Plugins["classifications"] = plugins.NewClassificationPlugin(config, plugins.NsPrefix)
	return p
}

func newTestSegWithLabels(id string, labels []string) domain.Seg {
	seg := domain.Seg{
		ID:     id,
		Name:   id,
		Labels: labels,
	}
	seg.ParseLabels()
	return seg
}

func label(key, value string) string {
	return testNs + "/" + key + ":" + value
}

func TestApplyInheritance_PluginLabels(t *testing.T) {
	t.Run("Inherits plugin labels from L1 parent to L2 child", func(t *testing.T) {
		parentSeg := newTestSegWithLabels("prod", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Production contains PII data"),
		})

		childSeg := newTestSegWithLabels("app", []string{
			label("other_key", "value"), // has labels in namespace but not sensitivity
		})
		childSeg.L1Parents = []string{"prod"}
		childSeg.L1Overrides = map[string]domain.L1Overrides{
			"prod": {},
		}

		txy := domain.Taxonomy{
			SegL1s:   map[string]domain.Seg{"prod": parentSeg},
			SegsL2s:  map[string]domain.Seg{"app": childSeg},
			CompReqs: map[string]domain.CompReq{},
		}

		p := newTestPlugins(true)
		ApplyInheritance(&txy, p)

		// Verify child inherited parent's sensitivity label
		if txy.SegsL2s["app"].LabelNamespaces[testNs]["sensitivity"] != "high" {
			t.Errorf("Expected child to inherit sensitivity=high, got %q",
				txy.SegsL2s["app"].LabelNamespaces[testNs]["sensitivity"])
		}
		// Verify ParsedLabels also updated
		expectedKey := testNs + "/sensitivity"
		if txy.SegsL2s["app"].ParsedLabels[expectedKey] != "high" {
			t.Errorf("Expected ParsedLabels[%q]=high, got %q",
				expectedKey, txy.SegsL2s["app"].ParsedLabels[expectedKey])
		}
	})

	t.Run("Child override label takes precedence over parent", func(t *testing.T) {
		parentSeg := newTestSegWithLabels("prod", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent rationale"),
		})

		childSeg := newTestSegWithLabels("app", []string{
			label("sensitivity", "low"), // child overrides
			label("sensitivity_rationale", "Child rationale"),
		})
		childSeg.L1Parents = []string{"prod"}
		childSeg.L1Overrides = map[string]domain.L1Overrides{
			"prod": {},
		}

		txy := domain.Taxonomy{
			SegL1s:   map[string]domain.Seg{"prod": parentSeg},
			SegsL2s:  map[string]domain.Seg{"app": childSeg},
			CompReqs: map[string]domain.CompReq{},
		}

		p := newTestPlugins(true)
		ApplyInheritance(&txy, p)

		// Child's explicit value should not be overridden
		if txy.SegsL2s["app"].LabelNamespaces[testNs]["sensitivity"] != "low" {
			t.Errorf("Expected child to keep sensitivity=low, got %q",
				txy.SegsL2s["app"].LabelNamespaces[testNs]["sensitivity"])
		}
	})

	t.Run("Multiple L1 parents - labels merged correctly", func(t *testing.T) {
		prodSeg := newTestSegWithLabels("prod", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Prod rationale"),
		})

		stagingSeg := newTestSegWithLabels("staging", []string{
			label("sensitivity", "low"),
			label("sensitivity_rationale", "Staging rationale"),
		})

		childSeg := newTestSegWithLabels("app", []string{
			label("other_key", "value"),
		})
		childSeg.L1Parents = []string{"prod", "staging"}
		childSeg.L1Overrides = map[string]domain.L1Overrides{
			"prod":    {},
			"staging": {},
		}

		txy := domain.Taxonomy{
			SegL1s:   map[string]domain.Seg{"prod": prodSeg, "staging": stagingSeg},
			SegsL2s:  map[string]domain.Seg{"app": childSeg},
			CompReqs: map[string]domain.CompReq{},
		}

		p := newTestPlugins(true)
		ApplyInheritance(&txy, p)

		// Child should have inherited from one of the parents
		// (order depends on map iteration, but value should be set)
		sens := txy.SegsL2s["app"].LabelNamespaces[testNs]["sensitivity"]
		if sens != "high" && sens != "low" {
			t.Errorf("Expected child to inherit sensitivity from a parent, got %q", sens)
		}
	})
}

func TestApplyInheritance_NilPlugins(t *testing.T) {
	t.Run("Works without plugin config", func(t *testing.T) {
		txy := domain.Taxonomy{
			SegL1s: map[string]domain.Seg{
				"prod": {
					ID:   "prod",
					Name: "Production",
				},
			},
			SegsL2s: map[string]domain.Seg{
				"app": {
					Name:        "Application",
					ID:          "app",
					Description: "Application domain",
					L1Parents:   []string{"prod"},
					L1Overrides: map[string]domain.L1Overrides{
						"prod": {},
					},
				},
			},
			CompReqs: map[string]domain.CompReq{},
		}

		// Should not panic with nil plugins
		err := ApplyInheritance(&txy, nil)
		if err != nil {
			t.Errorf("Expected no error with nil plugins, got: %v", err)
		}
	})
}
