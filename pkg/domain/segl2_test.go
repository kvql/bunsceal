package domain

import (
	"testing"
)

func TestMigrateL1Parents(t *testing.T) {
	t.Run("Migrates from L1Overrides when L1Parents empty", func(t *testing.T) {
		s := SegL2{
			ID: "test",
			L1Overrides: map[string]L1Overrides{
				"staging": {},
				"prod":    {},
			},
		}

		migrated := s.MigrateL1Parents()

		if !migrated {
			t.Error("Expected migration to occur")
		}
		if len(s.L1Parents) != 2 {
			t.Errorf("Expected 2 parents, got %d", len(s.L1Parents))
		}
		// Should be sorted alphabetically
		if s.L1Parents[0] != "prod" || s.L1Parents[1] != "staging" {
			t.Errorf("Expected sorted [prod, staging], got %v", s.L1Parents)
		}
	})

	t.Run("Does not migrate when L1Parents already populated", func(t *testing.T) {
		s := SegL2{
			ID:        "test",
			L1Parents: []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"prod":    {},
				"staging": {},
			},
		}

		migrated := s.MigrateL1Parents()

		if migrated {
			t.Error("Expected no migration when L1Parents already set")
		}
		if len(s.L1Parents) != 1 {
			t.Error("Expected L1Parents to remain unchanged")
		}
	})

	t.Run("Idempotent - safe to call multiple times", func(t *testing.T) {
		s := SegL2{
			ID: "test",
			L1Overrides: map[string]L1Overrides{
				"prod": {},
			},
		}

		s.MigrateL1Parents()
		firstResult := make([]string, len(s.L1Parents))
		copy(firstResult, s.L1Parents)

		s.MigrateL1Parents()
		secondResult := s.L1Parents

		if len(firstResult) != len(secondResult) || firstResult[0] != secondResult[0] {
			t.Error("Expected idempotent behaviour")
		}
	})

	t.Run("Returns false when L1Overrides is empty", func(t *testing.T) {
		s := SegL2{
			ID:          "test",
			L1Overrides: map[string]L1Overrides{},
		}

		migrated := s.MigrateL1Parents()

		if migrated {
			t.Error("Expected no migration when L1Overrides is empty")
		}
		if len(s.L1Parents) != 0 {
			t.Error("Expected L1Parents to remain empty")
		}
	})

	t.Run("Sorts keys for deterministic output", func(t *testing.T) {
		s := SegL2{
			ID: "test",
			L1Overrides: map[string]L1Overrides{
				"zebra": {},
				"alpha": {},
				"beta":  {},
			},
		}

		s.MigrateL1Parents()

		expected := []string{"alpha", "beta", "zebra"}
		for i, key := range expected {
			if s.L1Parents[i] != key {
				t.Errorf("Expected L1Parents[%d] = %s, got %s", i, key, s.L1Parents[i])
			}
		}
	})
}

func TestValidateL1Consistency(t *testing.T) {
	t.Run("Passes when L1Overrides subset of L1Parents", func(t *testing.T) {
		s := SegL2{
			ID:        "test",
			L1Parents: []string{"prod", "staging", "dev"},
			L1Overrides: map[string]L1Overrides{
				"prod":    {},
				"staging": {},
			},
		}

		err := s.ValidateL1Consistency()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Passes when L1Overrides equals L1Parents", func(t *testing.T) {
		s := SegL2{
			ID:        "test",
			L1Parents: []string{"prod", "staging"},
			L1Overrides: map[string]L1Overrides{
				"prod":    {},
				"staging": {},
			},
		}

		err := s.ValidateL1Consistency()
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})

	t.Run("Fails when L1Overrides has key not in L1Parents", func(t *testing.T) {
		s := SegL2{
			ID:        "test",
			L1Parents: []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"prod":    {},
				"staging": {}, // Not in L1Parents
			},
		}

		err := s.ValidateL1Consistency()
		if err == nil {
			t.Error("Expected error for inconsistent override key")
		}
	})

	t.Run("Passes when L1Parents empty (migration not complete)", func(t *testing.T) {
		s := SegL2{
			ID:        "test",
			L1Parents: []string{},
			L1Overrides: map[string]L1Overrides{
				"prod": {},
			},
		}

		err := s.ValidateL1Consistency()
		if err != nil {
			t.Error("Expected no error when L1Parents empty (graceful migration)")
		}
	})

	t.Run("Passes when both L1Parents and L1Overrides empty", func(t *testing.T) {
		s := SegL2{
			ID:          "test",
			L1Parents:   []string{},
			L1Overrides: map[string]L1Overrides{},
		}

		err := s.ValidateL1Consistency()
		if err != nil {
			t.Error("Expected no error when both empty")
		}
	})

	t.Run("Error message includes segment ID and mismatched key", func(t *testing.T) {
		s := SegL2{
			ID:        "my-segment",
			L1Parents: []string{"prod"},
			L1Overrides: map[string]L1Overrides{
				"invalid-key": {},
			},
		}

		err := s.ValidateL1Consistency()
		if err == nil {
			t.Fatal("Expected error")
		}

		errMsg := err.Error()
		if !contains(errMsg, "my-segment") {
			t.Errorf("Error message should contain segment ID 'my-segment', got: %s", errMsg)
		}
		if !contains(errMsg, "invalid-key") {
			t.Errorf("Error message should contain key 'invalid-key', got: %s", errMsg)
		}
	})
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > 0 && (s[:len(substr)] == substr || contains(s[1:], substr))))
}
