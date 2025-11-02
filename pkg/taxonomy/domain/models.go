package domain

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
