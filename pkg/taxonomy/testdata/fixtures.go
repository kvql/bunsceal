package testdata

// Valid SegL1 Fixtures
// These represent correctly structured security environments
// Note: These use the generic types to avoid import cycles
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
	L1Overrides map[string]L1Overrides `yaml:"l1_overrides"`
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

type Config struct {
	Terminology TermConfig `yaml:"terminology"`
}
type InvalidConfig struct {
	Terminology InvalidTermConfig `yaml:"terminology"`
}

// TermConfig holds terminology configuration for L1 and L2 segments
type TermConfig struct {
	L1 TermDef `yaml:"l1,omitempty"`
	L2 TermDef `yaml:"l2,omitempty"`
}
type InvalidTermConfig struct {
	L4 TermDef `yaml:"l4"`
}

// TermDef defines singular and plural forms for a segment level
type TermDef struct {
	Singular string `yaml:"singular"`
	Plural   string `yaml:"plural"`
}

var InvalidConfigSchema = InvalidConfig{
	Terminology: InvalidTermConfig{
		L4: TermDef{
			Singular: "dfas",
			Plural:   "fdasdfas",
		},
	},
}
var ValidConfigSchema = Config{
	Terminology: TermConfig{
		L1: TermDef{
			Singular: "dfas",
			Plural:   "fdasdfas",
		},
	},
}

var ValidSegL1Production = SegL1{
	Name:                 "Production",
	ID:                   "production",
	Description:          "Production environment with highest security controls and strict compliance requirements for customer-facing services.",
	Sensitivity:          "A",
	SensitivityRationale: "Production handles all customer data including personally identifiable information and financial transactions requiring the highest classification level.",
	Criticality:          "1",
	CriticalityRationale: "Production outages directly impact customers and revenue streams, requiring immediate incident response and highest availability standards.",
	ComplianceReqs:       []string{"pci-dss", "sox"},
}

var ValidSegL1Staging = SegL1{
	Name:                 "Staging",
	ID:                   "staging",
	Description:          "Pre-production staging environment used for final testing and validation before production deployment cycles.",
	Sensitivity:          "D",
	SensitivityRationale: "Staging contains no production or customer data, only synthetic test data generated for validation purposes.",
	Criticality:          "5",
	CriticalityRationale: "Staging downtime impacts development velocity but has no direct customer or revenue impact on business operations.",
	ComplianceReqs:       []string{},
}

var ValidSegL1SharedService = SegL1{
	Name:                 "Shared Service",
	ID:                   "shared-service",
	Description:          "Shared service environment hosting cross-account resources and centralized services with network connectivity across all environments.",
	Sensitivity:          "A",
	SensitivityRationale: "Shared services represent the highest risk from lateral movement perspective and could bridge between low and high security environments.",
	Criticality:          "1",
	CriticalityRationale: "All environments depend on shared services for core functionality and operations, making outages highly impactful across entire infrastructure.",
	ComplianceReqs:       []string{"pci-dss", "sox", "hipaa"},
}

// Invalid SegL1 Fixtures for negative testing

var InvalidSegL1_MissingName = SegL1{
	ID:                   "test",
	Description:          "This fixture is missing the required 'name' field which should cause schema validation to fail.",
	Sensitivity:          "A",
	SensitivityRationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
	Criticality:          "1",
	CriticalityRationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
}

var InvalidSegL1_InvalidID = SegL1{
	Name:                 "Invalid ID",
	ID:                   "Invalid_ID_With_Capitals",
	Description:          "This fixture has an ID with uppercase letters and underscores which violates the pattern constraint.",
	Sensitivity:          "A",
	SensitivityRationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
	Criticality:          "1",
	CriticalityRationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
}

var InvalidSegL1_ShortDescription = SegL1{
	Name:                 "Short Desc",
	ID:                   "short",
	Description:          "Too short",
	Sensitivity:          "A",
	SensitivityRationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
	Criticality:          "1",
	CriticalityRationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
}

var InvalidSegL1_InvalidSensitivity = SegL1{
	Name:                 "Invalid Sensitivity",
	ID:                   "invalid-sens",
	Description:          "This fixture has an invalid sensitivity value that is not in the allowed enum of A, B, C, D values.",
	Sensitivity:          "Z",
	SensitivityRationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
	Criticality:          "1",
	CriticalityRationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
}

var InvalidSegL1_InvalidCriticality = SegL1{
	Name:                 "Invalid Criticality",
	ID:                   "invalid-crit",
	Description:          "This fixture has an invalid criticality value that is not in the allowed enum of 1, 2, 3, 4, 5 values.",
	Sensitivity:          "A",
	SensitivityRationale: "Test sensitivity rationale with sufficient length to meet the minimum character requirement for descriptions.",
	Criticality:          "9",
	CriticalityRationale: "Test criticality rationale with sufficient length to meet the minimum character requirement for descriptions.",
}

var InvalidSegL1_ShortRationale = SegL1{
	Name:                 "Short Rationale",
	ID:                   "short-rat",
	Description:          "This fixture has rationale fields that are too short and don't meet the minimum length requirement.",
	Sensitivity:          "A",
	SensitivityRationale: "Too short",
	Criticality:          "1",
	CriticalityRationale: "Also too short",
}

// Valid SegL2 Fixtures

var ValidSegL2Security = SegL2{
	Name:        "Security",
	ID:          "sec",
	Description: "Security domain for security tooling and monitoring infrastructure",
	L1Overrides: map[string]L1Overrides{
		"production": {
			Sensitivity:          "B",
			SensitivityRationale: "Security logs contain sensitive metadata but not direct customer PII or financial data requiring lower classification.",
			Criticality:          "2",
			CriticalityRationale: "Security monitoring downtime impacts compliance and incident response but doesn't directly affect customer-facing services.",
			ComplianceReqs:       []string{"sox"},
		},
		"staging": {
			Sensitivity:          "D",
			SensitivityRationale: "Staging security environment uses only test data and synthetic logs with no production data present.",
			Criticality:          "5",
			CriticalityRationale: "Staging environment downtime only impacts testing cycles with no customer or business impact.",
			ComplianceReqs:       []string{},
		},
	},
}

var ValidSegL2Application = SegL2{
	Name:        "Application",
	ID:          "app",
	Description: "Application domain for core business application services",
	L1Overrides: map[string]L1Overrides{
		"production": {
			Sensitivity:          "A",
			SensitivityRationale: "Production applications handle customer PII, payment information, and other regulated data requiring highest protection.",
			Criticality:          "1",
			CriticalityRationale: "Application services are customer-facing and directly generate revenue, requiring maximum availability and performance.",
			ComplianceReqs:       []string{"pci-dss", "sox"},
		},
	},
}

// SegL2 with inheritance (empty env details)
var ValidSegL2WithInheritance = SegL2{
	Name:        "Infrastructure",
	ID:          "infra",
	Description: "Infrastructure domain that inherits all settings from environments",
	L1Overrides: map[string]L1Overrides{
		"production": {
			// Empty - should inherit from production SegL1
		},
		"staging": {
			// Empty - should inherit from staging SegL1
		},
	},
}

// Invalid SegL2 Fixtures

var InvalidSegL2_MissingName = SegL2{
	ID:          "invalid",
	Description: "Missing name field",
	L1Overrides: map[string]L1Overrides{
		"production": {},
	},
}

var InvalidSegL2_InvalidID = SegL2{
	Name:        "Invalid ID",
	ID:          "Invalid_ID!",
	Description: "ID contains invalid characters",
	L1Overrides: map[string]L1Overrides{
		"production": {},
	},
}

var InvalidSegL2_NoL1Overrides = SegL2{
	Name:        "No Environments",
	ID:          "no-env",
	Description: "Security domain with no environment details defined",
	L1Overrides: map[string]L1Overrides{},
}

// Valid CompReq Fixtures

var ValidCompReqs = map[string]CompReq{
	"pci-dss": {
		Name:        "PCI DSS",
		Description: "Payment Card Industry Data Security Standard",
		ReqsLink:    "https://www.pcisecuritystandards.org/",
	},
	"sox": {
		Name:        "SOX",
		Description: "Sarbanes-Oxley Act compliance requirements",
		ReqsLink:    "https://www.sox-online.com/",
	},
	"hipaa": {
		Name:        "HIPAA",
		Description: "Health Insurance Portability and Accountability Act",
		ReqsLink:    "https://www.hhs.gov/hipaa/",
	},
}

// Complete Valid Taxonomy Fixture

var ValidCompleteTaxonomy = Taxonomy{
	ApiVersion: "v1beta1",
	SegL1s: map[string]SegL1{
		"production":     ValidSegL1Production,
		"staging":        ValidSegL1Staging,
		"shared-service": ValidSegL1SharedService,
	},
	SegL2s: map[string]SegL2{
		"sec": ValidSegL2Security,
		"app": ValidSegL2Application,
	},
	SensitivityLevels: []string{"A", "B", "C", "D"},
	CriticalityLevels: []string{"1", "2", "3", "4", "5"},
	CompReqs:          ValidCompReqs,
}
