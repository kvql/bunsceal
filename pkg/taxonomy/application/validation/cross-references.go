package validation

import (
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

func ValidateL2Definition(txy *domain.Taxonomy, pluginMap plugins.Plugins) (bool, int) {
	valid := true
	failures := 0

	// Loop through Segs and validate L1 parent references
	for _, secDomain := range txy.SegsL2s {
		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range secDomain.L1Parents {
			// Validate parent L1 exists in taxonomy
			if _, ok := txy.SegL1s[l1ID]; !ok {
				o11y.Log.Printf("Invalid L1 parent for Seg %s: %s\n", secDomain.Name, l1ID)
				failures++
				valid = false
				continue
			}

			// Lookup override (should exist after inheritance)
			_, exists := secDomain.L1Overrides[l1ID]
			if !exists {
				// This should not happen if inheritance ran correctly
				o11y.Log.Printf("ERROR: Seg '%s' has parent '%s' but no override data after inheritance\n", secDomain.Name, l1ID)
				failures++
				valid = false
				continue
			}

			// Note: Compliance validation now happens via plugin label validation
		}
	}
	return valid, failures
}
