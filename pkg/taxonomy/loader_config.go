package taxonomy

import (
	"errors"
	"os"
	"strings"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

// LoadConfig loads configuration from the specified path
// If configPath is empty, loads from taxDir/config.yaml
// If the config file doesn't exist or has missing fields, uses defaults
func LoadConfig(configPath string, configSchemaPath string) (domain.Config, error) {
	defaults := domain.DefaultConfig()

	schemaValidator, err := validation.NewSchemaValidator(configSchemaPath)
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return domain.Config{}, errors.New("failed to initialize schema validator")
	}
	// Determine config file location
	if configPath == "" {
		configPath = filepath.Join("config.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file missing - use full defaults
		util.Log.Printf("Config file not found at %s, using defaults\n", configPath)
		return defaults, nil
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
	
	if merged.TaxonomyPath !="" && ( !strings.HasPrefix(merged.TaxonomyPath, "/") || !strings.HasPrefix(merged.TaxonomyPath, "\\")) {
		merged.TaxonomyPath = filepath.Join(configDir, merged.TaxonomyPath)
	}
	return merged, nil
}
