package validation

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	. "github.com/kvql/bunsceal/pkg/taxonomy/testhelpers"
)

func TestValidateL1Definitions(t *testing.T) {
	t.Run("Valid compliance requirements pass", func(t *testing.T) {
		txy := WithSegL1(
			WithStandardCompReqs(NewTestTaxonomy()),
			"prod",
			NewProdSegL1(),
		)

		valid := ValidateL1Comp(txy)
		AssertValidationPasses(t, valid, "Valid compliance requirements")
	})

	t.Run("Invalid compliance requirement fails", func(t *testing.T) {
		txy := WithSegL1(
			WithCompReq(NewTestTaxonomy(), "pci-dss", NewCompReq(
				"PCI DSS",
				"Payment Card Industry Data Security Standard",
				"https://www.pcisecuritystandards.org/",
			)),
			"prod",
			NewSegL1("prod", "Production", "A", "1", []string{"invalid-scope"}),
		)

		valid := ValidateL1Comp(txy)
		AssertValidationFails(t, valid, "Invalid compliance requirement")
	})

	t.Run("Empty compliance requirements pass", func(t *testing.T) {
		txy := WithSegL1(
			NewTestTaxonomy(),
			"staging",
			NewStagingSegL1(),
		)

		valid := ValidateL1Comp(txy)
		AssertValidationPasses(t, valid, "Empty compliance requirements")
	})
}

func TestValidateL2Definition(t *testing.T) {
	t.Run("Valid security domains pass", func(t *testing.T) {
		txy := WithSegL2(
			WithSegL1(
				WithCompReq(NewTestTaxonomy(), "pci-dss", NewCompReq(
					"PCI DSS",
					"Payment Card Industry Data Security Standard",
					"https://www.pcisecuritystandards.org/",
				)),
				"prod",
				NewSegL1("prod", "Production", "A", "1", []string{}),
			),
			"app",
			NewAppSegL2(),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationPasses(t, valid, "Valid security domains")
		AssertFailureCount(t, failures, 0, "Valid security domains")
	})

	t.Run("Invalid compliance requirement in SegL2 fails", func(t *testing.T) {
		txy := WithSegL2(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSegL2("app", "Application", map[string]domain.L1Overrides{
				"prod": {ComplianceReqs: []string{"invalid-scope"}},
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationFails(t, valid, "Invalid compliance requirement in SegL2")
		AssertMinFailures(t, failures, 1, "Invalid compliance requirement in SegL2")
	})

	t.Run("Invalid environment ID in SegL2 fails", func(t *testing.T) {
		txy := WithSegL2(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSegL2("app", "Application", map[string]domain.L1Overrides{
				"invalid-env": NewL1Override("A", "1", []string{}),
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationFails(t, valid, "Invalid environment ID in SegL2")
		AssertMinFailures(t, failures, 1, "Invalid environment ID in SegL2")
	})

	t.Run("Multiple validation failures counted", func(t *testing.T) {
		txy := WithSegL2(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSegL2("app", "Application", map[string]domain.L1Overrides{
				"prod":        {ComplianceReqs: []string{"invalid1", "invalid2"}},
				"invalid-env": {ComplianceReqs: []string{"invalid3"}},
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationFails(t, valid, "Multiple validation failures")
		// Should have failures for: 3 invalid comp reqs + 1 invalid env = 4 failures
		AssertMinFailures(t, failures, 3, "Multiple validation failures")
	})
}
