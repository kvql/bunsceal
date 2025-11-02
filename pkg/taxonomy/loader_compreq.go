package taxonomy

import (
	"os"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

func LoadCompScope(filePath string) (map[string]domain.CompReq, error) {
	// Initialize schema validator
	// TODO: Refactor to accept schema path as parameter for testability
	schemaValidator, err := validation.NewSchemaValidator("./schema")
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
