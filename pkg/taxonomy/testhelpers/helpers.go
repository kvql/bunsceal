package testhelpers

import (
	"reflect"
	"testing"
)

// AssertNoError fails the test if err is not nil
func AssertNoError(t *testing.T, err error, context string) {
	t.Helper()
	if err != nil {
		t.Fatalf("%s: unexpected error: %v", context, err)
	}
}

// AssertError fails the test if err is nil
func AssertError(t *testing.T, err error, context string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error but got nil", context)
	}
}

// AssertEqual fails the test if got != want
func AssertEqual(t *testing.T, got, want interface{}, context string) {
	t.Helper()
	if got != want {
		t.Errorf("%s: got %v, want %v", context, got, want)
	}
}

// AssertNotEqual fails the test if got == want
func AssertNotEqual(t *testing.T, got, want interface{}, context string) {
	t.Helper()
	if got == want {
		t.Errorf("%s: got %v, want it to be different", context, got)
	}
}

// AssertMapLength fails the test if map length doesn't match expected
// Works with any map[string]T type using reflection
func AssertMapLength(t *testing.T, m interface{}, expectedLen int, context string) {
	t.Helper()

	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		t.Fatalf("%s: expected map, got %T", context, m)
	}

	actualLen := v.Len()
	if actualLen != expectedLen {
		t.Errorf("%s: expected map length %d, got %d", context, expectedLen, actualLen)
	}
}

// AssertSliceLength fails the test if slice length doesn't match expected
func AssertSliceLength(t *testing.T, slice []string, expectedLen int, context string) {
	t.Helper()
	if len(slice) != expectedLen {
		t.Errorf("%s: expected slice length %d, got %d", context, expectedLen, len(slice))
	}
}

// AssertMapContainsKey fails the test if map doesn't contain the key
// Works with any map[string]T type using reflection
func AssertMapContainsKey(t *testing.T, m interface{}, key string, context string) {
	t.Helper()

	v := reflect.ValueOf(m)
	if v.Kind() != reflect.Map {
		t.Fatalf("%s: expected map, got %T", context, m)
	}

	keyVal := reflect.ValueOf(key)
	if !v.MapIndex(keyVal).IsValid() {
		t.Errorf("%s: expected map to contain key '%s'", context, key)
	}
}

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
