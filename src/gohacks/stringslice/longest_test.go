package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestLongest(t *testing.T) {
	t.Parallel()
	t.Run("various strings", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := stringslice.Longest(give)
		must.Eq(t, 5, have)
	})
	t.Run("empty string", func(t *testing.T) {
		t.Parallel()
		give := []string{""}
		have := stringslice.Longest(give)
		must.Eq(t, 0, have)
	})
	t.Run("empty slice", func(t *testing.T) {
		t.Parallel()
		give := []string{}
		have := stringslice.Longest(give)
		must.Eq(t, 0, have)
	})
}
