package taxonomy

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from the specified path.
// If configPath is empty, loads from taxDir/config.yaml.
// If the config file doesn't exist or has missing fields, uses defaults.
func LoadConfig(configPath, configSchemaPath string) (domain.Config, error) {
	defaults := domain.DefaultConfig()

	schemaValidator, err := validation.NewSchemaValidator(configSchemaPath)
	if err != nil {
		util.Log.Printf("Error initialising schema validator: %v\n", err)
		return domain.Config{}, errors.New("failed to initialise schema validator")
	}
	// Determine config file location
	if configPath == "" {
		configPath = filepath.Join("config.yaml")
	}

	// #nosec G304 -- configPath comes from CLI flag or default config location, not user input
	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file missing - use full defaults
		if os.IsNotExist(err) {
			util.Log.Printf("Config file not found at %s, using defaults\n", configPath)
			return defaults, nil
		}
		return domain.Config{}, err
	}
	if err := schemaValidator.ValidateData(data, "config.json"); err != nil {
		util.Log.Printf("Schema validation failed for Config")
		return domain.Config{}, errors.New("schema validation failed for Config")
	}

	var loadedConfig domain.Config
	if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
		return domain.Config{}, err
	}

	// Merge loaded config with defaults
	merged := loadedConfig.Merge()
	configDir := filepath.Dir(configPath)

	if merged.TaxonomyPath != "" && (!strings.HasPrefix(merged.TaxonomyPath, "/") || !strings.HasPrefix(merged.TaxonomyPath, "\\")) {
		merged.TaxonomyPath = filepath.Join(configDir, merged.TaxonomyPath)
	}
	return merged, nil
}
