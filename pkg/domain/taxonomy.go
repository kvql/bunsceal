package domain


// Taxonomy is the root aggregate containing all taxonomy data.
type Taxonomy struct {
	ApiVersion        string
	SegL1s            map[string]SegL1
	SegL2s            map[string]SegL2
	SensitivityLevels []string
	CriticalityLevels []string
	CompReqs          map[string]CompReq
}

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}

