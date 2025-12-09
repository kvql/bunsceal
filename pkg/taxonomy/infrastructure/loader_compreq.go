package infrastructure

import (
	"os"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/schemaValidation"

	"gopkg.in/yaml.v3"
)

// LoadCompScope loads compliance requirements from a YAML file and validates against schema
// schemaPath specifies the directory containing JSON schema files for validation
func LoadCompScope(filePath string, schemaPath string) (map[string]domain.CompReq, error) {
	// Initialise schema validator with provided path
	schemaValidator, err := schemaValidation.NewSchemaValidator(schemaPath)
	if err != nil {
		o11y.Log.Printf("Error initialising schema validator: %v\n", err)
		return nil, err
	}

	// Read the file
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		o11y.Log.Println("Error reading file:", err)
		return nil, err
	}

	// Validate against JSON schema first
	if validationErr := schemaValidator.ValidateData(data, "compliance-reqs.json"); validationErr != nil {
		o11y.Log.Printf("Schema validation failed for %s: %v\n", filePath, validationErr)
		return nil, validationErr
	}

	// Parse the file into a allCompScopes struct
	var compReqs map[string]domain.CompReq
	err = yaml.Unmarshal(data, &compReqs)
	if err != nil {
		o11y.Log.Println("Error parsing file:", err)
		return nil, err
	}

	return compReqs, nil
}
