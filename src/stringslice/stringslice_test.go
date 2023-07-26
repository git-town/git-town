package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("AppendAllMissing", func(t *testing.T) {
		t.Parallel()
		list := []string{"one", "two", "three"}
		give := []string{"two", "four", "five"}
		want := []string{"one", "two", "three", "four", "five"}
		have := stringslice.AppendAllMissing(list, give)
		assert.Equal(t, want, have)
	})

	t.Run("Connect", func(t *testing.T) {
		t.Run("no element", func(t *testing.T) {
			t.Parallel()
			give := []string{}
			want := ""
			have := stringslice.Connect(give)
			assert.Equal(t, want, have)
		})

		t.Run("single element", func(t *testing.T) {
			t.Parallel()
			give := []string{"one"}
			want := "\"one\""
			have := stringslice.Connect(give)
			assert.Equal(t, want, have)
		})

		t.Run("two elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two"}
			want := "\"one\" and \"two\""
			have := stringslice.Connect(give)
			assert.Equal(t, want, have)
		})

		t.Run("three elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three"}
			want := "\"one\", \"two\", and \"three\""
			have := stringslice.Connect(give)
			assert.Equal(t, want, have)
		})

		t.Run("four elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three", "four"}
			want := "\"one\", \"two\", \"three\", and \"four\""
			have := stringslice.Connect(give)
			assert.Equal(t, want, have)
		})
	})

	t.Run("Contains", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two"}
		assert.True(t, stringslice.Contains(give, "one"))
		assert.True(t, stringslice.Contains(give, "two"))
		assert.False(t, stringslice.Contains(give, "three"))
	})

	t.Run("Hoist", func(t *testing.T) {
		t.Parallel()

		t.Run("already hoisted", func(t *testing.T) {
			t.Parallel()
			give := []string{"initial", "one", "two"}
			want := []string{"initial", "one", "two"}
			have := stringslice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})

		t.Run("contains the element to hoist", func(t *testing.T) {
			t.Parallel()
			give := []string{"alpha", "initial", "omega"}
			want := []string{"initial", "alpha", "omega"}
			have := stringslice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})

		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := []string{}
			want := []string{}
			have := stringslice.Hoist(give, "initial")
			assert.Equal(t, want, have)
		})

		t.Run("Lines", func(t *testing.T) {
			t.Parallel()
			tests := map[string][]string{
				"":                {""},
				"single line":     {"single line"},
				"multiple\nlines": {"multiple", "lines"},
			}
			for give, want := range tests {
				have := stringslice.Lines(give)
				assert.Equal(t, want, have)
			}
		})
	})

	t.Run("Remove", func(t *testing.T) {
		t.Parallel()
		give := []string{"one", "two", "three"}
		have := stringslice.Remove(give, "two")
		want := []string{"one", "three"}
		assert.Equal(t, have, want)
	})
}
