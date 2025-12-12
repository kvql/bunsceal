package application

import (
	"fmt"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/validation"
)

// SegL1Service orchestrates loading and validation of SegL1 entities
type SegL1Service struct {
	repository SegL1Repository
}

// NewSegL1Service creates a new SegL1 service with the provided repository
func NewSegL1Service(repository SegL1Repository) *SegL1Service {
	return &SegL1Service{
		repository: repository,
	}
}

// Load loads SegL1 entities from the source, validates schema, ID uniqueness and performs postLoad, them,
// and returns a map indexed by ID
func (s *SegL1Service) Load(source string) (map[string]domain.Seg, error) {
	// Load entities from repository
	segList, err := s.repository.LoadAll(source)
	if err != nil {
		return nil, err
	}

	// Validate uniqueness of ID field before forming map
	validations := validation.IdentifierUniquenessValidation(segList)
	if len(validations) > 0 {
		for _, result := range validations {
			o11y.Log.Println(result)
		}
		return nil, fmt.Errorf("security environment validation failed, source: %s", source)
	}

	// Build map from validated list
	segL1Map := make(map[string]domain.Seg)
	for _, seg := range segList {
		segL1Map[seg.ID] = seg
	}

	return segL1Map, nil
}

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

// Load loads Seg entities from the source, validates them,
// and returns a map indexed by ID
func (s *SegService) Load(source string) (map[string]domain.Seg, error) {
	// Load entities from repository
	segList, err := s.repository.LoadAll(source)
	if err != nil {
		return nil, err
	}

	// Validate uniqueness
	validations := validation.IdentifierUniquenessValidation(segList)
	if len(validations) > 0 {
		for _, result := range validations {
			o11y.Log.Println(result)
		}
		return nil, fmt.Errorf("security domain validation failed, source: %s", source)
	}

	// Build map from validated list
	SegMap := make(map[string]domain.Seg)
	for _, seg := range segList {
		SegMap[seg.ID] = seg
	}

	return SegMap, nil
}
