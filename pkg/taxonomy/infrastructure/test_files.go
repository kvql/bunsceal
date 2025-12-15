package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
	"gopkg.in/yaml.v3"
)

// TestFiles manages temporary test files with automatic cleanup
// Used for testing file-based repositories
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
			if err := os.RemoveAll(dir); err != nil {
				t.Logf("Failed to remove temp directory %s: %v", dir, err)
			}
		}
	})

	return tf
}

// CreateSegFiles creates a directory with Seg YAML files from domain.Seg instances
// Returns the temporary directory path
// Infrastructure layer only handles file I/O - domain knowledge comes from domain package
func (tf *TestFiles) CreateSegFiles(segments []domain.Seg) string {
	tf.t.Helper()

	tmpDir, err := os.MkdirTemp("", "seg-test-*")
	if err != nil {
		tf.t.Fatalf("Failed to create temp directory: %v", err)
	}
	tf.tmpDirs = append(tf.tmpDirs, tmpDir)

	for i, seg := range segments {
		// Determine version based on segment level or presence of L1Overrides
		version := "1.0" // Default to L2 file format
		if seg.Level == "1" || (len(seg.L1Overrides) == 0 && len(seg.L1Parents) == 0) {
			// L1 segments don't have version field in their YAML
			version = ""
		}

		// Marshal domain.Seg to YAML using its yaml tags
		data, err := yaml.Marshal(seg)
		if err != nil {
			tf.t.Fatalf("Failed to marshal segment to YAML: %v", err)
		}

		// Prepend version if this is L2 format
		var content []byte
		if version != "" {
			content = append([]byte("version: \""+version+"\"\n"), data...)
		} else {
			content = data
		}

		// Use index-based filename to allow testing duplicate IDs
		// (different files can contain segments with the same ID)
		filename := filepath.Join(tmpDir, fmt.Sprintf("seg-%d.yaml", i))
		if err := os.WriteFile(filename, content, 0600); err != nil {
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
	defer func() {
		if err := tmpFile.Close(); err != nil {
			tf.t.Logf("Failed to close test file: %v", err)
		}
	}()

	if _, err := tmpFile.WriteString(content); err != nil {
		tf.t.Fatalf("Failed to write test file: %v", err)
	}

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
