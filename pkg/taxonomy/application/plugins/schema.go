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
				"label_inheritance": { "type": "boolean" }
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

// PluginsConfigSchema wraps all plugin schemas for the plugins section
const PluginsConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/plugins.json",
	"title": "Plugins Configuration",
	"type": "object",
	"additionalProperties": false,
	"properties": {
		"classifications": { "$ref": "./plugin-classifications.json" }
	}
}`
