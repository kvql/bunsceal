package application

import "github.com/kvql/bunsceal/pkg/domain"

// SegL1Repository defines the contract for loading SegL1 data from any source
type SegL1Repository interface {
	// LoadAll loads all SegL1 entities from the specified source
	// Returns a slice of SegL1 entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.Seg, error)
}

// SegRepository defines the contract for loading Seg data from any source
type SegRepository interface {
	// LoadAll loads all Seg entities from the specified source
	// Returns a slice of Seg entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.Seg, error)
}
