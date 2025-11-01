package taxonomy

import (
	"errors"
	"strings"

	"github.com/kvql/bunsceal/pkg/util"
)

// Use this to test compatibility in clients
const ApiVersion = "v1beta1"

// Define sensitivity levels
var SensitivityLevels = map[string]string{
	"A": "High",
	"B": "Medium",
	"C": "Low",
	"D": "N/A",
}

// Define Criticality levels
var CriticalityLevels map[string]string = map[string]string{
	"1": "Critical",
	"2": "High",
	"3": "Medium",
	"4": "Low",
	"5": "N/A",
}

// create slice of risk levels to track order
var SenseOrder = []string{"A", "B", "C", "D"}
var CritOrder = []string{"1", "2", "3", "4", "5"}

func (txy *Taxonomy) validateSecurityDomains() (bool, int) {
	valid := true
	failures := 0
	// Loop through SegL2s and validate default risk levels
	for _, secDomain := range txy.SegL2s {
		for envID, sdEnv := range secDomain.L1Overrides {
			// validate compliance scope for each SD
			for _, compReq := range sdEnv.ComplianceReqs {
				if _, ok := txy.CompReqs[compReq]; !ok {
					util.Log.Printf("Invalid compliance scope(%s) for SegL1(%s) in SD (%s)", compReq, envID, secDomain.Name)
					failures++
					valid = false
				}
			}
			// validate secEnv for each SD is a valid secEnv in the taxonomy,
			if _, ok := txy.SegL1s[envID]; !ok {
				util.Log.Printf("Invalid secEnv for SD %s: %s\n", secDomain.Name, envID)
				failures++
				valid = false
			}
		}
	}
	return valid, failures
}

func (txy *Taxonomy) validateEnv() bool {
	valid := true
	// Loop through environments and validate
	for _, env := range txy.SegL1s {
		// validate compliance scope for each SD
		for _, compReq := range env.ComplianceReqs {
			if _, ok := txy.CompReqs[compReq]; !ok {
				util.Log.Printf("Invalid compliance scope (%s) for env(%s)", env.ID, compReq)
				valid = false
			}
		}
	}
	return valid
}

func (txy *Taxonomy) ApplyInheritance() {
	// Loop through env details for each security domain and update risk compliance if not set based on env default
	for _, secDomain := range txy.SegL2s {
		for l1ID, l1Override := range secDomain.L1Overrides {
			if l1Override.Sensitivity == "" && l1Override.SensitivityRationale == "" {
				l1Override.Sensitivity = txy.SegL1s[l1ID].Sensitivity
				l1Override.SensitivityRationale = "Inherited: " + txy.SegL1s[l1ID].SensitivityRationale
			}
			if l1Override.Criticality == "" && l1Override.CriticalityRationale == "" {
				l1Override.Criticality = txy.SegL1s[l1ID].Criticality
				l1Override.CriticalityRationale = "Inherited: " + txy.SegL1s[l1ID].CriticalityRationale
			}
			// Inherit compliance requirements from environment if not set
			if l1Override.ComplianceReqs == nil {
				l1Override.ComplianceReqs = txy.SegL1s[l1ID].ComplianceReqs
			}
			// add compliance details to compReqs var for each complaince standard listed
			for _, compReq := range l1Override.ComplianceReqs {
				// only add details if listed standard is valid. If not, it will be caught in validation
				if _, ok := txy.CompReqs[compReq]; ok {
					if l1Override.CompReqs == nil {
						l1Override.CompReqs = make(map[string]CompReq)
					}
					l1Override.CompReqs[compReq] = txy.CompReqs[compReq]
				}
			}
			secDomain.L1Overrides[l1ID] = l1Override
		}
	}
}

// ValidateSharedServices validates the shared-services environment.
// Hard coding some specific checks for shared-services environment as it is an exception to the normal security domains.
// This secenv need to meet the strictest requirements.
func (txy *Taxonomy) ValidateSharedServices() (bool, int) {
	valid := true
	envName := "shared-service"
	failures := 0
	if _, ok := txy.SegL1s[envName]; !ok {
		util.Log.Printf("%s environment not found", envName)
		return false, 1
	}
	if txy.SegL1s[envName].Sensitivity != SenseOrder[0] ||
		txy.SegL1s[envName].Criticality != CritOrder[0] {
		util.Log.Printf("%s environment does not have the highest sensitivity or criticality", envName)
		failures++
		valid = false
	}

	if len(txy.SegL1s[envName].ComplianceReqs) != len(txy.CompReqs) {
		util.Log.Printf("%s environment does not have all compliance requirements", envName)
		failures++
		valid = false
	}
	return valid, failures
}

func (txy *Taxonomy) CompleteAndValidateTaxonomy() bool {
	// Loop through SegL1s and validate max risk level
	valid := false
	valid = txy.validateEnv()
	if !valid {
		return valid
	}
	valid, _ = txy.ValidateSharedServices()
	if !valid {
		return valid
	}
	// Apply inheritance rules
	txy.ApplyInheritance()

	// Validate the taxonomy
	valid, _ = txy.validateSecurityDomains()
	return valid
}

// LoadTaxonomy loads the taxonomy by loading the different files and combining them into one struct.
// Validates the loaded data is valid and meets requirements.
// fills in missing data based on inheritance rules
// cfg parameter provides terminology configuration for directory resolution
func LoadTaxonomy(taxDir string, cfg Config) (Taxonomy, error) {
	txy := Taxonomy{
		ApiVersion: ApiVersion,
		Config:     cfg,
	}
	var err error

	if !strings.HasSuffix(taxDir, "/") {
		taxDir = taxDir + "/"
	}

	// Load L1 segments using configured directory name
	l1Dir := taxDir + cfg.Terminology.L1.DirName()
	txy.SegL1s, err = LoadSegL1Files(l1Dir)
	if err != nil {
		util.Log.Printf("Error loading L1 files from %s, exiting\n", l1Dir)
		return Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Load risk levels
	txy.SensitivityLevels = SenseOrder
	txy.CriticalityLevels = CritOrder

	// Load L2 segments using configured directory name
	l2Dir := taxDir + cfg.Terminology.L2.DirName()
	txy.SegL2s, err = LoadSegL2Files(l2Dir)
	if err != nil {
		util.Log.Printf("Error loading L2 files from %s: %v\n", l2Dir, err)
		return Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Define compliance scopes
	txy.CompReqs, err = LoadCompScope(taxDir + "compliance_requirements.yaml")
	if err != nil {
		util.Log.Println("Error loading compliance scope files:", err)
		return Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Validate the taxonomy
	valid := txy.CompleteAndValidateTaxonomy()
	// TODO validate against compliance scopes acceptable risk levels
	if !valid {
		util.Log.Println("Taxonomy is invalid")
		return Taxonomy{}, errors.New("invalid Taxonomy")
	} else {
		// Return the taxonomy
		return txy, nil
	}
}
