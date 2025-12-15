package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Seg represents a Level 1 segment (Environment).
type Seg struct {
	Name                 string                       `yaml:"name" json:"name"`
	ID                   string                       `yaml:"id" json:"id"`
	Description          string                       `yaml:"description" json:"description"`
	Level                string                       `yaml:"level,omitempty" json:"level,omitempty"`
	Sensitivity          string                       `yaml:"sensitivity,omitempty" json:"sensitivity,omitempty"`                     //Depricated - migrating to label
	SensitivityRationale string                       `yaml:"sensitivity_rationale,omitempty" json:"sensitivity_rationale,omitempty"` //Depricated - migrating to label
	Criticality          string                       `yaml:"criticality,omitempty" json:"criticality,omitempty"`                     //Depricated - migrating to label
	CriticalityRationale string                       `yaml:"criticality_rationale,omitempty" json:"criticality_rationale,omitempty"` //Depricated - migrating to label
	ComplianceReqs       []string                     `yaml:"compliance_reqs,omitempty" json:"compliance_reqs,omitempty"`             //Depricated - migrating to label
	L1Parents            []string                     `yaml:"l1_parents,omitempty" json:"l1_parents,omitempty"`
	L1Overrides          map[string]L1Overrides       `yaml:"l1_overrides,omitempty" json:"l1_overrides,omitempty"`
	Prominence           int                          `yaml:"prominence,omitempty" json:"prominence,omitempty"`
	Labels               []string                     `yaml:"labels" json:"labels,omitempty"`
	ParsedLabels         map[string]string            `yaml:"-" json:"-"`
	LabelNamespaces      map[string]map[string]string `yaml:"-" json:"-"`
}

type L1Overrides struct {
	Sensitivity          string             `yaml:"sensitivity" json:"sensitivity,omitempty"`
	SensitivityRationale string             `yaml:"sensitivity_rationale" json:"sensitivity_rationale,omitempty"`
	Criticality          string             `yaml:"criticality" json:"criticality,omitempty"`
	CriticalityRationale string             `yaml:"criticality_rationale" json:"criticality_rationale,omitempty"`
	ComplianceReqs       []string           `yaml:"compliance_reqs" json:"compliance_reqs,omitempty"`
	CompReqs             map[string]CompReq `yaml:"comp_reqs,omitempty" json:"comp_reqs,omitempty"`
}

// ###################
// Segment Methods

func (s Seg) GetIdentities() Identifiers { return Identifiers{Name: s.Name, ID: s.ID} }

func (s Seg) GetKeyString(key string) (string, error) {
	var val string
	switch key {
	case "name":
		val = s.Name
	case "description":
		val = s.Description
	default:
		return "", errors.New("SegL1.GetKeyValue() unsupported key")
	}
	return val, nil
}

func (s *Seg) PostLoad(level string) error {
	// Set level if not already set
	if s.Level == "" {
		s.Level = level
	} else if s.Level != level {
		return fmt.Errorf("level (%s) passed as argument, doesn't match level field (%s)", level, s.Level)
	}

	// Validate level-specific required fields
	switch level {
	case "1":
		// L1 segments require top-level sensitivity/criticality
		if s.Criticality == "" {
			return fmt.Errorf("L1 segment missing required field: Criticality")
		}
		if s.CriticalityRationale == "" {
			return fmt.Errorf("L1 segment missing required field: CriticalityRationale")
		}
		if s.Sensitivity == "" {
			return fmt.Errorf("L1 segment missing required field: Sensitivity")
		}
		if s.SensitivityRationale == "" {
			return fmt.Errorf("L1 segment missing required field: SensitivityRationale")
		}
	case "2":
		// L2 segments require L1Parents
		if len(s.L1Parents) == 0 {
			return fmt.Errorf("L2 segment missing required field: L1Parents")
		}
		// Validate L1Overrides keys match L1Parents
		if err := s.ValidateL1Consistency(); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported segment level: %s", level)
	}

	// Apply defaults
	s.SetDefaults()

	// Parse labels into map
	return s.ParseLabels()
}

func (s *Seg) SetDefaults() {
	if s.Prominence == 0 {
		s.Prominence = 1
	}
}

// ValidateL1Consistency ensures that all L1Overrides keys are present in L1Parents.
// This prevents invalid YAML where overrides reference parents not in the parent list.
// Returns error if any override key is not in L1Parents.
func (s *Seg) ValidateL1Consistency() error {
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

// ParseLabels converts labels string array into ParsedLabels map.
// Expects format "key:value" where keys follow DNS-like naming (alphanumeric + ./_-)
// and values support AWS-compliant tag characters (alphanumeric + ./_-+=:@ and spaces).
// Values may contain colons (e.g., "url:https://example.com:8080").
// The function also extracts and group namespaced labels.
// example formot bunsceal.plugin.classification/key
func (s *Seg) ParseLabels() error {
	if s.ParsedLabels == nil {
		s.ParsedLabels = make(map[string]string)
	}
	if s.LabelNamespaces == nil {
		s.LabelNamespaces = make(map[string]map[string]string)
	}
	for _, label := range s.Labels {
		k := strings.SplitN(label, ":", 2)
		if len(k) == 2 {
			s.ParsedLabels[k[0]] = k[1]
		} else {
			return fmt.Errorf("label format invalid, expected key:value")
		}
		nsRaw := strings.SplitN(k[0], "/", 2)
		if len(nsRaw) == 2 {
			if _, exists := s.LabelNamespaces[nsRaw[0]]; exists {
				s.LabelNamespaces[nsRaw[0]][nsRaw[1]] = k[1]
			} else {
				s.LabelNamespaces[nsRaw[0]] = make(map[string]string)
				s.LabelNamespaces[nsRaw[0]][nsRaw[1]] = k[1]
			}
		}
	}
	return nil
}
