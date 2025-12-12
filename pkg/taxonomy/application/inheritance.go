package application

import (
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/validation"
)

// ApplyInheritance applies inheritance rules for taxonomy segments and validates cross-entity references.
// Returns true if all validations pass, false otherwise.
func ApplyInheritance(txy *domain.Taxonomy) bool {
	// Validate L1 definitions reference valid compliance requirements
	valid := validation.ValidateL1Comp(txy)
	if !valid {
		return valid
	}

	// Loop through env details for each security domain and update risk compliance if not set based on env default
	for _, Seg := range txy.Segs {
		// Initialize L1Overrides map if nil (enables parent-without-override pattern)
		if Seg.L1Overrides == nil {
			Seg.L1Overrides = make(map[string]domain.L1Overrides)
		}

		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range Seg.L1Parents {
			// Get existing override or create empty struct for full inheritance
			l1Override, exists := Seg.L1Overrides[l1ID]
			if !exists {
				l1Override = domain.L1Overrides{}
			}

			if l1Override.Sensitivity == "" && l1Override.SensitivityRationale == "" {
				l1Override.Sensitivity = txy.SegL1s[l1ID].Sensitivity
				l1Override.SensitivityRationale = "Inherited: " + txy.SegL1s[l1ID].SensitivityRationale
			}
			if l1Override.Criticality == "" && l1Override.CriticalityRationale == "" {
				l1Override.Criticality = txy.SegL1s[l1ID].Criticality
				l1Override.CriticalityRationale = "Inherited: " + txy.SegL1s[l1ID].CriticalityRationale
			}
			// Inherit compliance requirements from environment if not set
			if l1Override.ComplianceReqs == nil {
				l1Override.ComplianceReqs = txy.SegL1s[l1ID].ComplianceReqs
			}
			// add compliance details to compReqs var for each compliance standard listed
			for _, compReq := range l1Override.ComplianceReqs {
				// only add details if listed standard is valid. If not, it will be caught in validation
				if _, ok := txy.CompReqs[compReq]; ok {
					if l1Override.CompReqs == nil {
						l1Override.CompReqs = make(map[string]domain.CompReq)
					}
					l1Override.CompReqs[compReq] = txy.CompReqs[compReq]
				}
			}
			// Write back override (creates new entry if didn't exist)
			Seg.L1Overrides[l1ID] = l1Override
		}
	}

	// Validate L2 definitions after inheritance
	valid, _ = validation.ValidateL2Definition(txy)
	return valid
}
