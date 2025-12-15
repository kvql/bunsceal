// Package domain provides configuration models for the taxonomy system.
package domain

import (
	"strings"

	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
	"github.com/kvql/bunsceal/pkg/taxonomy/infrastructure"
)

// DefaultSchemaPath is the default location for JSON schemas
const DefaultSchemaPath = "./pkg/domain/schemas"

// Config represents the taxonomy configuration.
type Config struct {
	Terminology  TermConfig                         `yaml:"terminology"`
	SchemaPath   string                             `yaml:"schema_path,omitempty"`
	Visuals      VisualsDef                         `yaml:"visuals,omitempty"`
	Rules        LogicRulesConfig                   `yaml:"rules,omitempty"`
	FsRepository infrastructure.ConfigFsReposistory `yaml:"fs_repository,omitempty"`
	Plugins      plugins.ConfigPlugins              `yaml:"plugins"`
}

// TermConfig holds terminology configuration for L1 and L2 segments.
type TermConfig struct {
	L1 TermDef `yaml:"l1,omitempty"`
	L2 TermDef `yaml:"l2,omitempty"`
}

// TermDef defines singular and plural forms for a segment level.
type TermDef struct {
	Singular string `yaml:"singular"`
	Plural   string `yaml:"plural"`
}

// VisualsDef Config for how taxonomy is visualised
type VisualsDef struct {
	L1Layout map[string][]string `yaml:"l1_layout,omitempty"`
}

// Merge merges this TermDef with defaults, using defaults for any blank fields.
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

// Merge merges this TermConfig with defaults, using defaults for any blank fields.
func (tc TermConfig) Merge(defaults TermConfig) TermConfig {
	return TermConfig{
		L1: tc.L1.Merge(defaults.L1),
		L2: tc.L2.Merge(defaults.L2),
	}
}

// Merge merges this Config with defaults, using defaults for any blank fields.
func (c Config) Merge() Config {
	defaults := DefaultConfig()
	result := Config{
		Terminology:  c.Terminology.Merge(defaults.Terminology),
		SchemaPath:   defaults.SchemaPath,
		FsRepository: defaults.FsRepository,
		Visuals:      c.Visuals,
	}
	if c.SchemaPath != "" {
		result.SchemaPath = c.SchemaPath
	}

	if c.FsRepository.L1Dir != "" {
		result.FsRepository.L1Dir = c.FsRepository.L1Dir
	}
	if c.FsRepository.L2Dir != "" {
		result.FsRepository.L2Dir = c.FsRepository.L2Dir
	}
	if c.FsRepository.TaxonomyDir != "" {
		result.FsRepository.TaxonomyDir = c.FsRepository.TaxonomyDir
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

// LogicRulesConfig holds configuration for business logic validation rules.
type LogicRulesConfig struct {
	SharedService GeneralBooleanConfig `yaml:"shared_service,omitempty"`
	Uniqueness    UniquenessConfig     `yaml:"uniqueness,omitempty"`
}

// GeneralBooleanConfig provides a simple enabled/disabled configuration for rules.
type GeneralBooleanConfig struct {
	Enabled bool `yaml:"enabled"`
}

// UniquenessConfig holds configuration for uniqueness validation rules.
type UniquenessConfig struct {
	Enabled   bool     `yaml:"enabled"`
	CheckKeys []string `yaml:"check_keys,omitempty"`
}

// DefaultConfig returns the default configuration.
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
		SchemaPath: DefaultSchemaPath,
		Rules: LogicRulesConfig{
			SharedService: GeneralBooleanConfig{Enabled: true},
			Uniqueness: UniquenessConfig{
				Enabled:   true,
				CheckKeys: []string{"name"},
			},
		},
		FsRepository: infrastructure.ConfigFsReposistory{
			TaxonomyDir: "taxonomy",
			L1Dir:       "environments",
			L2Dir:       "segments",
		},
	}
}
