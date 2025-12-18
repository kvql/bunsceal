package plugins

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
)

func TestNewClassificationPlugin(t *testing.T) {
	t.Run("Creates plugin with correct namespace", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)

		expectedNs := "bunsceal.plugin.classifications"
		if plugin.Namespace != expectedNs {
			t.Errorf("Expected namespace %q, got %q", expectedNs, plugin.Namespace)
		}
	})

	t.Run("Stores config reference", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)

		if plugin.Config != config {
			t.Error("Expected plugin to store config reference")
		}
	})
}

func TestValidateLabels(t *testing.T) {
	t.Run("Returns valid for complete labels", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "This contains PII data"),
		})

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected valid result, got errors: %v", result.Errors)
		}
	})

	t.Run("Fails on missing classification label", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity_rationale", "This contains PII data"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for missing classification label")
		}
	})

	t.Run("Fails on invalid classification value", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity", "invalid-value"),
			label("sensitivity_rationale", "This contains PII data"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for invalid classification value")
		}
	})

	t.Run("Fails on missing rationale", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity", "high"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for missing rationale")
		}
	})

	t.Run("Fails on rationale below minimum length", func(t *testing.T) {
		config := newTestConfig(true, 50) // require 50 chars
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Too short"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for rationale below minimum length")
		}
	})

	t.Run("L1 fails when unknown replaces valid definition", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, true) // 2 definitions
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Valid rationale here"),
			label("unknown", "foo"),
			label("unknown_rationale", "Unknown rationale"),
		})
		seg.Level = "1" // L1 requires all definitions

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected L1 validation to fail when criticality replaced by unknown")
		}
	})

	t.Run("L2 with unknown label pair is valid", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			label("unknown", "foo"),
			label("unknown_rationale", "Unknown rationale"),
		})
		seg.Level = "2" // L2 only needs valid pairs

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected L2 with unknown pair to be valid, got: %v", result.Errors)
		}
	})

	t.Run("Override with unknown label pair is valid", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{})
		seg.Level = "2"
		override := domain.L1Overrides{Labels: []string{
			label("unknown", "foo"),
			label("unknown_rationale", "Unknown rationale"),
		}}
		override.ParseLabels()
		seg.L1Overrides = map[string]domain.L1Overrides{"parent": override}

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected override with unknown pair to be valid, got: %v", result.Errors)
		}
	})
}

func TestClassificationNamespace(t *testing.T) {
	t.Run("Verifies namespace format", func(t *testing.T) {
		config := newTestConfig(true, 10)
		plugin := NewClassificationPlugin(config, NsPrefix)

		if plugin.Namespace != testNs {
			t.Errorf("Expected namespace %q, got %q", testNs, plugin.Namespace)
		}
	})
}

func TestValidateRelationship(t *testing.T) {
	t.Run("Passes when parent >= child", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, true)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Parent is high sensitivity"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "medium"),
			label("sensitivity_rationale", "Child is medium sensitivity"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when parent >= child, got: %v", errs)
		}
	})

	t.Run("Fails when child > parent", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, true)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "low"),
			label("sensitivity_rationale", "Parent is low"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Child is high"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) == 0 {
			t.Error("Expected error when child > parent")
		}
	})

	t.Run("Skips when enforce_order is false", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, false)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "low"),
			label("sensitivity_rationale", "Parent is low"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Child is high"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when enforce_order=false, got: %v", errs)
		}
	})

	t.Run("Override sets child effective value", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, true)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "low"),
			label("sensitivity_rationale", "Parent is low"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "medium"),
			label("sensitivity_rationale", "Child base is medium"),
		})
		// Override makes child "high" for this parent - more severe than parent's "low"
		override := domain.L1Overrides{Labels: []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Override is high"),
		}}
		override.ParseLabels()
		child.L1Overrides = map[string]domain.L1Overrides{"parent": override}

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) == 0 {
			t.Error("Expected error when override makes child more severe than parent")
		}
	})

	t.Run("Skips when values not in order list", func(t *testing.T) {
		config := newTestConfigWithOrder(true, 10, true)
		plugin := NewClassificationPlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			label("sensitivity", "unknown"),
			label("sensitivity_rationale", "Parent has invalid value"),
		})
		child := newTestSeg("child", []string{
			label("sensitivity", "high"),
			label("sensitivity_rationale", "Child is high"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when value not in order, got: %v", errs)
		}
	})
}
