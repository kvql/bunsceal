package taxonomy

import (
	"errors"
	"strings"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
)

// ApplyInheritance applies inheritance rules for taxonomy segments and validates cross-entity references.
// Returns true if all validations pass, false otherwise.
func ApplyInheritance(txy *domain.Taxonomy) bool {
	// Validate L1 definitions reference valid compliance requirements
	valid := validation.ValidateL1Definitions(txy)
	if !valid {
		return valid
	}

	// Loop through env details for each security domain and update risk compliance if not set based on env default
	for _, segL2 := range txy.SegL2s {
		for l1ID, l1Override := range segL2.L1Overrides {
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
			// add compliance details to compReqs var for each compliance standard listed
			for _, compReq := range l1Override.ComplianceReqs {
				// only add details if listed standard is valid. If not, it will be caught in validation
				if _, ok := txy.CompReqs[compReq]; ok {
					if l1Override.CompReqs == nil {
						l1Override.CompReqs = make(map[string]domain.CompReq)
					}
					l1Override.CompReqs[compReq] = txy.CompReqs[compReq]
				}
			}
			segL2.L1Overrides[l1ID] = l1Override
		}
	}

	// Validate L2 definitions after inheritance
	valid, _ = validation.ValidateL2Definition(txy)
	return valid
}

// CompleteAndValidateTaxonomy completes and validates the taxonomy.
// Deprecated: Use ApplyInheritance and ValidateBusinessLogic separately for better control.
func CompleteAndValidateTaxonomy(txy *domain.Taxonomy) bool {
	return ApplyInheritance(txy)
}

// ValidateBusinessLogic validates business logic rules based on configuration.
// Returns true if all enabled rules pass, false if any rule fails.
func ValidateBusinessLogic(txy *domain.Taxonomy) bool {
	ruleSet := NewLogicRuleSet(txy.Config)
	results := ruleSet.ValidateAll(txy)

	if len(results) > 0 {
		util.Log.Printf("Business logic validation failed with %d rule(s) reporting errors:", len(results))
		for _, result := range results {
			util.Log.Printf("  Rule '%s' failed with %d error(s)", result.RuleName, len(result.Errors))
		}
		return false
	}

	return true
}

// InitTaxonomy defines the interface for taxonomy initialization.
var InitTaxonomy interface {
	Load()
}

// LoadTaxonomy loads the taxonomy by loading the different files and combining them into one struct.
// Validates the loaded data is valid and meets requirements.
// Fills in missing data based on inheritance rules.
// cfg parameter provides terminology configuration for directory resolution.
func LoadTaxonomy(cfg domain.Config) (domain.Taxonomy, error) {
	txy := domain.Taxonomy{
		ApiVersion: domain.ApiVersion,
		Config:     cfg,
	}
	var err error
	taxDir := cfg.TaxonomyPath

	if !strings.HasSuffix(taxDir, "/") {
		taxDir = taxDir + "/"
	}

	// Load L1 segments using configured directory name
	l1Dir := taxDir + cfg.Terminology.L1.DirName()
	schemaValidator, err := validation.NewSchemaValidator(cfg.SchemaPath)
	if err != nil {
		util.Log.Printf("Error initialising schema validator: %v\n", err)
		return domain.Taxonomy{}, errors.New("failed to initialise schema validator")
	}
	l1Repository := NewFileSegL1Repository(schemaValidator)
	l1Service := NewSegL1Service(l1Repository)
	txy.SegL1s, err = l1Service.LoadAndValidate(l1Dir)
	if err != nil {
		util.Log.Printf("Error loading L1 files from %s, exiting\n", l1Dir)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Load risk levels
	txy.SensitivityLevels = domain.SenseOrder
	txy.CriticalityLevels = domain.CritOrder

	// Load L2 segments using configured directory name
	l2Dir := taxDir + cfg.Terminology.L2.DirName()
	l2Repository := NewFileSegL2Repository(schemaValidator)
	l2Service := NewSegL2Service(l2Repository)
	txy.SegL2s, err = l2Service.LoadAndValidate(l2Dir)
	if err != nil {
		util.Log.Printf("Error loading L2 files from %s: %v\n", l2Dir, err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Define compliance scopes
	txy.CompReqs, err = LoadCompScope(taxDir+"compliance_requirements.yaml", cfg.SchemaPath)
	if err != nil {
		util.Log.Println("Error loading compliance scope files:", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Apply inheritance and validate cross-entity references
	valid := ApplyInheritance(&txy)
	if !valid {
		util.Log.Println("Taxonomy is invalid: cross-entity validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Validate business logic rules
	valid = ValidateBusinessLogic(&txy)
	// TODO validate against compliance scopes acceptable risk levels
	if !valid {
		util.Log.Println("Taxonomy is invalid: business logic validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Return the taxonomy
	return txy, nil
}
