package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestFirstLine(t *testing.T) {
	t.Parallel()
	t.Run("multi-line string", func(t *testing.T) {
		t.Parallel()
		give := "one\ntwo\nthree"
		have := stringslice.FirstLine(give)
		want := "one"
		must.EqOp(t, want, have)
	})
}
