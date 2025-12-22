package plugins

// ClassificationsConfigSchema defines the JSON schema for classifications plugin config
const ClassificationsConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/plugin-classifications.json",
	"title": "Classifications Plugin Configuration",
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"common_settings": {
			"type": "object",
			"additionalProperties": false,
			"properties": {
				"label_inheritance": { "type": "boolean" },
				"require_complete_l1": { "type": "boolean" }
			}
		},
		"rationale_length": { "type": "integer", "minimum": 0 },
		"definitions": {
			"type": "object",
			"additionalProperties": {
				"type": "object",
				"required": ["name", "values"],
				"additionalProperties": false,
				"properties": {
					"name": { "type": "string", "minLength": 1 },
					"description": { "type": "string" },
					"enforce_order": { "type": "boolean", "default": true },
					"values": {
						"type": "object",
						"additionalProperties": { "type": "string" }
					},
					"order": {
						"type": "array",
						"items": { "type": "string" }
					}
				}
			}
		}
	}
}`

// ComplianceConfigSchema defines the JSON schema for compliance plugin config
const ComplianceConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/plugin-compliance.json",
	"title": "Compliance Plugin Configuration",
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"common_settings": {
			"type": "object",
			"additionalProperties": false,
			"properties": {
				"label_inheritance": { "type": "boolean" },
				"require_complete_l1": { "type": "boolean" }
			}
		},
		"rationale_length": { "type": "integer", "minimum": 0 },
		"enforce_scope_hierarchy": { "type": "boolean" },
		"definitions": {
			"type": "object",
			"additionalProperties": {
				"type": "object",
				"required": ["name"],
				"additionalProperties": false,
				"properties": {
					"name": { "type": "string", "minLength": 1 },
					"description": { "type": "string" },
					"requirements_link": { "type": "string", "format": "uri" }
				}
			}
		}
	}
}`

// PluginsConfigSchema wraps all plugin schemas for the plugins section
const PluginsConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/plugins.json",
	"title": "Plugins Configuration",
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"classifications": { "$ref": "./plugin-classifications.json" },
		"compliance": { "$ref": "./plugin-compliance.json" }
	}
}`
