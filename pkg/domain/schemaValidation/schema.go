package schemaValidation

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

const SchemaBaseURL = "https://github.com/kvql/bunsceal/pkg/domain/schemas/"

// SchemaValidator provides JSON schema validation for taxonomy entities
type SchemaValidator struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
}

// ExternalSchema represents a schema to be registered from a JSON string
type ExternalSchema struct {
	JSON string // JSON schema content
	ID   string // Schema $id URL
}

// NewSchemaValidator creates and initialises a schema validator with all taxonomy schemas
// externalSchemas are registered before compilation to support $ref resolution
func NewSchemaValidator(schemaDir string, schemaBaseURL string, externalSchemas ...ExternalSchema) (*SchemaValidator, error) {
	compiler := jsonschema.NewCompiler()
	// Note: Draft version is auto-detected from $schema field in each JSON schema file
	info, err := os.Stat(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("error %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("error schema path not dir: %w", err)
	}
	// Load all schema files
	schemaFiles, err := os.ReadDir(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("error %w", err)
	}
	// Add all schemas to compiler - must resolve to absolute paths
	absSchemaDir, err := filepath.Abs(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve schema directory: %w", err)
	}

	for _, file := range schemaFiles {
		schemaPath := filepath.Join(absSchemaDir, file.Name())

		// Read schema file
		// #nosec G304 -- schemaPath is constructed from trusted schema directory and known schema files
		data, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema %s: %w", file, err)
		}

		// Parse JSON
		var schemaDoc interface{}
		if err := json.Unmarshal(data, &schemaDoc); err != nil {
			return nil, fmt.Errorf("failed to parse schema %s: %w", file, err)
		}

		// Register schema using its $id if present, otherwise use file URL
		// This allows relative $ref paths in schemas to resolve correctly
		if schemaMap, ok := schemaDoc.(map[string]interface{}); ok {
			if id, ok := schemaMap["$id"].(string); ok {
				// Register with the $id URL so relative refs work
				if err := compiler.AddResource(id, schemaDoc); err != nil {
					return nil, fmt.Errorf("failed to add schema %s: %w", file, err)
				}
			} else {
				// Fallback to file URL if no $id
				schemaURL := fmt.Sprintf("file://%s", schemaPath)
				if err := compiler.AddResource(schemaURL, schemaDoc); err != nil {
					return nil, fmt.Errorf("failed to add schema %s: %w", file, err)
				}
			}
		}
	}

	// Register external schemas before compilation (needed for $ref resolution)
	for _, ext := range externalSchemas {
		var schemaDoc interface{}
		if err := json.Unmarshal([]byte(ext.JSON), &schemaDoc); err != nil {
			return nil, fmt.Errorf("failed to parse external schema %s: %w", ext.ID, err)
		}
		if err := compiler.AddResource(ext.ID, schemaDoc); err != nil {
			return nil, fmt.Errorf("failed to add external schema %s: %w", ext.ID, err)
		}
	}

	// Pre-compile schemas for performance
	// Compile using the base schema directory URL for resolution
	schemas := make(map[string]*jsonschema.Schema)

	for _, file := range schemaFiles {
		schemaURL := schemaBaseURL + file.Name()
 		schema, err := compiler.Compile(schemaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to compile schema %s: %w", file, err)
		}
		schemas[file.Name()] = schema
	}

	return &SchemaValidator{
		compiler: compiler,
		schemas:  schemas,
	}, nil
}

// ValidateData validates data (YAML/JSON) against the specified schema file
// Accepts raw bytes in YAML or JSON format and validates against JSON Schema
func (sv *SchemaValidator) ValidateData(data []byte, schemaFile string) error {
	// Parse data (supports both YAML and JSON via yaml.v3)
	var parsedData interface{}
	if err := yaml.Unmarshal(data, &parsedData); err != nil {
		return fmt.Errorf("failed to parse data: %w", err)
	}

	// Convert to JSON-compatible format (yaml.v3 uses map[string]interface{} but we need proper JSON types)
	parsedData = convertYAMLToJSON(parsedData)

	// Get the compiled schema
	schema, ok := sv.schemas[schemaFile]
	if !ok {
		return fmt.Errorf("schema not found: %s", schemaFile)
	}

	// Validate
	if err := schema.Validate(parsedData); err != nil {
		return formatValidationError(err)
	}

	return nil
}

// convertYAMLToJSON converts YAML-parsed data to JSON-compatible types
// yaml.v3 can produce map[interface{}]interface{} which JSON doesn't support
func convertYAMLToJSON(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range x {
			m[fmt.Sprintf("%v", k)] = convertYAMLToJSON(v)
		}
		return m
	case []interface{}:
		for i, v := range x {
			x[i] = convertYAMLToJSON(v)
		}
	}
	return i
}

// formatValidationError formats jsonschema validation errors in a readable way
func formatValidationError(err error) error {
	var ve *jsonschema.ValidationError
	if errors.As(err, &ve) {
		// v6 has built-in error formatting, but we can enhance it
		// Use the built-in Error() method which provides good formatting
		return fmt.Errorf("schema validation failed:\n%s", ve.Error())
	}
	return err
}
