package domain

import (
	"strings"
)

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

const TermsConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/terms.json",
	"title": "Terminology Configuration",
	"$defs": {
		"termDef": {
			"type": "object",
			"description": "terminology Override to utilise terminology that is understandable to your audience",
			"properties": {
				"singular": {
					"type": "string",
					"minLength": 1,
					"description": "Singular form of the term"
				},
				"plural": {
					"type": "string",
					"minLength": 1,
					"description": "Plural form of the term"
				}
			},
			"required": [
				"singular",
				"plural"
			],
			"additionalProperties": false
		},
		"terms": {
			"type": "object",
			"description": "Terminology configuration for L1 and L2 segments",
			"properties": {
				"l1": {
				"$ref": "#/$defs/termDef" 
				},
				"l2": {
					"$ref": "#/$defs/termDef" 
				}
			},
			"additionalProperties": false
		}	
	}
}`
