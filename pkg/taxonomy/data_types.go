package taxonomy

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
	L1Overrides  map[string]L1Overrides `yaml:"l1_overrides"`
}

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}

type Taxonomy struct {
	ApiVersion        string             `yaml:"api_version"`
	SegL1s            map[string]SegL1   `yaml:"seg_l1s"`
	SegL2s            map[string]SegL2   `yaml:"seg_l2s"`
	SensitivityLevels []string           `yaml:"sensitivity_levels"`
	CriticalityLevels []string           `yaml:"criticality_levels"`
	CompReqs          map[string]CompReq `yaml:"comp_reqs"`
}