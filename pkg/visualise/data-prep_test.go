package visualise

import (
	"testing"

	tx "github.com/kvql/bunsceal/pkg/taxonomy"
)

func TestValidateRows(t *testing.T) {
	rowsMap = map[int][]string{
		0: {"production", "ci", "sandbox", "staging", "dev"},
		1: {"shared-service"},
	}
	tx := &tx.Taxonomy{
		SegL1s: map[string]tx.SegL1{
			"shared-service": {},
			"production":     {},
			"ci":             {},
			"sandbox":        {},
			"staging":        {},
			"de":             {},
		},
	}

	err := validateRows(tx, rowsMap)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
func TestValidateRowsPass(t *testing.T) {
	rowsMap = map[int][]string{
		0: {"production", "ci", "sandbox", "staging", "dev"},
		1: {"shared-service"},
	}
	tx := &tx.Taxonomy{
		SegL1s: map[string]tx.SegL1{
			"shared-service": {},
			"production":     {},
			"ci":             {},
			"sandbox":        {},
			"staging":        {},
			"dev":            {},
		},
	}

	err := validateRows(tx, rowsMap)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestValidateRowCount(t *testing.T) {
	rowsMap = map[int][]string{
		0: {"production", "ci", "sandbox", "staging", "dev"},
		1: {"shared-service"},
	}
	tx := &tx.Taxonomy{
		SegL1s: map[string]tx.SegL1{
			"shared-service": {},
			"production":     {},
			"ci":             {},
			"sandbox":        {},
			"staging":        {},
			"dev":            {},
			"dev2":           {},
		},
	}

	err := validateRows(tx, rowsMap)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}
