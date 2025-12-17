package plugins

import (
	"testing"
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
