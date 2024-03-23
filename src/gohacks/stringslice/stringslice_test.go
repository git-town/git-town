package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("Connect", func(t *testing.T) {
		t.Run("no element", func(t *testing.T) {
			t.Parallel()
			give := []string{}
			have := stringslice.Connect(give)
			want := ""
			must.EqOp(t, want, have)
		})

		t.Run("single element", func(t *testing.T) {
			t.Parallel()
			give := []string{"one"}
			have := stringslice.Connect(give)
			want := "\"one\""
			must.EqOp(t, want, have)
		})

		t.Run("two elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two"}
			have := stringslice.Connect(give)
			want := "\"one\" and \"two\""
			must.EqOp(t, want, have)
		})

		t.Run("three elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three"}
			have := stringslice.Connect(give)
			want := "\"one\", \"two\", and \"three\""
			must.EqOp(t, want, have)
		})

		t.Run("four elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three", "four"}
			have := stringslice.Connect(give)
			want := "\"one\", \"two\", \"three\", and \"four\""
			must.EqOp(t, want, have)
		})
	})

	t.Run("Lines", func(t *testing.T) {
		t.Parallel()
		tests := map[string][]string{
			"":                {},
			"single line":     {"single line"},
			"multiple\nlines": {"multiple", "lines"},
		}
		for give, want := range tests {
			have := stringslice.Lines(give)
			must.Eq(t, want, have)
		}
	})

	t.Run("Longest", func(t *testing.T) {
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
	})

	t.Run("SurroundEmptyWith", func(t *testing.T) {
		t.Parallel()
		give := []string{"git", "config", "perennial-branches", ""}
		have := stringslice.SurroundEmptyWith(give, `"`)
		want := []string{"git", "config", "perennial-branches", `""`}
		must.Eq(t, want, have)
	})
}
