package list_test

import "testing"

func TestList(t *testing.T) {
	t.Parallel()
	t.Run("MoveCursorDown", func(t *testing.T) {
		t.Parallel()
		t.Run("entry above is enabled", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("first and second entry above are disabled, third entry above is enabled", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("already at beginning of the list", func(t *testing.T) {
			t.Parallel()
		})
		t.Run("all entries above are disabled", func(t *testing.T) {
			t.Parallel()
		})
	})
}
