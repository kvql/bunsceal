package domain

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

// SegL2 represents a Level 2 segment.
type SegL2 struct {
	Name         string                 `yaml:"name"`
	ID           string                 `yaml:"id"`
	Description  string                 `yaml:"description"`
	L1Parents    []string               `yaml:"l1_parents,omitempty"`
	L1Overrides  map[string]L1Overrides `yaml:"l1_overrides"`
	Prominence   int                    `yaml:"prominence"`
	Labels       []string               `yaml:"labels"`
	ParsedLabels map[string]string
}

func (s SegL2) GetIdentities() Identifiers { return Identifiers{Name: s.Name, ID: s.ID} }

func (s *SegL2) SetDefaults() {
	if s.Prominence == 0 {
		s.Prominence = 1
	}
}

func (s SegL2) GetKeyString(key string) (string, error) {
	var val string
	switch key {
	case "name":
		val = s.Name
	case "description":
		val = s.Description
	default:
		return "", errors.New("SegL2.GetKeyValue() unsupported key")
	}
	return val, nil
}

// ParseLabels converts labels string array into ParsedLabels map.
// Expects format "key:value" where keys follow DNS-like naming (alphanumeric + ./_-)
// and values support AWS-compliant tag characters (alphanumeric + ./_-+=:@ and spaces).
// Values may contain colons (e.g., "url:https://example.com:8080").
func (s *SegL2) ParseLabels() error {
	if s.ParsedLabels == nil {
		s.ParsedLabels = make(map[string]string)
	}
	for _, label := range s.Labels {
		k := strings.SplitN(label, ":", 2)
		if len(k) == 2 {
			s.ParsedLabels[k[0]] = k[1]
		} else {
			return fmt.Errorf("label format invalid, expected key:value")
		}
	}
	return nil
}

// UpdateLabels converts ParsedLabels into deliminated strings and updates the Label value.
// TODO: To be run after ParsedLabels are updated via inheritance
func (s *SegL2) UpdateLabels() error {
	return nil
}

// MigrateL1Parents populates L1Parents from L1Overrides keys if L1Parents is empty.
// This provides backward compatibility during the migration from implicit to explicit parent relationships.
// Returns true if migration occurred, false if L1Parents was already populated.
func (s *SegL2) MigrateL1Parents() bool {
	if len(s.L1Parents) > 0 {
		return false // Already migrated
	}

	if len(s.L1Overrides) == 0 {
		return false // Nothing to migrate
	}

	// Extract keys from L1Overrides
	s.L1Parents = make([]string, 0, len(s.L1Overrides))
	for l1ID := range s.L1Overrides {
		s.L1Parents = append(s.L1Parents, l1ID)
	}

	// Sort for deterministic behaviour
	sort.Strings(s.L1Parents)
	return true
}

// ValidateL1Consistency checks that L1Overrides keys are a subset of L1Parents.
// During migration, both fields may be present - this ensures they're consistent.
// Returns error if L1Overrides contains keys not in L1Parents.
func (s *SegL2) ValidateL1Consistency() error {
	if len(s.L1Parents) == 0 {
		return nil // Migration not complete, skip validation
	}

	// Build lookup map for efficient checking
	parentSet := make(map[string]bool, len(s.L1Parents))
	for _, l1ID := range s.L1Parents {
		parentSet[l1ID] = true
	}

	// Check all override keys exist in parents
	for l1ID := range s.L1Overrides {
		if !parentSet[l1ID] {
			return fmt.Errorf("L1Override key '%s' not found in L1Parents for SegL2 '%s'", l1ID, s.ID)
		}
	}

	return nil
}

type L1Overrides struct {
	Sensitivity          string             `yaml:"sensitivity"`
	SensitivityRationale string             `yaml:"sensitivity_rationale"`
	Criticality          string             `yaml:"criticality"`
	CriticalityRationale string             `yaml:"criticality_rationale"`
	ComplianceReqs       []string           `yaml:"compliance_reqs"`
	CompReqs             map[string]CompReq `yaml:"comp_reqs,omitempty"`
}
