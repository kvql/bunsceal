package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	configDomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
	"gopkg.in/yaml.v3"
)

var defaultConfigSchemaPath = "pkg/config/schemas"

const configSchemaBaseURL = "https://github.com/kvql/bunsceal/pkg/config/schemas/"

// LoadConfig loads configuration from the specified path.
// If configPath is empty, loads from taxDir/config.yaml.
// If configSchemaPath is empty, uses domain.DefaultSchemaPath.
// If the config file doesn't exist or has missing fields, uses defaults.
func LoadConfig(configPath, configSchemaPath string) (configDomain.Config, error) {
	defaults := configDomain.DefaultConfig()

	// Use default schema path if not provided
	if configSchemaPath == "" {
		configSchemaPath = defaultConfigSchemaPath
	}

	// Plugin schemas must be registered before compilation (config.json refs plugins.json)
	pluginSchemas := []schemaValidation.ExternalSchema{
		{JSON: plugins.ClassificationsConfigSchema, ID: "https://github.com/kvql/bunsceal/pkg/config/schemas/plugin-classifications.json"},
		{JSON: plugins.PluginsConfigSchema, ID: "https://github.com/kvql/bunsceal/pkg/config/schemas/plugins.json"},
		{JSON: domain.TermsConfigSchema, ID: "https://github.com/kvql/bunsceal/pkg/config/schemas/terms.json"},
	}
	schemaValidator, err := schemaValidation.NewSchemaValidator(configSchemaPath, configSchemaBaseURL, pluginSchemas...)
	if err != nil {
		o11y.Log.Printf("Error initialising schema validator: %v\n", err)
		return configDomain.Config{}, errors.New("failed to initialise schema validator")
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
			o11y.Log.Printf("Config file not found at %s, using defaults\n", configPath)
			return defaults, nil
		}
		return configDomain.Config{}, err
	}
	if err := schemaValidator.ValidateData(data, "config.json"); err != nil {
		o11y.Log.Printf("Schema validation failed for Config %s", err)
		return configDomain.Config{}, fmt.Errorf("schema validation failed for Config. Error: %w", err)
	}

	var loadedConfig configDomain.Config
	if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
		return configDomain.Config{}, err
	}

	// Merge loaded config with defaults
	merged := loadedConfig.Merge()
	configDir := filepath.Dir(configPath)

	// Update Taxonomy path if relative
	if merged.FsRepository.TaxonomyDir != "" && (!strings.HasPrefix(merged.FsRepository.TaxonomyDir, "/") || !strings.HasPrefix(merged.FsRepository.TaxonomyDir, "\\")) {
		merged.FsRepository.TaxonomyDir = filepath.Join(configDir, merged.FsRepository.TaxonomyDir)
	}

	return merged, nil
}
