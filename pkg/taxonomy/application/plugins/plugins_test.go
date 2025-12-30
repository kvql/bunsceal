package plugins

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
)

func TestLoadPlugins(t *testing.T) {
	t.Run("Loads classifications plugin when config present", func(t *testing.T) {
		p := make(Plugins)
		cfg := ConfigPlugins{
			Classifications: newTestConfig(true, 10),
		}

		err := p.LoadPlugins(cfg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if _, exists := p["classifications"]; !exists {
			t.Error("Expected classifications plugin to be loaded")
		}
	})

	t.Run("Returns nil error for nil config", func(t *testing.T) {
		p := make(Plugins)
		cfg := ConfigPlugins{
			Classifications: nil,
		}

		err := p.LoadPlugins(cfg)

		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}
		if len(p) != 0 {
			t.Error("Expected no plugins to be loaded when config is nil")
		}
	})
}

func TestApplyInheritanceAndValidate(t *testing.T) {
	t.Run("Skips inheritance when LabelInheritance=false", func(t *testing.T) {
		config := newTestConfig(false, 10) // inheritance disabled
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent rationale here"),
		})
		child := newTestSeg("child", []string{})

		plugs := make(Plugins)
		plugs["classifications"] = plugin
		errs := plugs.ApplyPluginInheritanceAndValidate(*parent, child)

		if len(errs) > 0 {
			t.Fatalf("Expected no errors, got %v", errs)
		}
		// Child should not have inherited labels
		if len(child.LabelNamespaces[testNs]) != 0 {
			t.Error("Expected child to have no inherited labels when inheritance is disabled")
		}
	})

	t.Run("Inherits parent labels to child", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent rationale here"),
		})
		child := newTestSeg("child", []string{
			label("other_key", "child_value"), // child has some labels in namespace
		})

		plugs := make(Plugins)
		plugs["classifications"] = plugin
		errs := plugs.ApplyPluginInheritanceAndValidate(*parent, child)

		if len(errs) > 0 {
			t.Fatalf("Expected no errors, got %v", errs)
		}
		// Child should have inherited parent's labels
		if child.LabelNamespaces[testNs]["sensitivity"] != "high" {
			t.Errorf("Expected child to inherit sensitivity=high, got %q", child.LabelNamespaces[testNs]["sensitivity"])
		}
	})

	t.Run("Updates ParsedLabels during inheritance", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent rationale here"),
		})
		child := newTestSeg("child", []string{
			label("other_key", "child_value"),
		})

		plugs := make(Plugins)
		plugs["classifications"] = plugin
		errs := plugs.ApplyPluginInheritanceAndValidate(*parent, child)

		if len(errs) > 0 {
			t.Fatalf("Expected no errors, got %v", errs)
		}
		// ParsedLabels should also be updated
		expectedKey := testNs + "/sensitivity"
		if child.ParsedLabels[expectedKey] != "high" {
			t.Errorf("Expected ParsedLabels[%q]=high, got %q", expectedKey, child.ParsedLabels[expectedKey])
		}
	})

	t.Run("Does not override existing child labels", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent rationale here"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "low"), // child has explicit value
			label("sensitivity_rationale", "Child rationale here"),
		})

		plugs := make(Plugins)
		plugs["classifications"] = plugin
		errs := plugs.ApplyPluginInheritanceAndValidate(*parent, child)

		if len(errs) > 0 {
			t.Fatalf("Expected no errors, got %v", errs)
		}
		// Child's explicit value should not be overridden
		if child.LabelNamespaces[testNs]["sensitivity"] != "low" {
			t.Errorf("Expected child to keep sensitivity=low, got %q", child.LabelNamespaces[testNs]["sensitivity"])
		}
	})

	t.Run("Child retains explicit labels when parent has none", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{}) // parent has no labels in namespace
		child := newTestSeg("child", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Child rationale"),
		})

		plugs := make(Plugins)
		plugs["classifications"] = plugin
		errs := plugs.ApplyPluginInheritanceAndValidate(*parent, child)

		if len(errs) > 0 {
			t.Fatalf("Expected no errors, got %v", errs)
		}
		// Child should retain its explicit labels
		if child.LabelNamespaces[testNs]["sensitivity"] != "high" {
			t.Errorf("Expected child to keep sensitivity=high, got %q", child.LabelNamespaces[testNs]["sensitivity"])
		}
	})
}

func TestNsPrefix(t *testing.T) {
	t.Run("Verifies NsPrefix constant value", func(t *testing.T) {
		expected := "bunsceal.plugin."
		if NsPrefix != expected {
			t.Errorf("Expected NsPrefix %q, got %q", expected, NsPrefix)
		}
	})
}

func TestValidateAllSegments(t *testing.T) {
	t.Run("Returns no errors for empty plugins list", func(t *testing.T) {
		p := make(Plugins)
		l1s := map[string]domain.Seg{
			"seg1": *newTestSeg("seg1", []string{label("sensitivity", "high")}),
		}
		l2s := map[string]domain.Seg{}

		errs := p.ValidateAllSegments(l1s, l2s)

		if len(errs) != 0 {
			t.Errorf("Expected no errors with empty plugins, got %d", len(errs))
		}
	})

	t.Run("Skips segments with no labels", func(t *testing.T) {
		config := newTestConfig(true, 10)
		p := make(Plugins)
		p["classifications"] = NewClassificationPlugin(config, NsPrefix)

		l1s := map[string]domain.Seg{
			"seg1": *newTestSeg("seg1", []string{}), // no labels
		}
		l2s := map[string]domain.Seg{}

		errs := p.ValidateAllSegments(l1s, l2s)

		// Should skip validation for segment with no labels
		if len(errs) != 0 {
			t.Errorf("Expected no errors for segments without labels, got %d", len(errs))
		}
	})

	t.Run("Collects errors from multiple segments", func(t *testing.T) {
		config := newTestConfig(true, 10)
		p := make(Plugins)
		p["classifications"] = NewClassificationPlugin(config, NsPrefix)

		// Both segments have invalid labels (missing rationale)
		l1s := map[string]domain.Seg{
			"seg1": *newTestSeg("seg1", []string{label("sensitivity", "high")}),
		}
		l2s := map[string]domain.Seg{
			"seg2": *newTestSeg("seg2", []string{label("sensitivity", "low")}),
		}

		errs := p.ValidateAllSegments(l1s, l2s)

		// Should collect errors from both segments (not fail fast)
		if len(errs) == 0 {
			t.Error("Expected validation errors for segments with invalid labels")
		}
		// Should have errors from both seg1 and seg2
		hasL1Error := false
		hasL2Error := false
		for _, err := range errs {
			if contains(err.Error(), "L1 segment seg1") {
				hasL1Error = true
			}
			if contains(err.Error(), "L2 segment seg2") {
				hasL2Error = true
			}
		}
		if !hasL1Error {
			t.Error("Expected error for L1 segment seg1")
		}
		if !hasL2Error {
			t.Error("Expected error for L2 segment seg2")
		}
	})

	t.Run("Returns no errors for valid segments", func(t *testing.T) {
		config := newTestConfig(true, 10)
		p := make(Plugins)
		p["classifications"] = NewClassificationPlugin(config, NsPrefix)

		l1s := map[string]domain.Seg{
			"seg1": *newTestSeg("seg1", []string{
				label("sensitivity", "high"),
				label("sensitivity_rationale", "Valid rationale with enough length"),
			}),
		}
		l2s := map[string]domain.Seg{}

		errs := p.ValidateAllSegments(l1s, l2s)

		if len(errs) != 0 {
			t.Errorf("Expected no errors for valid segment, got %d: %v", len(errs), errs)
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[:len(substr)] == substr || (len(s) > len(substr) && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
