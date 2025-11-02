package testhelpers

import (
	"testing"
)

// AssertValidationPasses fails the test if validation did not pass
func AssertValidationPasses(t *testing.T, valid bool, context string) {
	t.Helper()
	if !valid {
		t.Errorf("%s: expected validation to pass", context)
	}
}

// AssertValidationFails fails the test if validation passed when it should have failed
func AssertValidationFails(t *testing.T, valid bool, context string) {
	t.Helper()
	if valid {
		t.Errorf("%s: expected validation to fail", context)
	}
}

// AssertFailureCount fails the test if failure count doesn't match expected
func AssertFailureCount(t *testing.T, failures, expected int, context string) {
	t.Helper()
	if failures != expected {
		t.Errorf("%s: expected %d failures, got %d", context, expected, failures)
	}
}

// AssertMinFailures fails the test if failure count is less than minimum
func AssertMinFailures(t *testing.T, failures, minimum int, context string) {
	t.Helper()
	if failures < minimum {
		t.Errorf("%s: expected at least %d failures, got %d", context, minimum, failures)
	}
}
