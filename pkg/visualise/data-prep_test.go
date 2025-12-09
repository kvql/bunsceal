package visualise

import (
	"testing"

	configdomain "github.com/kvql/bunsceal/pkg/config/domain"
	"github.com/kvql/bunsceal/pkg/domain"
)

func TestBuildRowsMap(t *testing.T) {
	t.Run("Uses config L1Layout when provided", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{
				L1Layout: map[string][]string{
					"0": {"prod", "staging"},
					"1": {"dev", "test"},
				},
			},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod":    {},
				"staging": {},
				"dev":     {},
				"test":    {},
			},
		}

		result, err := buildRowsMap(cfg, txy)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 rows, got %d", len(result))
		}

		if len(result[0]) != 2 || result[0][0] != "prod" || result[0][1] != "staging" {
			t.Errorf("Expected row 0 to be [prod, staging], got %v", result[0])
		}

		if len(result[1]) != 2 || result[1][0] != "dev" || result[1][1] != "test" {
			t.Errorf("Expected row 1 to be [dev, test], got %v", result[1])
		}
	})

	t.Run("Adds missing L1s to last row", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{
				L1Layout: map[string][]string{
					"0": {"prod", "staging"},
				},
			},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod":           {},
				"staging":        {},
				"dev":            {},
				"test":           {},
				"shared-service": {},
			},
		}

		result, err := buildRowsMap(cfg, txy)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		if len(result) != 2 {
			t.Errorf("Expected 2 rows (1 configured + 1 for missing), got %d", len(result))
		}

		if len(result[0]) != 2 {
			t.Errorf("Expected row 0 to have 2 items, got %d", len(result[0]))
		}

		// Row 1 should contain the missing L1s
		if len(result[1]) != 3 {
			t.Errorf("Expected row 1 to have 3 missing L1s, got %d", len(result[1]))
		}

		// Check that all missing L1s are in row 1
		missingL1s := map[string]bool{"dev": false, "test": false, "shared-service": false}
		for _, l1Id := range result[1] {
			if _, ok := missingL1s[l1Id]; ok {
				missingL1s[l1Id] = true
			}
		}

		for l1Id, found := range missingL1s {
			if !found {
				t.Errorf("Expected missing L1 %s to be in last row", l1Id)
			}
		}
	})

	t.Run("Defaults to single row when no L1Layout configured", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"production": {},
				"dev":        {},
			},
		}

		result, err := buildRowsMap(cfg, txy)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Should default to all L1s on row 0
		if len(result) != 1 {
			t.Errorf("Expected 1 row (default), got %d", len(result))
		}

		if len(result[0]) != 2 {
			t.Errorf("Expected row 0 to have 2 L1s, got %d", len(result[0]))
		}
	})

	t.Run("Defaults to single row when L1Layout is nil", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{
				L1Layout: nil,
			},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"production": {},
				"staging":    {},
			},
		}

		result, err := buildRowsMap(cfg, txy)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Should default to all L1s on row 0
		if len(result) != 1 {
			t.Errorf("Expected 1 row (default), got %d", len(result))
		}

		if len(result[0]) != 2 {
			t.Errorf("Expected row 0 to have 2 L1s, got %d", len(result[0]))
		}
	})

	t.Run("Config has more l1 than taxonomy", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{
				L1Layout: map[string][]string{
					"0": {"prod", "staging"},
					"1": {"dev", "unknown"},
				},
			},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod":    {},
				"staging": {},
				"dev":     {},
			},
		}

		_, err := buildRowsMap(cfg, txy)

		// Taxonomy must be source of truth, error out if config includes l1 ids which don't exist
		if err == nil {
			t.Error("Expected error when config includes L1 IDs that don't exist in taxonomy")
		}
	})

	t.Run("Handles non-sequential row numbers", func(t *testing.T) {
		cfg := &configdomain.Config{
			Visuals: configdomain.VisualsDef{
				L1Layout: map[string][]string{
					"0": {"prod"},
					"2": {"staging"},
				},
			},
		}

		txy := &domain.Taxonomy{
			SegL1s: map[string]domain.SegL1{
				"prod":    {},
				"staging": {},
				"dev":     {},
			},
		}

		result, err := buildRowsMap(cfg, txy)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		// Should have rows 0, 2, and 3 (for missing L1s after max row 2)
		if _, ok := result[0]; !ok {
			t.Error("Expected row 0 to exist")
		}
		if _, ok := result[2]; !ok {
			t.Error("Expected row 2 to exist")
		}
		if _, ok := result[3]; !ok {
			t.Error("Expected row 3 to exist for missing L1s")
		}

		// Row 3 should contain dev
		found := false
		for _, l1Id := range result[3] {
			if l1Id == "dev" {
				found = true
				break
			}
		}
		if !found {
			t.Error("Expected 'dev' to be in row 3 (last row for missing L1s)")
		}
	})
}
