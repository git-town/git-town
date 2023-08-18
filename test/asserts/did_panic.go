package asserts

import "testing"

// Paniced verifies in tests that the unit test resulted in a panic.
func Paniced(t *testing.T) {
	t.Helper()
	if r := recover(); r == nil {
		t.Errorf("The code did not panic")
	}
}
