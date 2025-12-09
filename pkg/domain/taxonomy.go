package domain

import (
	"errors"
)

// Taxonomy is the root aggregate containing all taxonomy data.
type Taxonomy struct {
	ApiVersion        string
	SegL1s            map[string]SegL1
	SegL2s            map[string]SegL2
	SensitivityLevels []string
	CriticalityLevels []string
	CompReqs          map[string]CompReq
}

// SegL1 represents a Level 1 segment (Environment).
type SegL1 struct {
	Name                 string   `yaml:"name"`
	ID                   string   `yaml:"id"`
	Description          string   `yaml:"description"`
	Sensitivity          string   `yaml:"sensitivity"`
	SensitivityRationale string   `yaml:"sensitivity_rationale"`
	Criticality          string   `yaml:"criticality"`
	CriticalityRationale string   `yaml:"criticality_rationale"`
	ComplianceReqs       []string `yaml:"compliance_reqs"`
}

func (s SegL1) GetIdentities() Identifiers { return Identifiers{Name: s.Name, ID: s.ID} }

func (s SegL1) GetKeyString(key string) (string, error) {
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

// SegL2 represents a Level 2 segment.
type SegL2 struct {
	Name        string                 `yaml:"name"`
	ID          string                 `yaml:"id"`
	Description string                 `yaml:"description"`
	L1Overrides map[string]L1Overrides `yaml:"l1_overrides"`
	Prominence  int                    `yaml:"prominence"`
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

type L1Overrides struct {
	Sensitivity          string             `yaml:"sensitivity"`
	SensitivityRationale string             `yaml:"sensitivity_rationale"`
	Criticality          string             `yaml:"criticality"`
	CriticalityRationale string             `yaml:"criticality_rationale"`
	ComplianceReqs       []string           `yaml:"compliance_reqs"`
	CompReqs             map[string]CompReq `yaml:"comp_reqs,omitempty"`
}

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}
