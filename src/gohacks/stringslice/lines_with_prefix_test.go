package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestLinesWithPrefix(t *testing.T) {
	t.Parallel()

	t.Run("multiple matching lines", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"a line",
			"another line",
			"a line with more text",
		}
		have := stringslice.LinesWithPrefix(lines, "a line")
		want := []string{
			"a line",
			"a line with more text",
		}
		must.Eq(t, want, have)
	})

	t.Run("one matching line", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"* (no branch, rebasing feature)",
			"  feature",
			"+ main",
		}
		have := stringslice.LinesWithPrefix(lines, "* ")
		want := []string{
			"* (no branch, rebasing feature)",
		}
		must.Eq(t, want, have)
	})

	t.Run("no matching line", func(t *testing.T) {
		t.Parallel()
		lines := []string{
			"one",
			"two",
		}
		have := stringslice.LinesWithPrefix(lines, "zonk")
		want := []string{}
		must.Eq(t, want, have)
	})
}
