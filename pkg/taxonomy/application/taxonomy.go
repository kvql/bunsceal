package application

import (
	"errors"
	"fmt"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
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

	schemaValidator, err := schemaValidation.NewSchemaValidator(cfg.SchemaPath, schemaValidation.SchemaBaseURL)
	if err != nil {
		o11y.Log.Printf("Error initialising schema validator: %v\n", err)
		return domain.Taxonomy{}, errors.New("failed to initialise schema validator")
	}

	// Load L2 segments using configured directory name
	FsRepository := infrastructure.NewFileSegRepository(schemaValidator, cfg.FsRepository)
	FsService := NewSegService(FsRepository)

	txy.SegL1s, err = FsService.LoadLevel("1")
	if err != nil {
		o11y.Log.Printf("Error loading L1 files. %s", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	txy.SegsL2s, err = FsService.LoadLevel("2")
	if err != nil {
		o11y.Log.Printf("Error loading L2 files. %s", err)
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Load plugins from config
	pluginsList := make(plugins.Plugins)
	if cfg.Plugins.Classifications != nil || cfg.Plugins.Compliance != nil {
		err = pluginsList.LoadPlugins(cfg.Plugins)
		if err != nil {
			o11y.Log.Printf("error loading plugins: %s", err)
			return domain.Taxonomy{}, errors.New("failed to load plugins")
		}
	}

	// Validate plugin labels BEFORE inheritance (L1 completeness, pairing for all)
	if err = ValidatePluginLabels(&txy, pluginsList); err != nil {
		o11y.Log.Println("Taxonomy is invalid: plugin label validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Apply inheritance (includes plugin label inheritance and order validation)
	if err = ApplyInheritance(&txy, pluginsList); err != nil {
		o11y.Log.Println("Taxonomy is invalid: inheritance failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Validate L2 definitions after inheritance
	var valid bool
	valid, _ = validation.ValidateL2Definition(&txy, pluginsList)
	if !valid {
		o11y.Log.Println("Taxonomy is invalid: Validate L2 definitions after inheritance")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Validate business logic rules
	valid = ValidateCoreLogic(&txy, cfg, pluginsList)
	if !valid {
		o11y.Log.Println("Taxonomy is invalid: business logic validation failed")
		return domain.Taxonomy{}, errors.New("invalid Taxonomy")
	}

	// Return the taxonomy
	return txy, nil
}

// ValidateCoreLogic validates business logic rules based on configuration.
// Returns true if all enabled rules pass, false if any rule fails.
func ValidateCoreLogic(txy *domain.Taxonomy, cfg configdomain.Config, pluginMap plugins.Plugins) bool {
	ruleSet := validation.NewLogicRuleSet(cfg, pluginMap)
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

// ValidatePluginLabels validates all segment labels against loaded plugins.
// Must be called BEFORE ApplyInheritance to catch malformed labels (missing rationale pairs).
// L1 segments must have all classification definitions. L2/overrides only need valid pairs.
func ValidatePluginLabels(txy *domain.Taxonomy, pluginsList plugins.Plugins) error {
	if pluginsList == nil {
		return nil
	}

	errs := pluginsList.ValidateAllSegments(txy.SegL1s, txy.SegsL2s)
	if len(errs) > 0 {
		for _, err := range errs {
			o11y.Log.Println(err)
		}
		return fmt.Errorf("plugin label validation failed with %d error(s)", len(errs))
	}

	return nil
}
