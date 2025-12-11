package domain

import (
	"errors"
	"fmt"
	"strings"
)

// SegL2 represents a Level 2 segment.
type SegL2 struct {
	Name         string                 `yaml:"name"`
	ID           string                 `yaml:"id"`
	Description  string                 `yaml:"description"`
	L1Parents    []string               `yaml:"l1_parents"`
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

// ValidateL1Consistency ensures that all L1Overrides keys are present in L1Parents.
// This prevents invalid YAML where overrides reference parents not in the parent list.
// Returns error if any override key is not in L1Parents.
func (s *SegL2) ValidateL1Consistency() error {
	if len(s.L1Overrides) == 0 {
		return nil
	}

	// Build set of L1Parents for O(1) lookup
	parentSet := make(map[string]bool, len(s.L1Parents))
	for _, parent := range s.L1Parents {
		parentSet[parent] = true
	}

	// Check each override key exists in parent list
	for overrideKey := range s.L1Overrides {
		if !parentSet[overrideKey] {
			return fmt.Errorf("l1_overrides contains key '%s' which is not in l1_parents", overrideKey)
		}
	}

	return nil
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

type L1Overrides struct {
	Sensitivity          string             `yaml:"sensitivity"`
	SensitivityRationale string             `yaml:"sensitivity_rationale"`
	Criticality          string             `yaml:"criticality"`
	CriticalityRationale string             `yaml:"criticality_rationale"`
	ComplianceReqs       []string           `yaml:"compliance_reqs"`
	CompReqs             map[string]CompReq `yaml:"comp_reqs,omitempty"`
}
