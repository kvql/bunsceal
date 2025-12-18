package domain

import (
	"strings"
	"testing"
)

// Test overview
// - Label parsing
// - l1 override conn

func TestSegL1_ParseLabels(t *testing.T) {
	t.Run("Parses valid key:value pairs", func(t *testing.T) {
		seg := &Seg{
			Labels: []string{"env:prod", "region:us-east-1"},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(seg.ParsedLabels) != 2 {
			t.Errorf("Expected 2 labels, got %d", len(seg.ParsedLabels))
		}
		if seg.ParsedLabels["env"] != "prod" {
			t.Errorf("Expected env=prod, got %s", seg.ParsedLabels["env"])
		}
		if seg.ParsedLabels["region"] != "us-east-1" {
			t.Errorf("Expected region=us-east-1, got %s", seg.ParsedLabels["region"])
		}
	})

	t.Run("Handles values containing colons", func(t *testing.T) {
		seg := &Seg{
			Labels: []string{"url:https://example.com:8080"},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if seg.ParsedLabels["url"] != "https://example.com:8080" {
			t.Errorf("Expected url=https://example.com:8080, got %s", seg.ParsedLabels["url"])
		}
	})

	t.Run("Empty labels slice succeeds", func(t *testing.T) {
		seg := &Seg{
			Labels: []string{},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if len(seg.ParsedLabels) != 0 {
			t.Errorf("Expected empty map, got %d entries", len(seg.ParsedLabels))
		}
	})

	t.Run("Initialises nil ParsedLabels map", func(t *testing.T) {
		seg := &Seg{
			Labels:       []string{"key:value"},
			ParsedLabels: nil,
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		if seg.ParsedLabels == nil {
			t.Fatal("Expected ParsedLabels to be initialised")
		}
		if seg.ParsedLabels["key"] != "value" {
			t.Errorf("Expected key=value, got %s", seg.ParsedLabels["key"])
		}
	})

	t.Run("Returns error for invalid format", func(t *testing.T) {
		seg := &Seg{
			Labels: []string{"invalid-no-colon"},
		}

		err := seg.ParseLabels()
		if err == nil {
			t.Fatal("Expected error for invalid format, got nil")
		}

		if !strings.Contains(err.Error(), "invalid label format") {
			t.Errorf("Expected error message to contain 'invalid label format', got %s", err.Error())
		}
	})

	t.Run("Handles AWS-compliant special characters in values", func(t *testing.T) {
		testCases := []struct {
			name     string
			labels   []string
			expected map[string]string
		}{
			{
				"Value with spaces",
				[]string{"env:production environment"},
				map[string]string{"env": "production environment"},
			},
			{
				"Value with plus sign",
				[]string{"version:v1.0.0+build.123"},
				map[string]string{"version": "v1.0.0+build.123"},
			},
			{
				"Value with equals",
				[]string{"formula:x=y+z"},
				map[string]string{"formula": "x=y+z"},
			},
			{
				"Value with at sign",
				[]string{"owner:team@example.com"},
				map[string]string{"owner": "team@example.com"},
			},
			{
				"Combined special characters",
				[]string{"deploy:user@host:/path v1.0+patch"},
				map[string]string{"deploy": "user@host:/path v1.0+patch"},
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				seg := &Seg{Labels: tc.labels}
				err := seg.ParseLabels()
				if err != nil {
					t.Fatalf("Expected no error, got %v", err)
				}

				for key, expectedValue := range tc.expected {
					if seg.ParsedLabels[key] != expectedValue {
						t.Errorf("Expected %s=%s, got %s", key, expectedValue, seg.ParsedLabels[key])
					}
				}
			})
		}
	})
}

func TestSeg_ValidateL1Consistency(t *testing.T) {
	tests := []struct {
		name        string
		Seg         Seg
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid - all override keys in parents",
			Seg: Seg{
				L1Parents: []string{"prod", "staging"},
				L1Overrides: map[string]L1Overrides{
					"prod":    {},
					"staging": {},
				},
			},
			expectError: false,
		},
		{
			name: "Valid - subset of parents have overrides",
			Seg: Seg{
				L1Parents: []string{"prod", "staging", "dev"},
				L1Overrides: map[string]L1Overrides{
					"prod": {},
				},
			},
			expectError: false,
		},
		{
			name: "Valid - empty overrides",
			Seg: Seg{
				L1Parents:   []string{"prod"},
				L1Overrides: map[string]L1Overrides{},
			},
			expectError: false,
		},
		{
			name: "Valid - nil overrides",
			Seg: Seg{
				L1Parents:   []string{"prod"},
				L1Overrides: nil,
			},
			expectError: false,
		},
		{
			name: "Invalid - override key not in parents",
			Seg: Seg{
				L1Parents: []string{"prod"},
				L1Overrides: map[string]L1Overrides{
					"prod":    {},
					"staging": {},
				},
			},
			expectError: true,
			errorMsg:    "l1_overrides contains key 'staging' which is not in l1_parents",
		},
		{
			name: "Invalid - multiple override keys not in parents",
			Seg: Seg{
				L1Parents: []string{"prod"},
				L1Overrides: map[string]L1Overrides{
					"staging": {},
					"dev":     {},
				},
			},
			expectError: true,
		},
		{
			name: "Invalid - no parents but has overrides",
			Seg: Seg{
				L1Parents: []string{},
				L1Overrides: map[string]L1Overrides{
					"prod": {},
				},
			},
			expectError: true,
			errorMsg:    "l1_overrides contains key 'prod' which is not in l1_parents",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.Seg.ValidateL1Consistency()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("Expected error message '%s' but got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got: %v", err)
				}
			}
		})
	}
}

func TestPostLoad_L1Segment(t *testing.T) {
	t.Run("Valid L1 segment passes validation", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Description:          "Test description",
			Sensitivity:          "A",
			SensitivityRationale: "Test rationale with sufficient length",
			Criticality:          "1",
			CriticalityRationale: "Test rationale with sufficient length",
		}

		err := seg.PostLoad("1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Level != "1" {
			t.Errorf("Expected Level='1', got: %s", seg.Level)
		}
		if seg.Prominence != 1 {
			t.Errorf("Expected Prominence=1 (default), got: %d", seg.Prominence)
		}
	})

	t.Run("Missing Criticality fails validation", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test rationale",
			// Missing Criticality
			CriticalityRationale: "Test rationale",
		}

		err := seg.PostLoad("1")
		if err == nil {
			t.Error("Expected error for missing Criticality")
		}
		if !strings.Contains(err.Error(), "Criticality") {
			t.Errorf("Expected error about Criticality, got: %v", err)
		}
	})

	t.Run("Missing CriticalityRationale fails validation", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test rationale",
			Criticality:          "1",
			// Missing CriticalityRationale
		}

		err := seg.PostLoad("1")
		if err == nil {
			t.Error("Expected error for missing CriticalityRationale")
		}
		if !strings.Contains(err.Error(), "CriticalityRationale") {
			t.Errorf("Expected error about CriticalityRationale, got: %v", err)
		}
	})

	t.Run("Missing Sensitivity fails validation", func(t *testing.T) {
		seg := Seg{
			ID:   "prod",
			Name: "Production",
			// Missing Sensitivity
			SensitivityRationale: "Test rationale",
			Criticality:          "1",
			CriticalityRationale: "Test rationale",
		}

		err := seg.PostLoad("1")
		if err == nil {
			t.Error("Expected error for missing Sensitivity")
		}
		if !strings.Contains(err.Error(), "Sensitivity") {
			t.Errorf("Expected error about Sensitivity, got: %v", err)
		}
	})

	t.Run("Missing SensitivityRationale fails validation", func(t *testing.T) {
		seg := Seg{
			ID:          "prod",
			Name:        "Production",
			Sensitivity: "A",
			// Missing SensitivityRationale
			Criticality:          "1",
			CriticalityRationale: "Test rationale",
		}

		err := seg.PostLoad("1")
		if err == nil {
			t.Error("Expected error for missing SensitivityRationale")
		}
		if !strings.Contains(err.Error(), "SensitivityRationale") {
			t.Errorf("Expected error about SensitivityRationale, got: %v", err)
		}
	})
}

func TestPostLoad_L2Segment(t *testing.T) {
	t.Run("Valid L2 segment passes validation", func(t *testing.T) {
		seg := Seg{
			ID:          "app",
			Name:        "Application",
			Description: "Test description",
			L1Parents:   []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Sensitivity: "A",
					Criticality: "1",
				},
			},
		}

		err := seg.PostLoad("2")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Level != "2" {
			t.Errorf("Expected Level='2', got: %s", seg.Level)
		}
	})

	t.Run("Missing L1Parents fails validation", func(t *testing.T) {
		seg := Seg{
			ID:   "app",
			Name: "Application",
			// Missing L1Parents
		}

		err := seg.PostLoad("2")
		if err == nil {
			t.Error("Expected error for missing L1Parents")
		}
		if !strings.Contains(err.Error(), "L1Parents") {
			t.Errorf("Expected error about L1Parents, got: %v", err)
		}
	})

	t.Run("Empty L1Parents fails validation", func(t *testing.T) {
		seg := Seg{
			ID:        "app",
			Name:      "Application",
			L1Parents: []string{}, // Empty
		}

		err := seg.PostLoad("2")
		if err == nil {
			t.Error("Expected error for empty L1Parents")
		}
	})

	t.Run("L1Overrides inconsistency fails validation", func(t *testing.T) {
		seg := Seg{
			ID:        "app",
			Name:      "Application",
			L1Parents: []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"staging": {}, // Not in L1Parents
			},
		}

		err := seg.PostLoad("2")
		if err == nil {
			t.Error("Expected error for L1Overrides inconsistency")
		}
		if !strings.Contains(err.Error(), "l1_overrides") {
			t.Errorf("Expected error about l1_overrides, got: %v", err)
		}
	})

	t.Run("L2 with labels parses correctly", func(t *testing.T) {
		seg := Seg{
			ID:          "app",
			Name:        "Application",
			Description: "Test",
			L1Parents:   []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"prod": {Sensitivity: "A", Criticality: "1"},
			},
			Labels: []string{"env:test", "team:platform"},
		}

		err := seg.PostLoad("2")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if len(seg.ParsedLabels) != 2 {
			t.Errorf("Expected 2 parsed labels, got: %d", len(seg.ParsedLabels))
		}
		if seg.ParsedLabels["env"] != "test" {
			t.Errorf("Expected env=test, got: %s", seg.ParsedLabels["env"])
		}
	})
}

func TestPostLoad_LevelAssignment(t *testing.T) {
	t.Run("Sets Level when empty", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test rationale",
			Criticality:          "1",
			CriticalityRationale: "Test rationale",
		}

		err := seg.PostLoad("1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Level != "1" {
			t.Errorf("Expected Level='1', got: %s", seg.Level)
		}
	})

	t.Run("Preserves existing Level", func(t *testing.T) {
		seg := Seg{
			Level:                "1",
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test rationale",
			Criticality:          "1",
			CriticalityRationale: "Test rationale",
		}

		err := seg.PostLoad("1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Level != "1" {
			t.Errorf("Expected Level='1', got: %s", seg.Level)
		}
	})

	t.Run("Unsupported level fails validation", func(t *testing.T) {
		seg := Seg{
			ID:   "test",
			Name: "Test",
		}

		err := seg.PostLoad("99")
		if err == nil {
			t.Error("Expected error for unsupported level")
		}
		if !strings.Contains(err.Error(), "unsupported segment level") {
			t.Errorf("Expected error about unsupported level, got: %v", err)
		}
	})
}

func TestPostLoad_SetDefaults(t *testing.T) {
	t.Run("Sets default prominence", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test",
			Criticality:          "1",
			CriticalityRationale: "Test",
			// Prominence not set, should default to 1
		}

		err := seg.PostLoad("1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Prominence != 1 {
			t.Errorf("Expected Prominence=1, got: %d", seg.Prominence)
		}
	})

	t.Run("Preserves non-zero prominence", func(t *testing.T) {
		seg := Seg{
			ID:                   "prod",
			Name:                 "Production",
			Sensitivity:          "A",
			SensitivityRationale: "Test",
			Criticality:          "1",
			CriticalityRationale: "Test",
			Prominence:           5,
		}

		err := seg.PostLoad("1")
		if err != nil {
			t.Errorf("Expected no error, got: %v", err)
		}
		if seg.Prominence != 5 {
			t.Errorf("Expected Prominence=5, got: %d", seg.Prominence)
		}
	})
}

func TestParseLabels_Overrides(t *testing.T) {
	t.Run("Parses valid override labels into ParsedLabels map", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Labels: []string{
						"example.plugin/key1:value1",
						"example.plugin/key2:value2",
					},
				},
			},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		override := seg.L1Overrides["prod"]
		if len(override.ParsedLabels) != 2 {
			t.Errorf("Expected 2 parsed labels, got: %d", len(override.ParsedLabels))
		}
		if override.ParsedLabels["example.plugin/key1"] != "value1" {
			t.Errorf("Expected key1=value1, got: %s", override.ParsedLabels["example.plugin/key1"])
		}
		if override.ParsedLabels["example.plugin/key2"] != "value2" {
			t.Errorf("Expected key2=value2, got: %s", override.ParsedLabels["example.plugin/key2"])
		}
	})

	t.Run("Succeeds with empty override labels", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Labels: []string{},
				},
			},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error for empty labels, got: %v", err)
		}
	})

	t.Run("Succeeds with no override labels field", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					// Labels field not set
				},
			},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error when labels field missing, got: %v", err)
		}
	})

	t.Run("Fails on invalid override label format", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Labels: []string{
						"invalid-no-colon",
					},
				},
			},
		}

		err := seg.ParseLabels()
		if err == nil {
			t.Fatal("Expected error for invalid format, got nil")
		}
		if !strings.Contains(err.Error(), "invalid label format") {
			t.Errorf("Expected error about invalid format, got: %v", err)
		}
	})

	t.Run("Works with multiple parent overrides", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Labels: []string{
						"example.plugin/env:production",
						"example.plugin/tier:high",
					},
				},
				"staging": {
					Labels: []string{
						"example.plugin/env:staging",
						"example.plugin/tier:low",
					},
				},
			},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error with multiple parents, got: %v", err)
		}

		prodOverride := seg.L1Overrides["prod"]
		stagingOverride := seg.L1Overrides["staging"]

		if prodOverride.ParsedLabels["example.plugin/env"] != "production" {
			t.Error("Expected prod override to have env=production")
		}
		if stagingOverride.ParsedLabels["example.plugin/env"] != "staging" {
			t.Error("Expected staging override to have env=staging")
		}
	})

	t.Run("Handles override values containing colons", func(t *testing.T) {
		seg := Seg{
			ID: "test-seg",
			L1Overrides: map[string]L1Overrides{
				"prod": {
					Labels: []string{
						"example.plugin/url:https://example.com:8080",
					},
				},
			},
		}

		err := seg.ParseLabels()
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		override := seg.L1Overrides["prod"]
		if override.ParsedLabels["example.plugin/url"] != "https://example.com:8080" {
			t.Errorf("Expected url to preserve colons in value, got: %s", override.ParsedLabels["example.plugin/url"])
		}
	})
}
