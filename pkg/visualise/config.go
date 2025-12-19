package visualise

// VisualsDef Config for how taxonomy is visualised
type VisualsDef struct {
	L1Layout map[string][]string `yaml:"l1_layout,omitempty"`
}

const VisualiseConfigSchema = `{
	"$schema": "https://json-schema.org/draft/2020-12/schema",
	"$id": "https://github.com/kvql/bunsceal/pkg/config/schemas/visualise.json",
	"title": "Visualise Configuration",
	"$defs": {
		"visuals": {
			"type": "object",
			"description": "Configuration options for visualisation functions",
			"additionalProperties": false,
			"properties": {
				"l1_layout": {
					"type": "object",
					"description": "Map of row numbers to ordered arrays of L1 identifiers for layout control",
					"additionalProperties": false,
					"patternProperties": {
						"^[0-9]+$": {
							"type": "array",
							"items": {
								"type": "string"
							}
						}
					}
				}
			}
		}
	}
}`