package taxonomy

import (
	"os"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

// LoadCompScope loads compliance requirements from a YAML file and validates against schema
// schemaPath specifies the directory containing JSON schema files for validation
func LoadCompScope(filePath string, schemaPath string) (map[string]domain.CompReq, error) {
	// Initialize schema validator with provided path
	schemaValidator, err := validation.NewSchemaValidator(schemaPath)
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return nil, err
	}

	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		util.Log.Println("Error reading file:", err)
		return nil, err
	}

	// Validate against JSON schema first
	if err := schemaValidator.ValidateData(data, "compliance-reqs.json"); err != nil {
		util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
		return nil, err
	}

	// Parse the file into a allCompScopes struct
	var compReqs map[string]domain.CompReq
	err = yaml.Unmarshal(data, &compReqs)
	if err != nil {
		util.Log.Println("Error parsing file:", err)
		return nil, err
	}

	return compReqs, nil
}
