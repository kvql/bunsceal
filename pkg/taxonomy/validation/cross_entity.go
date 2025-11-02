package validation

import (
	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/util"
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
					util.Log.Printf("Invalid compliance scope(%s) for SegL1(%s) in SD (%s)", compReq, envID, secDomain.Name)
					failures++
					valid = false
				}
			}
			// validate secEnv for each SD is a valid secEnv in the taxonomy,
			if _, ok := txy.SegL1s[envID]; !ok {
				util.Log.Printf("Invalid secEnv for SD %s: %s\n", secDomain.Name, envID)
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
				util.Log.Printf("Invalid compliance scope (%s) for env(%s)", env.ID, compReq)
				valid = false
			}
		}
	}
	return valid
}

// ValidateSharedServices validates the shared-services environment.
// Hard coding some specific checks for shared-services environment as it is an exception to the normal security domains.
// This secenv need to meet the strictest requirements.
func ValidateSharedServices(txy *domain.Taxonomy) (bool, int) {
	valid := true
	envName := "shared-service"
	failures := 0
	if _, ok := txy.SegL1s[envName]; !ok {
		util.Log.Printf("%s environment not found", envName)
		return false, 1
	}
	if txy.SegL1s[envName].Sensitivity != domain.SenseOrder[0] ||
		txy.SegL1s[envName].Criticality != domain.CritOrder[0] {
		util.Log.Printf("%s environment does not have the highest sensitivity or criticality", envName)
		failures++
		valid = false
	}

	if len(txy.SegL1s[envName].ComplianceReqs) != len(txy.CompReqs) {
		util.Log.Printf("%s environment does not have all compliance requirements", envName)
		failures++
		valid = false
	}
	return valid, failures
}
