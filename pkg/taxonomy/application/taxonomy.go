package application

import (
	"errors"

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

	// Load L1 segments using configured directory name

	schemaValidator, err := schemaValidation.NewSchemaValidator(cfg.SchemaPath)
	if err != nil {
		o11y.Log.Printf("Error initialising schema validator: %v\n", err)
		return domain.Taxonomy{}, errors.New("failed to initialise schema validator")
	}

	// Load risk levels
	txy.SensitivityLevels = domain.SenseOrder
	txy.CriticalityLevels = domain.CritOrder

	// Load L2 segments using configured directory name
	FsRepository := infrastructure.NewFileSegRepository(schemaValidator, cfg.FsRepository)
	FsService := NewSegService(FsRepository)

	txy.Segs, err = FsService.LoadLevel("1")
	if err != nil {
		o11y.Log.Printf("Error loading L1 files. %s", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	txy.Segs, err = FsService.LoadLevel("2")
	if err != nil {
		o11y.Log.Printf("Error loading L2 files. %s", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Define compliance scopes
	txy.CompReqs, err = infrastructure.LoadCompScope(cfg.FsRepository.TaxonomyDir+"compliance_requirements.yaml", schemaValidator)
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
