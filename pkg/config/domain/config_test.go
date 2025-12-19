package domain

import (
	"testing"

	"github.com/kvql/bunsceal/pkg/domain"
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
		termDef  domain.TermDef
		expected string
	}{
		{
			name:     "Simple lowercase",
			termDef:  domain.TermDef{Plural: "environments"},
			expected: "environments",
		},
		{
			name:     "Capitalised word",
			termDef:  domain.TermDef{Plural: "Environments"},
			expected: "environments",
		},
		{
			name:     "Multi-word with space",
			termDef:  domain.TermDef{Plural: "Security Environments"},
			expected: "security-environments",
		},
		{
			name:     "Already kebab-case",
			termDef:  domain.TermDef{Plural: "security-domains"},
			expected: "security-domains",
		},
		{
			name:     "Multiple spaces",
			termDef:  domain.TermDef{Plural: "My Custom Zones"},
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
