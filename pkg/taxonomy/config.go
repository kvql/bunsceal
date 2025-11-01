package taxonomy

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kvql/bunsceal/pkg/util"
	"errors"
	"gopkg.in/yaml.v3"
)

// Config represents the taxonomy configuration
type Config struct {
	Terminology TermConfig `yaml:"terminology"`
}

// TermConfig holds terminology configuration for L1 and L2 segments
type TermConfig struct {
	L1 TermDef `yaml:"l1"`
	L2 TermDef `yaml:"l2"`
}

// TermDef defines singular and plural forms for a segment level
type TermDef struct {
	Singular string `yaml:"singular"`
	Plural   string `yaml:"plural"`
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

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	// Validate completeness - both singular and plural required for each level
	// If either field is missing, use the entire default for that level
	if cfg.Terminology.L1.Singular != "" || cfg.Terminology.L1.Plural != "" {
		if cfg.Terminology.L1.Singular == "" || cfg.Terminology.L1.Plural == "" {
			return Config{}, errors.New("incomplete L1 Terminology definition, singular and plural forms required if either is set")
		}
	}else {
		cfg.Terminology.L1 = defaults.Terminology.L1
	}
	if cfg.Terminology.L2.Singular != "" || cfg.Terminology.L2.Plural != "" {
		if cfg.Terminology.L2.Singular == "" || cfg.Terminology.L2.Plural == "" {
			return Config{}, errors.New("incomplete L2 Terminology definition, singular and plural forms required if either is set")
		}
	}else {
		cfg.Terminology.L2 = defaults.Terminology.L2
	}

	return cfg, nil
}

// DirName converts the plural form to a kebab-case directory name
// Examples:
//   "Environments" → "environments"
//   "Security Environments" → "security-environments"
//   "My Custom Zones" → "my-custom-zones"
func (td TermDef) DirName() string {
	lower := strings.ToLower(td.Plural)
	kebab := strings.ReplaceAll(lower, " ", "-")
	return kebab
}
