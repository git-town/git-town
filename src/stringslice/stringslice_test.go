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

}
