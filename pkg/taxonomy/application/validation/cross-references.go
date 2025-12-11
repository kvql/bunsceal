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
		for envID, sdEnv := range secDomain.L1Overrides {
			// validate compliance scope for each SD
			for _, compReq := range sdEnv.ComplianceReqs {
				if _, ok := txy.CompReqs[compReq]; !ok {
					o11y.Log.Printf("Invalid compliance scope(%s) for SegL1(%s) in SD (%s)", compReq, envID, secDomain.Name)
					failures++
					valid = false
				}
			}
			// validate secEnv for each SD is a valid secEnv in the taxonomy,
			if _, ok := txy.SegL1s[envID]; !ok {
				o11y.Log.Printf("Invalid secEnv for SD %s: %s\n", secDomain.Name, envID)
				failures++
				valid = false
			}
		}
	}
	return valid, failures
}

func ValidateL1Definitions(txy *domain.Taxonomy) bool {
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
