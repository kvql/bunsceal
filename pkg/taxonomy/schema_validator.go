package taxonomy

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

// SchemaValidator provides JSON schema validation for taxonomy entities
type SchemaValidator struct {
	compiler *jsonschema.Compiler
	schemas  map[string]*jsonschema.Schema
}

// NewSchemaValidator creates and initializes a schema validator with all taxonomy schemas
func NewSchemaValidator(schemaDir string) (*SchemaValidator, error) {
	compiler := jsonschema.NewCompiler()
	// Note: Draft version is auto-detected from $schema field in each JSON schema file

	// Load all schema files
	schemaFiles := []string{
		"common.json",
		"seg-level1.json",
		"seg-level2.json",
		"env-details.json",
		"comp-req.json",
		"compliance-requirements.json",
		"taxonomy.json",
	}

	// Add all schemas to compiler - must resolve to absolute paths
	absSchemaDir, err := filepath.Abs(schemaDir)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve schema directory: %w", err)
	}

	for _, file := range schemaFiles {
		schemaPath := filepath.Join(absSchemaDir, file)
		schemaURL := fmt.Sprintf("file://%s", schemaPath)

		// Read schema file
		data, err := os.ReadFile(schemaPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read schema %s: %w", file, err)
		}

		// Parse JSON
		var schemaDoc interface{}
		if err := json.Unmarshal(data, &schemaDoc); err != nil {
			return nil, fmt.Errorf("failed to parse schema %s: %w", file, err)
		}

		if err := compiler.AddResource(schemaURL, schemaDoc); err != nil {
			return nil, fmt.Errorf("failed to add schema %s: %w", file, err)
		}
	}

	// Pre-compile schemas for performance
	schemas := make(map[string]*jsonschema.Schema)
	for _, file := range schemaFiles {
		schemaPath := filepath.Join(absSchemaDir, file)
		schemaURL := fmt.Sprintf("file://%s", schemaPath)
		schema, err := compiler.Compile(schemaURL)
		if err != nil {
			return nil, fmt.Errorf("failed to compile schema %s: %w", file, err)
		}
		schemas[file] = schema
	}

	return &SchemaValidator{
		compiler: compiler,
		schemas:  schemas,
	}, nil
}

// ValidateYAML validates YAML data against the specified schema file
func (sv *SchemaValidator) ValidateYAML(yamlData []byte, schemaFile string) error {
	// Convert YAML to JSON for validation
	var data interface{}
	if err := yaml.Unmarshal(yamlData, &data); err != nil {
		return fmt.Errorf("failed to parse YAML: %w", err)
	}

	// Convert to JSON-compatible format (yaml.v3 uses map[string]interface{} but we need proper JSON types)
	data = convertYAMLToJSON(data)

	// Get the compiled schema
	schema, ok := sv.schemas[schemaFile]
	if !ok {
		return fmt.Errorf("schema not found: %s", schemaFile)
	}

	// Validate
	if err := schema.Validate(data); err != nil {
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
	if ve, ok := err.(*jsonschema.ValidationError); ok {
		// v6 has built-in error formatting, but we can enhance it
		// Use the built-in Error() method which provides good formatting
		return fmt.Errorf("schema validation failed:\n%s", ve.Error())
	}
	return err
}
