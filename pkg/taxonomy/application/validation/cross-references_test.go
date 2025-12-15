package validation

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	. "github.com/kvql/bunsceal/pkg/domain/testhelpers"
)

func TestValidateL1Definitions(t *testing.T) {
	t.Run("Valid compliance requirements pass", func(t *testing.T) {
		txy := WithSegL1(
			WithStandardCompReqs(NewTestTaxonomy()),
			"prod",
			NewSegL1("prod", "Test", "A", "1", []string{"pci-dss", "sox"}),
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
			NewSegL1("prod", "Test", "A", "1", []string{"invalid-scope"}),
		)

		valid := ValidateL1Comp(txy)
		AssertValidationFails(t, valid, "Invalid compliance requirement")
	})

	t.Run("Empty compliance requirements pass", func(t *testing.T) {
		txy := WithSegL1(
			NewTestTaxonomy(),
			"staging",
			NewSegL1("staging", "Test", "D", "5", nil),
		)

		valid := ValidateL1Comp(txy)
		AssertValidationPasses(t, valid, "Empty compliance requirements")
	})
}

func TestValidateL2Definition(t *testing.T) {
	t.Run("Valid security domains pass", func(t *testing.T) {
		txy := WithSeg(
			WithSegL1(
				WithCompReq(NewTestTaxonomy(), "pci-dss", NewCompReq(
					"PCI DSS",
					"Payment Card Industry Data Security Standard",
					"https://www.pcisecuritystandards.org/",
				)),
				"prod",
				NewSegL1("prod", "Test", "A", "1", nil),
			),
			"app",
			NewSeg("app", "Test", map[string]domain.L1Overrides{
				"prod": NewL1Override("A", "1", []string{"pci-dss"}),
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationPasses(t, valid, "Valid security domains")
		AssertFailureCount(t, failures, 0, "Valid security domains")
	})

	t.Run("Invalid compliance requirement in Seg fails", func(t *testing.T) {
		txy := WithSeg(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSeg("app", "Application", map[string]domain.L1Overrides{
				"prod": {ComplianceReqs: []string{"invalid-scope"}},
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationFails(t, valid, "Invalid compliance requirement in Seg")
		AssertMinFailures(t, failures, 1, "Invalid compliance requirement in Seg")
	})

	t.Run("Invalid environment ID in Seg fails", func(t *testing.T) {
		txy := WithSeg(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSeg("app", "Application", map[string]domain.L1Overrides{
				"invalid-env": NewL1Override("A", "1", []string{}),
			}),
		)

		valid, failures := ValidateL2Definition(txy)
		AssertValidationFails(t, valid, "Invalid environment ID in Seg")
		AssertMinFailures(t, failures, 1, "Invalid environment ID in Seg")
	})

	t.Run("Multiple validation failures counted", func(t *testing.T) {
		txy := WithSeg(
			WithSegL1(NewTestTaxonomy(), "prod", NewSegL1("prod", "Production", "A", "1", []string{})),
			"app",
			NewSeg("app", "Application", map[string]domain.L1Overrides{
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
