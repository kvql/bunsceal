package application

import (
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
)

// ApplyInheritance applies inheritance rules for taxonomy segments and validates cross-entity references.
// Pass nil for pluginsList to skip plugin label inheritance (backwards compatible).
func ApplyInheritance(txy *domain.Taxonomy, pluginsList plugins.Plugins) error {
	for _, seg := range txy.SegsL2s {
		// Initialize L1Overrides map if nil (enables parent-without-override pattern)
		if seg.L1Overrides == nil {
			seg.L1Overrides = make(map[string]domain.L1Overrides)
		}

		// REFACTORED: Iterate over L1Parents instead of L1Overrides keys
		for _, l1ID := range seg.L1Parents {
			// Get existing override or create empty struct for full inheritance
			l1Override, exists := seg.L1Overrides[l1ID]
			if !exists {
				l1Override = domain.L1Overrides{}
			}

			// Write back override (creates new entry if didn't exist)
			seg.L1Overrides[l1ID] = l1Override

			// Plugin labels are inherited via plugin system
			if pluginsList != nil {
				errs := pluginsList.ApplyPluginInheritanceAndValidate(
					txy.SegL1s[l1ID],
					&seg,
				)
				for _, err := range errs {
					o11y.Log.Println(err)
				}
				if len(errs) > 0 {
					return errs[0]
				}
			}
		}
	}
	return nil
}
