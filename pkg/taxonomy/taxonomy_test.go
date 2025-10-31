package taxonomy

import (
	"testing"
)

func TestTaxonomy_ValidateSharedServices(t *testing.T) {
	t.Run("Shared service environment not found", func(t *testing.T) {
		txy := &Taxonomy{
			SegL1s: make(map[string]SegL1),
		}
		valid, _ := txy.ValidateSharedServices()
		if valid {
			t.Error("Expected valid to be false")
		}
	})

	t.Run("Shared service environment with incorrect sensitivity, criticality and comp requirements", func(t *testing.T) {
		txy := &Taxonomy{
			SegL1s: map[string]SegL1{
				"shared-service": {
					Sensitivity:    "b",
					Criticality:    "2",
					ComplianceReqs: []string{"req1", "req2"},
				},
			},
			CompReqs: map[string]CompReq{
				"req1": {},
				"req2": {},
				"req3": {},
			},
		}
		valid, failures := txy.ValidateSharedServices()
		if valid {
			t.Error("Expected valid to be false")
		}
		if failures != 2 {
			t.Errorf("Expected failures to be 2, got %d", failures)
		}
	})

	t.Run("Valid shared service environment", func(t *testing.T) {
		txy := &Taxonomy{
			SegL1s: map[string]SegL1{
				"shared-service": {
					Sensitivity:    SenseOrder[0],
					Criticality:    CritOrder[0],
					ComplianceReqs: []string{"req1", "req2"},
				},
			},
			CompReqs: map[string]CompReq{
				"req1": {},
				"req2": {},
			},
		}
		valid, failures := txy.ValidateSharedServices()
		if !valid {
			t.Error("Expected valid to be true")
		}
		if failures != 0 {
			t.Errorf("Expected failures to be 0, got %d", failures)
		}
	})
}
