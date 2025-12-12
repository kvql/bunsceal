package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kvql/bunsceal/pkg/domain"
	"github.com/kvql/bunsceal/pkg/domain/schemaValidation"
	"github.com/kvql/bunsceal/pkg/o11y"
	"gopkg.in/yaml.v3"
)

// FileSegL1Repository implements taxonomy.SegL1Repository for file-based data sources
type FileSegL1Repository struct {
	schemaValidator *schemaValidation.SchemaValidator
}

// NewFileSegL1Repository creates a new file-based SegL1 repository
// schemaValidator is used to validate each file against the JSON schema
func NewFileSegL1Repository(schemaValidator *schemaValidation.SchemaValidator) *FileSegL1Repository {
	return &FileSegL1Repository{
		schemaValidator: schemaValidator,
	}
}

// LoadAll loads all SegL1 files from the specified directory
// Returns a slice of SegL1 entities or an error if loading fails
// Performs schema validation but NOT business rule validation
func (r *FileSegL1Repository) LoadAll(segL1Dir string) ([]domain.Seg, error) {
	files, err := os.ReadDir(segL1Dir)
	if err != nil {
		return nil, err
	}

	var segList []domain.Seg
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

func (r *FileSegL1Repository) parseSegL1File(filePath string) (domain.Seg, error) {
	return parseSegL1File(filePath, r.schemaValidator)
}

type version struct {
	Version string `yaml:"version"`
}

// parseSegL1File parses a single SegL1 file with schema validation
func parseSegL1File(filePath string, schemaValidator *schemaValidation.SchemaValidator) (domain.Seg, error) {
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.Seg{}, err
	}

	if validationErr := schemaValidator.ValidateData(data, "seg-level.json"); validationErr != nil {
		return domain.Seg{}, fmt.Errorf("schema validation failed for %s: %w", filePath, validationErr)
	}

	var segL1 domain.Seg
	if err = yaml.Unmarshal(data, &segL1); err != nil {
		return domain.Seg{}, err
	}

	if err = segL1.ParseLabels(); err != nil {
		return domain.Seg{}, err
	}

	return segL1, nil
}

// FileSegRepository implements taxonomy.SegRepository for file-based data sources
type FileSegRepository struct {
	schemaValidator *schemaValidation.SchemaValidator
}

// NewFileSegRepository creates a new file-based Seg repository
// schemaValidator is used to validate each file against the JSON schema
func NewFileSegRepository(schemaValidator *schemaValidation.SchemaValidator) *FileSegRepository {
	return &FileSegRepository{
		schemaValidator: schemaValidator,
	}
}

// LoadAll loads all Seg files from the specified directory
// Returns a slice of Seg entities or an error if loading fails
// Performs schema validation but NOT business rule validation
func (r *FileSegRepository) LoadAll(SegDir string) ([]domain.Seg, error) {
	var segList []domain.Seg
	var parseErrors []error

	err := filepath.WalkDir(SegDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			Seg, err := r.parseSegFile(path)
			if err != nil {
				o11y.Log.Printf("Error parsing file %s: %v\n", path, err)
				parseErrors = append(parseErrors, err)
				return nil // Continue walking despite parse error
			}
			segList = append(segList, Seg)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(parseErrors) > 0 {
		return nil, fmt.Errorf("failed to parse %d file(s) in directory: %s", len(parseErrors), SegDir)
	}

	return segList, nil
}

func (r *FileSegRepository) parseSegFile(filePath string) (domain.Seg, error) {
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.Seg{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	var fileVersion version
	if err = yaml.Unmarshal(data, &fileVersion); err != nil {
		return domain.Seg{}, fmt.Errorf("failed to parse version from file %s: %w", filePath, err)
	}

	switch fileVersion.Version {
	case "1.0":
		if validationErr := r.schemaValidator.ValidateData(data, "seg-level.json"); validationErr != nil {
			return domain.Seg{}, fmt.Errorf("schema validation failed for %s: %w", filePath, validationErr)
		}

		var Seg domain.Seg
		if err = yaml.Unmarshal(data, &Seg); err != nil {
			return domain.Seg{}, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
		}

		Seg.SetDefaults()

		// Validate that l1_overrides keys are subset of l1_parents
		if err = Seg.ValidateL1Consistency(); err != nil {
			return domain.Seg{}, fmt.Errorf("L1 consistency validation failed for %s: %w", filePath, err)
		}

		if err = Seg.ParseLabels(); err != nil {
			return domain.Seg{}, err
		}
		return Seg, nil
	default:
		return domain.Seg{}, fmt.Errorf("unsupported version %s in file %s", fileVersion.Version, filePath)
	}
}
