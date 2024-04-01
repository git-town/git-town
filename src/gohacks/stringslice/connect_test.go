package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestConnect(t *testing.T) {
	t.Parallel()
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
}
