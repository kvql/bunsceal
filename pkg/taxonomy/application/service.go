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

// LoadAndValidate loads SegL1 entities from the source, validates them,
// and returns a map indexed by ID
func (s *SegL1Service) LoadAndValidate(source string) (map[string]domain.SegL1, error) {
	// Load entities from repository
	segList, err := s.repository.LoadAll(source)
	if err != nil {
		return nil, err
	}

	// Validate uniqueness
	validations := validation.UniquenessValidator(segList)
	if len(validations) > 0 {
		for _, result := range validations {
			o11y.Log.Println(result)
		}
		return nil, fmt.Errorf("security environment validation failed, source: %s", source)
	}

	// Build map from validated list
	segL1Map := make(map[string]domain.SegL1)
	for _, seg := range segList {
		segL1Map[seg.ID] = seg
	}

	return segL1Map, nil
}

// SegL2Service orchestrates loading and validation of SegL2 entities
type SegL2Service struct {
	repository SegL2Repository
}

// NewSegL2Service creates a new SegL2 service with the provided repository
func NewSegL2Service(repository SegL2Repository) *SegL2Service {
	return &SegL2Service{
		repository: repository,
	}
}

// LoadAndValidate loads SegL2 entities from the source, validates them,
// and returns a map indexed by ID
func (s *SegL2Service) LoadAndValidate(source string) (map[string]domain.SegL2, error) {
	// Load entities from repository
	segList, err := s.repository.LoadAll(source)
	if err != nil {
		return nil, err
	}

	// Validate uniqueness
	validations := validation.UniquenessValidator(segList)
	if len(validations) > 0 {
		for _, result := range validations {
			o11y.Log.Println(result)
		}
		return nil, fmt.Errorf("security domain validation failed, source: %s", source)
	}

	// Build map from validated list
	segL2Map := make(map[string]domain.SegL2)
	for _, seg := range segList {
		segL2Map[seg.ID] = seg
	}

	return segL2Map, nil
}
