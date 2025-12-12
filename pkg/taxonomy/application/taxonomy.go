package application

import (
	"errors"
	"strings"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/validation"
	"github.com/kvql/bunsceal/pkg/taxonomy/infrastructure"
)

// LoadTaxonomy loads the taxonomy by loading the different files and combining them into one struct.
// Validates the loaded data is valid and meets requirements.
// Fills in missing data based on inheritance rules.
// cfg parameter provides terminology configuration for directory resolution.
func LoadTaxonomy(cfg configdomain.Config) (domain.Taxonomy, error) {
	txy := domain.Taxonomy{
		ApiVersion: domain.ApiVersion,
	}
	var err error
	taxDir := cfg.TaxonomyPath

	if !strings.HasSuffix(taxDir, "/") {
		taxDir = taxDir + "/"
	}

	// Load L1 segments using configured directory name
	l1Dir := taxDir + cfg.Terminology.L1.DirName()
	schemaValidator, err := schemaValidation.NewSchemaValidator(cfg.SchemaPath)
	if err != nil {
		o11y.Log.Printf("Error initialising schema validator: %v\n", err)
		return domain.Taxonomy{}, errors.New("failed to initialise schema validator")
	}
	l1Repository := infrastructure.NewFileSegL1Repository(schemaValidator)
	l1Service := NewSegL1Service(l1Repository)
	txy.SegL1s, err = l1Service.Load(l1Dir)
	if err != nil {
		o11y.Log.Printf("Error loading L1 files from %s, exiting\n", l1Dir)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Load risk levels
	txy.SensitivityLevels = domain.SenseOrder
	txy.CriticalityLevels = domain.CritOrder

	// Load L2 segments using configured directory name
	l2Dir := taxDir + cfg.Terminology.L2.DirName()
	l2Repository := infrastructure.NewFileSegRepository(schemaValidator)
	l2Service := NewSegService(l2Repository)
	txy.Segs, err = l2Service.Load(l2Dir)
	if err != nil {
		o11y.Log.Printf("Error loading L2 files from %s: %v\n", l2Dir, err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Define compliance scopes
	txy.CompReqs, err = infrastructure.LoadCompScope(taxDir+"compliance_requirements.yaml", schemaValidator)
	if err != nil {
		o11y.Log.Println("Error loading compliance scope files:", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Apply inheritance and validate cross-entity references
	valid := ApplyInheritance(&txy)
	if !valid {
		o11y.Log.Println("Taxonomy is invalid: cross-entity validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Validate business logic rules
	valid = ValidateCoreLogic(&txy, cfg)
	// TODO validate against compliance scopes acceptable risk levels
	if !valid {
		o11y.Log.Println("Taxonomy is invalid: business logic validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Return the taxonomy
	return txy, nil
}

// ValidateCoreLogic validates business logic rules based on configuration.
// Returns true if all enabled rules pass, false if any rule fails.
func ValidateCoreLogic(txy *domain.Taxonomy, cfg configdomain.Config) bool {
	ruleSet := validation.NewLogicRuleSet(cfg)
	results := ruleSet.ValidateAll(txy)

	if len(results) > 0 {
		o11y.Log.Printf("Business logic validation failed with %d rule(s) reporting errors:", len(results))
		for _, result := range results {
			o11y.Log.Printf("  Rule '%s' failed with %d error(s)", result.RuleName, len(result.Errors))
		}
		return false
	}

	return true
}
