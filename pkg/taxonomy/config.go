package taxonomy

import (
	"os"
	"path/filepath"
	"strings"

	"errors"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

// Config represents the taxonomy configuration
type Config struct {
	Terminology TermConfig `yaml:"terminology"`
}

// TermConfig holds terminology configuration for L1 and L2 segments
type TermConfig struct {
	L1 TermDef `yaml:"l1,omitempty"`
	L2 TermDef `yaml:"l2,omitempty"`
}

// TermDef defines singular and plural forms for a segment level
type TermDef struct {
	Singular string `yaml:"singular"`
	Plural   string `yaml:"plural"`
}

// Merge merges this TermDef with defaults, using defaults for any blank fields
func (td TermDef) Merge(defaults TermDef) TermDef {
	result := defaults
	if td.Singular != "" {
		result.Singular = td.Singular
	}
	if td.Plural != "" {
		result.Plural = td.Plural
	}
	return result
}

// Merge merges this TermConfig with defaults, using defaults for any blank fields
func (tc TermConfig) Merge(defaults TermConfig) TermConfig {
	return TermConfig{
		L1: tc.L1.Merge(defaults.L1),
		L2: tc.L2.Merge(defaults.L2),
	}
}

// DefaultConfig returns the default configuration
func DefaultConfig() Config {
	return Config{
		Terminology: TermConfig{
			L1: TermDef{
				Singular: "Environment",
				Plural:   "Environments",
			},
			L2: TermDef{
				Singular: "Segment",
				Plural:   "Segments",
			},
		},
	}
}

// LoadConfig loads configuration from the specified path
// If configPath is empty, loads from taxDir/config.yaml
// If the config file doesn't exist or has missing fields, uses defaults
func LoadConfig(configPath string, taxDir string) (Config, error) {
	defaults := DefaultConfig()

	schemaValidator, err := NewSchemaValidator("../../schema")
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return Config{}, errors.New("failed to initialize schema validator")
	}
	// Determine config file location
	if configPath == "" {
		configPath = filepath.Join(taxDir, "config.yaml")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		// Config file missing - use full defaults
		util.Log.Printf("Config file not found at %s, using defaults\n", configPath)
		return defaults, nil
	}
	if err := schemaValidator.ValidateData(data, "config.json"); err != nil {
		util.Log.Printf("Schema validation failed for Config")
		return Config{}, errors.New("schema validation failed for Config")
	}

	var loadedConfig Config
	if err := yaml.Unmarshal(data, &loadedConfig); err != nil {
		return Config{}, err
	}

	// Merge loaded config with defaults
	merged := Config{
		Terminology: loadedConfig.Terminology.Merge(defaults.Terminology),
	}

	return merged, nil
}

// DirName converts the plural form to a kebab-case directory name
// Examples:
//
//	"Environments" → "environments"
//	"Security Environments" → "security-environments"
//	"My Custom Zones" → "my-custom-zones"
func (td TermDef) DirName() string {
	lower := strings.ToLower(td.Plural)
	kebab := strings.ReplaceAll(lower, " ", "-")
	return kebab
}
