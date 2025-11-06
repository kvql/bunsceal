package validation

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
)

// UniquenessValidator validates that IDs and Names are unique across a collection of objects
// Returns a slice of error messages if validation fails, empty slice if all validations pass
func UniquenessValidator[T domain.UnqSegKeys](objects []T) []string {
	idMap := make(map[string]bool)
	var validations []string

	for _, item := range objects {
		identifiers := item.GetIdentities()

		// Validate ID is unique
		if _, exists := idMap[identifiers.ID]; exists {
			validations = append(validations, fmt.Sprintf("ID for %s is not unique: %s", identifiers.Name, identifiers.ID))
		} else {
			idMap[identifiers.ID] = true
		}
	}

	return validations
}
