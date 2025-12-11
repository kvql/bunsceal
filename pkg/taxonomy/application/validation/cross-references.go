package validation

import (
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
)

func ValidateL2Definition(txy *domain.Taxonomy) (bool, int) {
	valid := true
	failures := 0
	// Loop through SegL2s and validate default risk levels
	for _, secDomain := range txy.SegL2s {
		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range secDomain.L1Parents {
			// Validate parent L1 exists in taxonomy
			if _, ok := txy.SegL1s[l1ID]; !ok {
				o11y.Log.Printf("Invalid L1 parent for SegL2 %s: %s\n", secDomain.Name, l1ID)
				failures++
				valid = false
				continue // Skip compliance validation for invalid parent
			}

			// Lookup override (should exist after inheritance)
			sdEnv, exists := secDomain.L1Overrides[l1ID]
			if !exists {
				// This should not happen if inheritance ran correctly
				o11y.Log.Printf("ERROR: SegL2 '%s' has parent '%s' but no override data after inheritance\n", secDomain.Name, l1ID)
				failures++
				valid = false
				continue
			}

			// Validate compliance scope for each SD
			for _, compReq := range sdEnv.ComplianceReqs {
				if _, ok := txy.CompReqs[compReq]; !ok {
					o11y.Log.Printf("Invalid compliance scope(%s) for SegL1(%s) in SD (%s)", compReq, l1ID, secDomain.Name)
					failures++
					valid = false
				}
			}
		}
	}
	return valid, failures
}

func ValidateL1Comp(txy *domain.Taxonomy) bool {
	valid := true
	// Loop through environments and validate
	for _, env := range txy.SegL1s {
		// validate compliance scope for each SD
		for _, compReq := range env.ComplianceReqs {
			if _, ok := txy.CompReqs[compReq]; !ok {
				o11y.Log.Printf("Invalid compliance scope (%s) for env(%s)", env.ID, compReq)
				valid = false
			}
		}
	}
	return valid
}
