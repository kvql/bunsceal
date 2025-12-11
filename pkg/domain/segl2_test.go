package domain

import (
	"testing"
)

func TestSegL2_ValidateL1Consistency(t *testing.T) {
	tests := []struct {
		name        string
		segL2       SegL2
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid - all override keys in parents",
			segL2: SegL2{
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
			segL2: SegL2{
				L1Parents: []string{"prod", "staging", "dev"},
				L1Overrides: map[string]L1Overrides{
					"prod": {},
				},
			},
			expectError: false,
		},
		{
			name: "Valid - empty overrides",
			segL2: SegL2{
				L1Parents:   []string{"prod"},
				L1Overrides: map[string]L1Overrides{},
			},
			expectError: false,
		},
		{
			name: "Valid - nil overrides",
			segL2: SegL2{
				L1Parents:   []string{"prod"},
				L1Overrides: nil,
			},
			expectError: false,
		},
		{
			name: "Invalid - override key not in parents",
			segL2: SegL2{
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
			segL2: SegL2{
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
			segL2: SegL2{
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
			err := tt.segL2.ValidateL1Consistency()

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
