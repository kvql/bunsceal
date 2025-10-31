package taxonomy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type SegL1 struct {
	Name                 string   `yaml:"name"`
	ID                   string   `yaml:"id"`
	Description          string   `yaml:"description"`
	Sensitivity          string   `yaml:"sensitivity"`
	SensitivityRationale string   `yaml:"sensitivity_rationale"`
	Criticality          string   `yaml:"criticality"`
	CriticalityRationale string   `yaml:"criticality_rationale"`
	ComplianceReqs       []string `yaml:"compliance_reqs"`
}

// LoadSegL1Files Parse all security environment files from the provided directory,
// validate and return a map of SegL1structs
func LoadSegL1Files(segL1Dir string) (map[string]SegL1, error) {
	// Initialize schema validator
	schemaValidator, err := NewSchemaValidator("./schema")
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return nil, errors.New("failed to initialize schema validator")
	}

	files, err := os.ReadDir(segL1Dir)
	if err != nil {
		return nil, err
	}
	valid := true
	segL1s := make(map[string]SegL1)
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(segL1Dir, file.Name())
			// Load the file and parse it into a SegL1struct
			segL1, err := parseSegL1File(filePath, schemaValidator)
			if err != nil {
				util.Log.Printf("Error parsing file: %s\n", filePath)
				valid = false
			}
			segL1s[segL1.ID] = segL1
		}
	}
	if valid {
		return segL1s, nil
	} else {
		return nil, errors.New("loading security environments failed")
	}
}

func parseSegL1File(filePath string, schemaValidator *SchemaValidator) (SegL1, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return SegL1{}, err
	}

	// Validate against JSON schema first
	if err := schemaValidator.ValidateYAML(data, "seg-level1.json"); err != nil {
		util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
		return SegL1{}, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
	}

	// Unmarshal the YAML data into a SegL1struct
	var segL1 SegL1
	err = yaml.Unmarshal(data, &segL1)
	if err != nil {
		return SegL1{}, err
	}

	return segL1, nil
}
