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

// LoadSegL1Files Parse all security environment files from the provided directory,
// validate and return a map of SegL1structs
func LoadSegL1Files(segL1Dir string) (map[string]domain.SegL1, error) {
	// Initialize schema validator
	// TODO: Refactor to accept schema path as parameter for testability
	schemaValidator, err := validation.NewSchemaValidator("./schema")
	if err != nil {
		util.Log.Printf("Error initializing schema validator: %v\n", err)
		return nil, errors.New("failed to initialize schema validator")
	}

	files, err := os.ReadDir(segL1Dir)
	if err != nil {
		return nil, err
	}

	var segList []domain.SegL1
	parseErrors := false

	// Parse all files into a list
	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(segL1Dir, file.Name())
			// Load the file and parse it into a SegL1struct
			segL1, err := parseSegL1File(filePath, schemaValidator)
			if err != nil {
				util.Log.Printf("Error parsing file: %s\n", filePath)
				parseErrors = true
				continue
			}
			segList = append(segList, segL1)
		}
	}

	// Validate uniqueness using extracted validator
	validations := validation.UniquenessValidator(segList)

	if parseErrors || len(validations) > 0 {
		for _, result := range validations {
			util.Log.Println(result)
		}
		return nil, fmt.Errorf("security environment validation failed, directory: %s", segL1Dir)
	}

	// Build map from validated list
	segL1Map := make(map[string]domain.SegL1)
	for _, seg := range segList {
		segL1Map[seg.ID] = seg
	}

	return segL1Map, nil
}

func parseSegL1File(filePath string, schemaValidator *validation.SchemaValidator) (domain.SegL1, error) {
	// Read the file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.SegL1{}, err
	}

	// Validate against JSON schema first
	if err := schemaValidator.ValidateData(data, "seg-level1.json"); err != nil {
		util.Log.Printf("Schema validation failed for %s: %v\n", filePath, err)
		return domain.SegL1{}, fmt.Errorf("schema validation failed for %s: %w", filePath, err)
	}

	// Unmarshal the YAML data into a SegL1struct
	var segL1 domain.SegL1
	err = yaml.Unmarshal(data, &segL1)
	if err != nil {
		return domain.SegL1{}, err
	}

	return segL1, nil
}
