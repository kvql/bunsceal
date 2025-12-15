package application

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/validation"
)

// SegService orchestrates loading and validation of Seg entities
type SegService struct {
	repository SegRepository
}

// NewSegService creates a new Seg service with the provided repository
func NewSegService(repository SegRepository) *SegService {
	return &SegService{
		repository: repository,
	}
}

// LoadLevel loads Seg entities from the source, validates them,
// and returns a map indexed by ID
func (s *SegService) LoadLevel(level string) (map[string]domain.Seg, error) {
	// Load entities from repository
	segList, err := s.repository.LoadLevel(level)
	if err != nil {
		return nil, err
	}

	// Validate uniqueness
	validations := validation.IdentifierUniquenessValidation(segList)
	if len(validations) > 0 {
		for _, result := range validations {
			o11y.Log.Println(result)
		}
		return nil, fmt.Errorf("security domain validation failed, source: %s", level)
	}

	// Build map from validated list
	SegMap := make(map[string]domain.Seg)
	for _, seg := range segList {
		SegMap[seg.ID] = seg
	}

	return SegMap, nil
}
