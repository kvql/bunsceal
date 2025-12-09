package taxonomy

import "github.com/kvql/bunsceal/pkg/domain"

// SegL1Repository defines the contract for loading SegL1 data from any source
type SegL1Repository interface {
	// LoadAll loads all SegL1 entities from the specified source
	// Returns a slice of SegL1 entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.SegL1, error)
}

// SegL2Repository defines the contract for loading SegL2 data from any source
type SegL2Repository interface {
	// LoadAll loads all SegL2 entities from the specified source
	// Returns a slice of SegL2 entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.SegL2, error)
}
