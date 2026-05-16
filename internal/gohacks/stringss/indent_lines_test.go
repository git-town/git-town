package stringss_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	"github.com/shoenig/test/must"
)

func TestIndentLines(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		have := stringss.IndentLines("", 2)
		must.EqOp(t, "", have)
	})

	t.Run("indent by 2 spaces", func(t *testing.T) {
		t.Parallel()
		give := "one\ntwo\nthree"
		have := stringss.IndentLines(give, 2)
		want := "  one\n  two\n  three"
		must.EqOp(t, want, have)
	})

	t.Run("indent by 4 spaces", func(t *testing.T) {
		t.Parallel()
		give := "one\ntwo\nthree"
		have := stringss.IndentLines(give, 4)
		want := "    one\n    two\n    three"
		must.EqOp(t, want, have)
	})
}
