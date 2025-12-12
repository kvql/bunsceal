package testhelpers

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestFiles manages temporary test files with automatic cleanup
type TestFiles struct {
	t       *testing.T
	tmpDirs []string
}

// NewTestFiles creates a new test file manager
// Automatically registers cleanup to remove all created directories
func NewTestFiles(t *testing.T) *TestFiles {
	t.Helper()

	tf := &TestFiles{t: t}

	// Register cleanup function to remove all temp directories
	t.Cleanup(func() {
		for _, dir := range tf.tmpDirs {
			os.RemoveAll(dir)
		}
	})

	return tf
}

// SegL1Fixture represents a minimal SegL1 for test file creation
type SegL1Fixture struct {
	Name        string
	ID          string
	Sensitivity string
	Criticality string
}

// SegFixture represents a minimal Seg for test file creation
type SegFixture struct {
	Name string
	ID   string
}

// CreateSegL1Files creates a directory with SegL1 YAML files
// Returns the temporary directory path
func (tf *TestFiles) CreateSegL1Files(items []SegL1Fixture) string {
	tf.t.Helper()

	tmpDir, err := os.MkdirTemp("", "segl1-test-*")
	if err != nil {
		tf.t.Fatalf("Failed to create temp directory: %v", err)
	}
	tf.tmpDirs = append(tf.tmpDirs, tmpDir)

	template := `name: "%s"
id: "%s"
description: "This is a test environment with sufficient description length to meet minimum requirements."
sensitivity: "%s"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
criticality: "%s"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation."
`

	for i, item := range items {
		content := fmt.Sprintf(template, item.Name, item.ID, item.Sensitivity, item.Criticality)
		filename := fmt.Sprintf("env%d.yaml", i+1)
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0600); err != nil {
			tf.t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

// CreateSegFiles creates a directory with Seg YAML files
// Returns the temporary directory path
func (tf *TestFiles) CreateSegFiles(items []SegFixture) string {
	tf.t.Helper()

	tmpDir, err := os.MkdirTemp("", "Seg-test-*")
	if err != nil {
		tf.t.Fatalf("Failed to create temp directory: %v", err)
	}
	tf.tmpDirs = append(tf.tmpDirs, tmpDir)

	template := `version: "1.0"
name: "%s"
id: "%s"
description: "Test domain for validating file loading and parsing with minimum required length"
sensitivity: "A"
sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
criticality: "1"
criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement for validation purposes."
l1_parents:
  - production
l1_overrides:
  production:
    sensitivity: "A"
    sensitivity_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
    criticality: "1"
    criticality_rationale: "Test rationale with sufficient length to meet the minimum character requirement."
`

	for i, item := range items {
		content := fmt.Sprintf(template, item.Name, item.ID)
		filename := fmt.Sprintf("domain%d.yaml", i+1)
		if err := os.WriteFile(filepath.Join(tmpDir, filename), []byte(content), 0600); err != nil {
			tf.t.Fatalf("Failed to write test file %s: %v", filename, err)
		}
	}

	return tmpDir
}

// CreateYAMLFile creates a single YAML file with the given content
// Returns the file path
func (tf *TestFiles) CreateYAMLFile(prefix, content string) string {
	tf.t.Helper()

	tmpFile, err := os.CreateTemp("", prefix+"-*.yaml")
	if err != nil {
		tf.t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.WriteString(content); err != nil {
		tmpFile.Close()
		tf.t.Fatalf("Failed to write test file: %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

// CreateEmptyDir creates an empty temporary directory
// Useful for testing error conditions
func (tf *TestFiles) CreateEmptyDir() string {
	tf.t.Helper()

	tmpDir, err := os.MkdirTemp("", "empty-test-*")
	if err != nil {
		tf.t.Fatalf("Failed to create temp directory: %v", err)
	}
	tf.tmpDirs = append(tf.tmpDirs, tmpDir)

	return tmpDir
}
