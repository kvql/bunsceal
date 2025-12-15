package application

import "github.com/kvql/bunsceal/pkg/domain"

// SegRepository defines the contract for loading Seg data from any source
type SegRepository interface {
	// LoadAll loads all Seg entities from the specified source
	// Returns a slice of Seg entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadLevel(level string) ([]domain.Seg, error)
}
