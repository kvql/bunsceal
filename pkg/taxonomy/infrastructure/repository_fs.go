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

type LevelPaths map[string]string

// ConfigFsReposistory If relative path set for TaxonomyDir, config file path used as the base.
type ConfigFsReposistory struct {
	TaxonomyDir string `yaml:"taxonomy_path,omitempty"`
	L1Dir       string `yaml:"l1_dir"`
	L2Dir       string `yaml:"l2_dir"`
}

func (cfs *ConfigFsReposistory) GetLevelPath(level string) (string, error) {
	var path string
	switch level {
	case "1":
		path = filepath.Join(cfs.TaxonomyDir, cfs.L1Dir)
	case "2":
		path = filepath.Join(cfs.TaxonomyDir, cfs.L2Dir)
	default:
		return "", fmt.Errorf("no path found for level %s", level)
	}
	return path, nil
}

// FileSegRepository implements taxonomy.SegRepository for file-based data sources
type FileSegRepository struct {
	schemaValidator *schemaValidation.SchemaValidator
	config          ConfigFsReposistory
}

// NewFileSegRepository creates a new file-based Seg repository
// schemaValidator is used to validate each file against the JSON schema
func NewFileSegRepository(schemaValidator *schemaValidation.SchemaValidator, cfg ConfigFsReposistory) *FileSegRepository {
	return &FileSegRepository{
		schemaValidator: schemaValidator,
		config:          cfg,
	}
}

func (r *FileSegRepository) LoadLevel(level string) ([]domain.Seg, error) {
	var segList []domain.Seg
	var parseErrors []error

	path, err := r.config.GetLevelPath(level)
	if err != nil {
		return nil, err
	}

	err = filepath.WalkDir(path, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			seg, err := r.parseSegFile(path, level)
			if err != nil {
				o11y.Log.Printf("Error parsing file %s: %v\n", path, err)
				parseErrors = append(parseErrors, err)
				return nil // Continue walking despite parse error
			}
			segList = append(segList, seg)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	if len(parseErrors) > 0 {
		return nil, fmt.Errorf("failed to parse %d file(s) in directory: %s", len(parseErrors), path)
	}

	return segList, nil
}

func (r *FileSegRepository) parseSegFile(filePath string, level string) (domain.Seg, error) {
	// #nosec G304 -- filePath comes from config-specified taxonomy directory, not user input
	data, err := os.ReadFile(filePath)
	if err != nil {
		return domain.Seg{}, fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	if validationErr := r.schemaValidator.ValidateData(data, "seg-level.json"); validationErr != nil {
		return domain.Seg{}, fmt.Errorf("schema validation failed for %s: %w", filePath, validationErr)
	}

	var seg domain.Seg
	if err = yaml.Unmarshal(data, &seg); err != nil {
		return domain.Seg{}, fmt.Errorf("failed to unmarshal file %s: %w", filePath, err)
	}

	// PostLoad handles defaults, validation, and label parsing
	if err = seg.PostLoad(level); err != nil {
		return domain.Seg{}, fmt.Errorf("PostLoad validation failed for %s: %w", filePath, err)
	}

	return seg, nil
}
