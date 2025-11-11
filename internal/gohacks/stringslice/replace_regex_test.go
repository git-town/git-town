package stringslice_test

import (
	"regexp"
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestReplaceRegex(t *testing.T) {
	t.Parallel()

	t.Run("has matches", func(t *testing.T) {
		give := []string{"one", "two 123", "three"}
		want := []string{"one", "two SHA", "three"}
		have := stringslice.ReplaceRegex(give, regexp.MustCompile(`\d+`), "SHA")
		must.Eq(t, want, have)
	})

	t.Run("no matches", func(t *testing.T) {
		give := []string{"one", "two"}
		want := []string{"one", "two"}
		have := stringslice.ReplaceRegex(give, regexp.MustCompile(`\d+`), "SHA")
		must.Eq(t, want, have)
	})
}
