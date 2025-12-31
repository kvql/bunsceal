package validation

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	. "github.com/kvql/bunsceal/pkg/domain/testhelpers"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

func TestValidateL2Definition(t *testing.T) {
	t.Run("Valid L2 with valid L1 parent passes", func(t *testing.T) {
		pluginMap := make(plugins.Plugins)
		txy := WithSeg(
			WithSegL1(
				NewTestTaxonomy(),
				"prod",
				NewSegL1("prod", "Test", "A", "1", nil),
			),
			"app",
			NewSeg("app", "Test", map[string]domain.L1Overrides{
				"prod": NewL1Override("A", "1", nil),
			}),
		)

		valid, failures := ValidateL2Definition(txy, pluginMap)
		AssertValidationPasses(t, valid, "Valid L2 with valid L1 parent")
		AssertFailureCount(t, failures, 0, "Valid L2 with valid L1 parent")
	})

	t.Run("Invalid L1 parent reference fails", func(t *testing.T) {
		pluginMap := make(plugins.Plugins)
		txy := WithSeg(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", nil)),
			"app",
			NewSeg("app", "Application", map[string]domain.L1Overrides{
				"invalid-env": NewL1Override("A", "1", nil),
			}),
		)

		valid, failures := ValidateL2Definition(txy, pluginMap)
		AssertValidationFails(t, valid, "Invalid L1 parent reference")
		AssertMinFailures(t, failures, 1, "Invalid L1 parent reference")
	})

	t.Run("Multiple validation failures counted", func(t *testing.T) {
		pluginMap := make(plugins.Plugins)
		txy := WithSeg(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", nil)),
			"app",
			NewSeg("app", "Application", map[string]domain.L1Overrides{
				"invalid-env1": NewL1Override("A", "1", nil),
				"invalid-env2": NewL1Override("A", "1", nil),
			}),
		)

		valid, failures := ValidateL2Definition(txy, pluginMap)
		AssertValidationFails(t, valid, "Multiple validation failures")
		AssertMinFailures(t, failures, 2, "Multiple validation failures")
	})
}
