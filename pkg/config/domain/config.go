// Package domain provides configuration models for the taxonomy system.
package domain

import (
	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/application/plugins"
	"github.com/kvql/bunsceal/pkg/taxonomy/infrastructure"
	"github.com/kvql/bunsceal/pkg/visualise"
)

// DefaultSchemaPath is the default location for JSON schemas
const DefaultSchemaPath = "./pkg/domain/schemas"

// Config represents the taxonomy configuration.
type Config struct {
	Terminology  domain.TermConfig                  `yaml:"terminology"`
	SchemaPath   string                             `yaml:"schema_path,omitempty"`
	Visuals      visualise.VisualsDef               `yaml:"visuals,omitempty"`
	Rules        LogicRulesConfig                   `yaml:"rules,omitempty"`
	FsRepository infrastructure.ConfigFsReposistory `yaml:"fs_repository,omitempty"`
	Plugins      plugins.ConfigPlugins              `yaml:"plugins"`
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
		Terminology: domain.TermConfig{
			L1: domain.TermDef{
				Singular: "Environment",
				Plural:   "Environments",
			},
			L2: domain.TermDef{
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
