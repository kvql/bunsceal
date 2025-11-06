package domain

import (
	"errors"
)

type Identifiers struct {
	Name string
	ID   string
}

type UnqSegKeys interface {
	GetIdentities() Identifiers
}

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

type L1Overrides struct {
	Sensitivity          string             `yaml:"sensitivity"`
	SensitivityRationale string             `yaml:"sensitivity_rationale"`
	Criticality          string             `yaml:"criticality"`
	CriticalityRationale string             `yaml:"criticality_rationale"`
	ComplianceReqs       []string           `yaml:"compliance_reqs"`
	CompReqs             map[string]CompReq `yaml:"comp_reqs,omitempty"`
}

type SegL2 struct {
	Name        string                 `yaml:"name"`
	ID          string                 `yaml:"id"`
	Description string                 `yaml:"description"`
	L1Overrides map[string]L1Overrides `yaml:"l1_overrides"`
}

func (s SegL2) GetIdentities() Identifiers { return Identifiers{Name: s.Name, ID: s.ID} }

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}

type Taxonomy struct {
	ApiVersion        string
	SegL1s            map[string]SegL1
	SegL2s            map[string]SegL2
	SensitivityLevels []string
	CriticalityLevels []string
	CompReqs          map[string]CompReq
	Config            Config
}

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
