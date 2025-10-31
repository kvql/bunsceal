package taxonomy

import (
	"os"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type CompReq struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	ReqsLink    string `yaml:"requirements_link"`
}

func LoadCompScope(filePath string) (map[string]CompReq, error) {
	// Initialize schema validator
	schemaValidator, err := NewSchemaValidator("./schema")
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
	if err := schemaValidator.ValidateYAML(data, "compliance-requirements.json"); err != nil {
		util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
		return nil, err
	}

	// Parse the file into a allCompScopes struct
	var compReqs map[string]CompReq
	err = yaml.Unmarshal(data, &compReqs)
	if err != nil {
		util.Log.Println("Error parsing file:", err)
		return nil, err
	}

	return compReqs, nil
}
