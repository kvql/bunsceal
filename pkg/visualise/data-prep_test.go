package visualise

import (
	"testing"

	tx "github.com/kvql/bunsceal/pkg/taxonomy"
)

func TestValidateRows(t *testing.T) {
	rowsMap = map[int][]string{
		0: []string{"production", "ci", "sandbox", "staging", "dev"},
		1: []string{"shared-service"},
	}
	tx := &tx.Taxonomy{
		SecEnvironments: map[string]tx.SecEnv{
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
		0: []string{"production", "ci", "sandbox", "staging", "dev"},
		1: []string{"shared-service"},
	}
	tx := &tx.Taxonomy{
		SecEnvironments: map[string]tx.SecEnv{
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
		0: []string{"production", "ci", "sandbox", "staging", "dev"},
		1: []string{"shared-service"},
	}
	tx := &tx.Taxonomy{
		SecEnvironments: map[string]tx.SecEnv{
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
