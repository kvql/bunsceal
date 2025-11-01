package taxonomy

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	// Verify defaults are set (actual values defined in DefaultConfig)
	if cfg.Terminology.L1.Singular == "" {
		t.Error("Expected L1 singular to be set")
	}
	if cfg.Terminology.L1.Plural == "" {
		t.Error("Expected L1 plural to be set")
	}
	if cfg.Terminology.L2.Singular == "" {
		t.Error("Expected L2 singular to be set")
	}
	if cfg.Terminology.L2.Plural == "" {
		t.Error("Expected L2 plural to be set")
	}
}

func TestTermDef_DirName(t *testing.T) {
	tests := []struct {
		name     string
		termDef  TermDef
		expected string
	}{
		{
			name:     "Simple lowercase",
			termDef:  TermDef{Plural: "environments"},
			expected: "environments",
		},
		{
			name:     "Capitalized word",
			termDef:  TermDef{Plural: "Environments"},
			expected: "environments",
		},
		{
			name:     "Multi-word with space",
			termDef:  TermDef{Plural: "Security Environments"},
			expected: "security-environments",
		},
		{
			name:     "Already kebab-case",
			termDef:  TermDef{Plural: "security-domains"},
			expected: "security-domains",
		},
		{
			name:     "Multiple spaces",
			termDef:  TermDef{Plural: "My Custom Zones"},
			expected: "my-custom-zones",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.termDef.DirName()
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	t.Run("Returns defaults when config file missing", func(t *testing.T) {
		defaults := DefaultConfig()

		// Use non-existent path
		cfg, err := LoadConfig("", "/nonexistent/path")
		if err != nil {
			t.Errorf("Expected no error for missing config, got: %v", err)
		}

		// Should return defaults
		if cfg.Terminology.L1.Singular != defaults.Terminology.L1.Singular {
			t.Errorf("Expected default L1 singular '%s', got '%s'", defaults.Terminology.L1.Singular, cfg.Terminology.L1.Singular)
		}
		if cfg.Terminology.L2.Singular != defaults.Terminology.L2.Singular {
			t.Errorf("Expected default L2 singular '%s', got '%s'", defaults.Terminology.L2.Singular, cfg.Terminology.L2.Singular)
		}
	})
}

func TestLoadConfig_CompleteConfig(t *testing.T) {
	t.Run("Loads complete custom config", func(t *testing.T) {
		// Create temporary config file
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Zone"
    plural: "Zones"
  l2:
    singular: "Application"
    plural: "Applications"
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		cfg, err := LoadConfig(configPath, "")
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}

		if cfg.Terminology.L1.Singular != "Zone" {
			t.Errorf("Expected L1 singular 'Zone', got '%s'", cfg.Terminology.L1.Singular)
		}
		if cfg.Terminology.L1.Plural != "Zones" {
			t.Errorf("Expected L1 plural 'Zones', got '%s'", cfg.Terminology.L1.Plural)
		}
		if cfg.Terminology.L2.Singular != "Application" {
			t.Errorf("Expected L2 singular 'Application', got '%s'", cfg.Terminology.L2.Singular)
		}
		if cfg.Terminology.L2.Plural != "Applications" {
			t.Errorf("Expected L2 plural 'Applications', got '%s'", cfg.Terminology.L2.Plural)
		}
	})
}

func TestLoadConfig_PartialL1Config(t *testing.T) {
	t.Run("Falls back to defaults when L1 missing singular", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    plural: "Zones"
  l2:
    singular: "Application"
    plural: "Applications"
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configPath, "")
		if err == nil {
			t.Errorf("Expected error on load for not defining singular value")
		}
	})

	t.Run("Falls back to defaults when L1 missing plural", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Zone"
  l2:
    singular: "Application"
    plural: "Applications"
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configPath, "")
		if err == nil {
			t.Errorf("Expected error for not defining plural values, go no error")
		}
	})
}

func TestLoadConfig_PartialL2Config(t *testing.T) {
	t.Run("Falls back to defaults when L2 incomplete", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Zone"
    plural: "Zones"
  l2:
    singular: "Application"
    # Missing plural
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configPath, "")
		if err == nil {
			t.Errorf("Expected error on load for not defining singular value")
		}
	})
}

func TestLoadConfig_L2DefinedButBlank(t *testing.T) {
	t.Run("Fail when blank terms defined", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Zone"
    plural: "Zones"
  l2:
    singular: ""
    plural: ""
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configPath, "")
		if err == nil {
			t.Errorf("Expected failed load based on schema validation for blank terms")
		}
	})
}

func TestLoadConfig_DefaultLocation(t *testing.T) {
	t.Run("Loads from default location when no explicit path", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Custom"
    plural: "Customs"
  l2:
    singular: "Service"
    plural: "Services"
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		// Pass empty string for configPath, tmpDir as taxDir
		cfg, err := LoadConfig("", tmpDir)
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}

		if cfg.Terminology.L1.Singular != "Custom" {
			t.Errorf("Expected L1 singular 'Custom', got '%s'", cfg.Terminology.L1.Singular)
		}
		if cfg.Terminology.L2.Singular != "Service" {
			t.Errorf("Expected L2 singular 'Service', got '%s'", cfg.Terminology.L2.Singular)
		}
	})
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	t.Run("Returns error for invalid YAML", func(t *testing.T) {
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		invalidYAML := `terminology:
  l1:
    singular: "Zone
    # Missing closing quote - invalid YAML
`
		if err := os.WriteFile(configPath, []byte(invalidYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		_, err := LoadConfig(configPath, "")
		if err == nil {
			t.Error("Expected error for invalid YAML, got nil")
		}
	})
}
func TestLoadConfig_missingLevel(t *testing.T) {
	t.Run("Uses defaults for L2 when only L1 defined", func(t *testing.T) {
		defaults := DefaultConfig()
		tmpDir := t.TempDir()
		configPath := filepath.Join(tmpDir, "config.yaml")

		configYAML := `terminology:
  l1:
    singular: "Zone"
    plural: "Zones"
`
		if err := os.WriteFile(configPath, []byte(configYAML), 0644); err != nil {
			t.Fatalf("Failed to write config file: %v", err)
		}

		cfg, err := LoadConfig(configPath, "")
		if err != nil {
			t.Fatalf("Expected successful load, got error: %v", err)
		}
		if cfg.Terminology.L2.Singular != defaults.Terminology.L2.Singular {
			t.Errorf("Expected default L2 singular '%s', got '%s'", defaults.Terminology.L2.Singular, cfg.Terminology.L2.Singular)
		}
	})
}
