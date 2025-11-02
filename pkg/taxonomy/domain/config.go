package domain

import "strings"

// Config represents the taxonomy configuration
type Config struct {
	Terminology TermConfig `yaml:"terminology"`
	SchemaPath  string     `yaml:"schema_path,omitempty"`
	TaxonomyPath  string     `yaml:"taxonomy_path,omitempty"`
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

// Merge merges this Config with defaults, using defaults for any blank fields
func (c Config) Merge() Config {
	defaults := DefaultConfig()
	result := Config{
		Terminology: c.Terminology.Merge(defaults.Terminology),
		SchemaPath:  defaults.SchemaPath,
		TaxonomyPath: defaults.TaxonomyPath,
	}
	if c.SchemaPath != "" {
		result.SchemaPath = c.SchemaPath
	}
	if c.TaxonomyPath != "" {
		result.TaxonomyPath = c.TaxonomyPath
	}
	return result
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
		SchemaPath: "./schema",
		TaxonomyPath: "taxonomy",
	}
}
