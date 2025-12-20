package domain

import (
	"errors"
	"fmt"
	"strings"
)

// Seg represents a Level 1 segment (Environment).
type Seg struct {
	Name            string                       `yaml:"name" json:"name"`
	ID              string                       `yaml:"id" json:"id"`
	Description     string                       `yaml:"description" json:"description"`
	Level           string                       `yaml:"level,omitempty" json:"level,omitempty"`
	ComplianceReqs  []string                     `yaml:"compliance_reqs,omitempty" json:"compliance_reqs,omitempty"`
	L1Parents       []string                     `yaml:"l1_parents,omitempty" json:"l1_parents,omitempty"`
	L1Overrides     map[string]L1Overrides       `yaml:"l1_overrides,omitempty" json:"l1_overrides,omitempty"`
	Prominence      int                          `yaml:"prominence,omitempty" json:"prominence,omitempty"`
	Labels          []string                     `yaml:"labels" json:"labels,omitempty"`
	ParsedLabels    map[string]string            `yaml:"-" json:"-"`
	LabelNamespaces map[string]map[string]string `yaml:"-" json:"-"`
}

type L1Overrides struct {
	ComplianceReqs  []string                     `yaml:"compliance_reqs" json:"compliance_reqs,omitempty"`
	CompReqs        map[string]CompReq           `yaml:"comp_reqs,omitempty" json:"comp_reqs,omitempty"`
	Labels          []string                     `yaml:"labels,omitempty" json:"labels,omitempty"`
	ParsedLabels    map[string]string            `yaml:"-" json:"-"`
	LabelNamespaces map[string]map[string]string `yaml:"-" json:"-"`
}

func (o *L1Overrides) ParseLabels() error {
	return parseLabelsIntoMaps(o.Labels, &o.ParsedLabels, &o.LabelNamespaces)
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

func (s Seg) GetNamespacedValue(parent string, ns string, key string) (string, error) {
	if _, hasOverride := s.L1Overrides[parent]; hasOverride {
		if _, hasNs := s.L1Overrides[parent].LabelNamespaces[ns]; hasNs {
			if _, haskey := s.L1Overrides[parent].LabelNamespaces[ns][key]; haskey {
				return s.L1Overrides[parent].LabelNamespaces[ns][key], nil
			}
		}
	}
	if _, hasNs := s.LabelNamespaces[ns]; hasNs {
		if _, haskey := s.LabelNamespaces[ns][key]; haskey {
			return s.LabelNamespaces[ns][key], nil
		}
	}
	return "", fmt.Errorf("no value found for ns(%s), key(%s) on seg(%s)", ns, key, s.ID)
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
		// L1 segments: classification now handled via plugin labels
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

	// Parse labels (handles both segment and override labels)
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

// parseLabelsIntoMaps parses label strings into both ParsedLabels and LabelNamespaces maps.
// Format: "namespace/key:value" -> ParsedLabels["namespace/key"] = "value"
//
//	-> LabelNamespaces["namespace"]["key"] = "value"
func parseLabelsIntoMaps(labels []string, parsed *map[string]string, namespaces *map[string]map[string]string) error {
	*parsed = make(map[string]string)
	*namespaces = make(map[string]map[string]string)

	for _, label := range labels {
		parts := strings.SplitN(label, ":", 2)
		if len(parts) != 2 {
			return fmt.Errorf("invalid label format: %s", label)
		}

		key := parts[0]
		value := parts[1]
		(*parsed)[key] = value

		// Group by namespace
		nsParts := strings.SplitN(key, "/", 2)
		if len(nsParts) == 2 {
			ns := nsParts[0]
			nsKey := nsParts[1]
			if (*namespaces)[ns] == nil {
				(*namespaces)[ns] = make(map[string]string)
			}
			(*namespaces)[ns][nsKey] = value
		}
	}

	return nil
}

// ParseLabels parses both segment labels AND L1 override labels.
// Expects format "key:value" where keys follow DNS-like naming (alphanumeric + ./_-)
// and values support AWS-compliant tag characters (alphanumeric + ./_-+=:@ and spaces).
// Values may contain colons (e.g., "url:https://example.com:8080").
// The function also extracts and groups namespaced labels.
// Example format: bunsceal.plugin.classification/key
func (s *Seg) ParseLabels() error {
	// Parse segment labels
	if err := parseLabelsIntoMaps(s.Labels, &s.ParsedLabels, &s.LabelNamespaces); err != nil {
		return fmt.Errorf("segment %s: %w", s.ID, err)
	}

	// Parse override labels
	for parentID, override := range s.L1Overrides {
		if err := parseLabelsIntoMaps(override.Labels, &override.ParsedLabels, &override.LabelNamespaces); err != nil {
			return fmt.Errorf("segment %s l1_override[%s]: %w", s.ID, parentID, err)
		}
		s.L1Overrides[parentID] = override
	}

	return nil
}
