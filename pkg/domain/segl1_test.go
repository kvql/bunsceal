package domain

import (
	"strings"
	"testing"
)

func TestSegL1_ParseLabels(t *testing.T) {
	t.Run("Parses valid key:value pairs", func(t *testing.T) {
		seg := &SegL1{
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
		seg := &SegL1{
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
		seg := &SegL1{
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
		seg := &SegL1{
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
		seg := &SegL1{
			Labels: []string{"invalid-no-colon"},
		}

		err := seg.ParseLabels()
		if err == nil {
			t.Fatal("Expected error for invalid format, got nil")
		}

		if !strings.Contains(err.Error(), "format invalid") {
			t.Errorf("Expected error message to contain 'format invalid', got %s", err.Error())
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
				seg := &SegL1{Labels: tc.labels}
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
