package plugins

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
)

func TestComplianceValidateLabels_ValidScope(t *testing.T) {
	t.Run("Accepts in-scope value", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "This environment processes payment card data"),
		})

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected valid result for in-scope, got errors: %v", result.Errors)
		}
	})

	t.Run("Accepts out-of-scope value", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "out-of-scope"),
			complianceLabel("pci-dss_rationale", "No payment card data processed"),
		})

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected valid result for out-of-scope, got errors: %v", result.Errors)
		}
	})

	t.Run("Rejects invalid scope value", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "maybe-scope"),
			complianceLabel("pci-dss_rationale", "Invalid scope value"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for invalid scope value")
		}
	})
}

func TestComplianceValidateLabels_MissingRationale(t *testing.T) {
	t.Run("Fails on missing rationale", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "in-scope"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for missing rationale")
		}
	})

	t.Run("Fails on rationale below minimum length", func(t *testing.T) {
		config := newComplianceTestConfig(true, 50, true) // require 50 chars
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Too short"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for rationale below minimum length")
		}
	})

	t.Run("Fails on scope without paired rationale", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("soc2_rationale", "Has rationale without scope"),
		})

		result := plugin.ValidateLabels(seg)

		if result.Valid {
			t.Error("Expected validation to fail for rationale without scope")
		}
	})
}

func TestComplianceValidateLabels_L1NotRequired(t *testing.T) {
	t.Run("L1 segments do not require complete compliance labels", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		seg := newTestSeg("test-seg", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Only PCI DSS is in scope, SOC2 is not defined"),
		})
		seg.Level = "1" // L1 segment

		result := plugin.ValidateLabels(seg)

		if !result.Valid {
			t.Errorf("Expected L1 with partial compliance labels to be valid, got: %v", result.Errors)
		}
	})
}

func TestComplianceValidateRelationship_HierarchyEnforced(t *testing.T) {
	t.Run("Passes when parent has requirement defined and child is in-scope", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Parent is in scope"),
		})
		child := newTestSeg("child", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Child is in scope"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when parent has requirement, got: %v", errs)
		}
	})

	t.Run("Fails when child is in-scope but parent doesn't have requirement defined", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{}) // No compliance labels
		child := newTestSeg("child", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Child is in scope"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) == 0 {
			t.Error("Expected error when child is in-scope but parent doesn't have requirement")
		}
	})

	t.Run("Passes when child is out-of-scope", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{}) // No compliance labels
		child := newTestSeg("child", []string{
			complianceLabel("pci-dss", "out-of-scope"),
			complianceLabel("pci-dss_rationale", "Child is out of scope"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when child is out-of-scope, got: %v", errs)
		}
	})

	t.Run("Override sets child effective value", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{}) // No compliance labels
		child := newTestSeg("child", []string{
			complianceLabel("pci-dss", "out-of-scope"),
			complianceLabel("pci-dss_rationale", "Child base is out of scope"),
		})
		// Override makes child "in-scope" for this parent
		override := domain.L1Overrides{Labels: []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Override makes it in-scope"),
		}}
		override.ParseLabels()
		child.L1Overrides = map[string]domain.L1Overrides{"parent": override}

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) == 0 {
			t.Error("Expected error when override makes child in-scope but parent doesn't have requirement")
		}
	})
}

func TestComplianceValidateRelationship_HierarchyDisabled(t *testing.T) {
	t.Run("No validation when enforce_scope_hierarchy is false", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, false) // hierarchy disabled
		plugin := NewCompliancePlugin(config, NsPrefix)
		parent := newTestSeg("parent", []string{}) // No compliance labels
		child := newTestSeg("child", []string{
			complianceLabel("pci-dss", "in-scope"),
			complianceLabel("pci-dss_rationale", "Child is in scope"),
		})

		errs := plugin.ValidateRelationship(parent, child)

		if len(errs) > 0 {
			t.Errorf("Expected no errors when hierarchy disabled, got: %v", errs)
		}
	})
}

func TestComplianceGetImageData(t *testing.T) {
	t.Run("Returns correct structure for visualisation", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)

		imageData := plugin.GetImageData()

		if len(imageData) != 2 {
			t.Errorf("Expected 2 image data items, got %d", len(imageData))
		}

		for _, data := range imageData {
			if data.Namespace != "bunsceal.plugin.compliance" {
				t.Errorf("Expected namespace 'bunsceal.plugin.compliance', got %q", data.Namespace)
			}
			if len(data.OrderedValues) != 2 {
				t.Errorf("Expected 2 ordered values, got %d", len(data.OrderedValues))
			}
			if data.OrderedValues[0] != "in-scope" || data.OrderedValues[1] != "out-of-scope" {
				t.Errorf("Expected ordered values ['in-scope', 'out-of-scope'], got %v", data.OrderedValues)
			}
		}
	})
}

func TestComplianceGetEnabled(t *testing.T) {
	t.Run("Returns label_inheritance setting", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)

		if !plugin.GetEnabled() {
			t.Error("Expected GetEnabled to return true when label_inheritance is true")
		}

		config.Common.LabelInheritance = false
		plugin = NewCompliancePlugin(config, NsPrefix)

		if plugin.GetEnabled() {
			t.Error("Expected GetEnabled to return false when label_inheritance is false")
		}
	})
}

func TestComplianceGetNamespace(t *testing.T) {
	t.Run("Returns correct namespace", func(t *testing.T) {
		config := newComplianceTestConfig(true, 10, true)
		plugin := NewCompliancePlugin(config, NsPrefix)

		expectedNs := "bunsceal.plugin.compliance"
		if plugin.GetNamespace() != expectedNs {
			t.Errorf("Expected namespace %q, got %q", expectedNs, plugin.GetNamespace())
		}
	})
}
