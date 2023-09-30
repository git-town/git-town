package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/stringslice"
	"github.com/stretchr/testify/assert"
)

func TestStringSlice(t *testing.T) {
	t.Parallel()

	t.Run("Connect", func(t *testing.T) {
		t.Run("no element", func(t *testing.T) {
			t.Parallel()
			give := []string{}
			have := stringslice.Connect(give)
			want := ""
			assert.Equal(t, want, have)
		})

		t.Run("single element", func(t *testing.T) {
			t.Parallel()
			give := []string{"one"}
			have := stringslice.Connect(give)
			want := "\"one\""
			assert.Equal(t, want, have)
		})

		t.Run("two elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two"}
			have := stringslice.Connect(give)
			want := "\"one\" and \"two\""
			assert.Equal(t, want, have)
		})

		t.Run("three elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three"}
			have := stringslice.Connect(give)
			want := "\"one\", \"two\", and \"three\""
			assert.Equal(t, want, have)
		})

		t.Run("four elements", func(t *testing.T) {
			t.Parallel()
			give := []string{"one", "two", "three", "four"}
			have := stringslice.Connect(give)
			want := "\"one\", \"two\", \"three\", and \"four\""
			assert.Equal(t, want, have)
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
			assert.Equal(t, want, have)
		}
	})
}
