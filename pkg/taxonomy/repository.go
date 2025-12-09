package taxonomy

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/o11y"
	"github.com/kvql/bunsceal/pkg/taxonomy/validation"
	"gopkg.in/yaml.v3"
)

// SegL1Repository defines the contract for loading SegL1 data from any source
type SegL1Repository interface {
	// LoadAll loads all SegL1 entities from the specified source
	// Returns a slice of SegL1 entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.SegL1, error)
}

// SegL2Repository defines the contract for loading SegL2 data from any source
type SegL2Repository interface {
	// LoadAll loads all SegL2 entities from the specified source
	// Returns a slice of SegL2 entities or an error if loading fails
	// Does NOT perform business rule validation (uniqueness, etc)
	LoadAll(source string) ([]domain.SegL2, error)
}

// FileSegL1Repository implements SegL1Repository for file-based data sources
type FileSegL1Repository struct {
	schemaValidator *validation.SchemaValidator
}

// NewFileSegL1Repository creates a new file-based SegL1 repository
// schemaValidator is used to validate each file against the JSON schema
func NewFileSegL1Repository(schemaValidator *validation.SchemaValidator) *FileSegL1Repository {
	return &FileSegL1Repository{
		schemaValidator: schemaValidator,
	}
}

// LoadAll loads all SegL1 files from the specified directory
// Returns a slice of SegL1 entities or an error if loading fails
// Performs schema validation but NOT business rule validation
func (r *FileSegL1Repository) LoadAll(segL1Dir string) ([]domain.SegL1, error) {
	files, err := os.ReadDir(segL1Dir)
	if err != nil {
		return nil, err
	}

	var segList []domain.SegL1
	var parseErrors []error

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(segL1Dir, file.Name())
			segL1, err := r.parseSegL1File(filePath)
			if err != nil {
				o11y.Log.Printf("Error parsing file %s: %v\n", filePath, err)
				parseErrors = append(parseErrors, err)
				continue
			}
			segList = append(segList, segL1)
		}
	}

	if len(parseErrors) > 0 {
		return nil, fmt.Errorf("failed to parse %d file(s) in directory: %s", len(parseErrors), segL1Dir)
	}

	return segList, nil
}

func (r *FileSegL1Repository) parseSegL1File(filePath string) (domain.SegL1, error) {
	return parseSegL1File(filePath, r.schemaValidator)
}

type version struct {
	Version string `yaml:"version"`
}

// parseSegL1File parses a single SegL1 file with schema validation
func parseSegL1File(filePath string, schemaValidator *validation.SchemaValidator) (domain.SegL1, error) {
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.SegL1{}, err
	}

	if validationErr := schemaValidator.ValidateData(data, "seg-level1.json"); validationErr != nil {
		return domain.SegL1{}, fmt.Errorf("schema validation failed for %s: %w", filePath, validationErr)
	}

	var segL1 domain.SegL1
	if err = yaml.Unmarshal(data, &segL1); err != nil {
		return domain.SegL1{}, err
	}

	return segL1, nil
}

// FileSegL2Repository implements SegL2Repository for file-based data sources
type FileSegL2Repository struct {
	schemaValidator *validation.SchemaValidator
}

// NewFileSegL2Repository creates a new file-based SegL2 repository
// schemaValidator is used to validate each file against the JSON schema
func NewFileSegL2Repository(schemaValidator *validation.SchemaValidator) *FileSegL2Repository {
	return &FileSegL2Repository{
		schemaValidator: schemaValidator,
	}
}

// LoadAll loads all SegL2 files from the specified directory
// Returns a slice of SegL2 entities or an error if loading fails
// Performs schema validation but NOT business rule validation
func (r *FileSegL2Repository) LoadAll(segL2Dir string) ([]domain.SegL2, error) {
	var segList []domain.SegL2
	var parseErrors []error

	err := filepath.WalkDir(segL2Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			segL2, err := r.parseSegL2File(path)
			if err != nil {
				o11y.Log.Printf("Error parsing file %s: %v\n", path, err)
				parseErrors = append(parseErrors, err)
				return nil // Continue walking despite parse error
			}
			segList = append(segList, segL2)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(parseErrors) > 0 {
		return nil, fmt.Errorf("failed to parse %d file(s) in directory: %s", len(parseErrors), segL2Dir)
	}

	return segList, nil
}

func (r *FileSegL2Repository) parseSegL2File(filePath string) (domain.SegL2, error) {
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.SegL2{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var fileVersion version
	if err = yaml.Unmarshal(data, &fileVersion); err != nil {
		return domain.SegL2{}, fmt.Errorf("failed to parse version from file %s: %w", filePath, err)
	}

	switch fileVersion.Version {
	case "1.0":
		if validationErr := r.schemaValidator.ValidateData(data, "seg-level2.json"); validationErr != nil {
			return domain.SegL2{}, fmt.Errorf("schema validation failed for %s: %w", filePath, validationErr)
		}

		var segL2 domain.SegL2
		if err = yaml.Unmarshal(data, &segL2); err != nil {
			return domain.SegL2{}, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
		}
		segL2.SetDefaults()
		return segL2, nil
	default:
		return domain.SegL2{}, fmt.Errorf("unsupported version %s in file %s", fileVersion.Version, filePath)
	}
}
