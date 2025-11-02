package taxonomy

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/taxonomy/domain"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"github.com/kvql/bunsceal/pkg/util"
	"gopkg.in/yaml.v3"
)

type version struct {
	Version string `yaml:"version"`
}


// LoadSegL2Files loads all security domain files from the given directory
func LoadSegL2Files(segL2Dir string) (map[string]domain.SegL2, error) {
	// Initialize schema validator
	// TODO: Refactor to accept schema path as parameter for testability
	schemaValidator, err := validation.NewSchemaValidator("./schema")
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return nil, errors.New("failed to initialize schema validator")
	}

	var segList []domain.SegL2
	parseErrors := false

	// Parse all files into a list
	err = filepath.WalkDir(segL2Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			// Load the file and parse it into a SegL2 struct
			segL2, err := parseSDFile(path, schemaValidator)
			if err != nil {
				util.Log.Printf("Error parsing file: %s\n", path)
				parseErrors = true
				return nil // Continue walking despite parse error
			}
			segList = append(segList, segL2)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// Validate uniqueness using extracted validator
	validations := validation.UniquenessValidator(segList)

	if parseErrors || len(validations) > 0 {
		for _, result := range validations {
			util.Log.Println(result)
		}
		return nil, fmt.Errorf("security domain validation failed, directory: %s", segL2Dir)
	}

	// Build map from validated list
	segL2Map := make(map[string]domain.SegL2)
	for _, seg := range segList {
		segL2Map[seg.ID] = seg
	}

	return segL2Map, nil
}

func parseSDFile(filePath string, schemaValidator *validation.SchemaValidator) (domain.SegL2, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	var fileVersion version
	err = yaml.Unmarshal(data, &fileVersion)
	if err != nil {
		return domain.SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
	}
	switch fileVersion.Version {
	case "1.0":
		// Validate against JSON schema first
		if err := schemaValidator.ValidateData(data, "seg-level2.json"); err != nil {
			util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
			return domain.SegL2{}, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
		}

		// Unmarshal the YAML data into a SegL2 struct
		var segL2 domain.SegL2
		err = yaml.Unmarshal(data, &segL2)
		if err != nil {
			return domain.SegL2{}, errors.New("Failed to parse file" + filePath + err.Error())
		}
		return segL2, nil
	default:
		return domain.SegL2{}, errors.New("Unsupported security domain file version: " + filePath + fileVersion.Version)
	}
}
