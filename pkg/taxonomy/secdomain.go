package taxonomy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type EnvDetails struct {
	Sensitivity          string   `yaml:"sensitivity"`
	SensitivityRationale string   `yaml:"sensitivity_rationale"`
	Criticality          string   `yaml:"criticality"`
	CriticalityRationale string   `yaml:"criticality_rationale"`
	ComplianceReqs       []string `yaml:"compliance_reqs"`
	CompReqs             map[string]CompReq
}
type version struct {
	Version string `yaml:"version"`
}

type SegL2 struct {
	Name        string                `yaml:"name"`
	ID          string                `yaml:"id"`
	Description string                `yaml:"description"`
	L1Overrides map[string]EnvDetails `yaml:"l1_overrides"`
}

// LoadSegL2Files loads all security domain files from the given directory
func LoadSegL2Files(segL2Dir string) (map[string]SegL2, error) {
	// Initialize schema validator
	// TODO: Refactor to accept schema path as parameter for testability
	schemaValidator, err := NewSchemaValidator("./schema")
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return nil, errors.New("failed to initialize schema validator")
	}

	segL2s := make(map[string]SegL2)
	err = filepath.WalkDir(segL2Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Load the file and parse it into a SegL2 struct
			segL2, err := parseSDFile(path, schemaValidator)
			if err != nil {
				return err
			}
			segL2s[segL2.ID] = segL2
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// validate id is unique across all domains. print error with non unique domains and exit if not
	validations := make([]string, 0)
	idMap := make(map[string]bool) //used to validate that the security domain id is unique
	outcome := true
	for _, segL2 := range segL2s {
		if _, ok := idMap[segL2.ID]; ok {
			validations = append(validations, "ID for "+segL2.Name+"is not unique: "+segL2.ID)
			outcome = false
		} else {
			idMap[segL2.ID] = true
		}
	}
	if !outcome {
		for _, result := range validations {
			util.Log.Println(result)
		}
		return nil, errors.New("Security domain validation failed, directory: " + segL2Dir)
	}

	return segL2s, nil
}

func parseSDFile(filePath string, schemaValidator *SchemaValidator) (SegL2, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	var fileVersion version
	err = yaml.Unmarshal(data, &fileVersion)
	if err != nil {
		return SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	switch fileVersion.Version {
	case "1.0":
		// Validate against JSON schema first
		if err := schemaValidator.ValidateData(data, "seg-level2.json"); err != nil {
			util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
			return SegL2{}, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
		}

		// Unmarshal the YAML data into a SegL2 struct
		var segL2 SegL2
		err = yaml.Unmarshal(data, &segL2)
		if err != nil {
			return SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
		}
		return segL2, nil
	default:
		return SegL2{}, errors.New("Unsupported security domain file version: " + filePath + fileVersion.Version)
	}
}
