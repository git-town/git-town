package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestSurroundSpacesWith(t *testing.T) {
	t.Parallel()
	give := []string{"git", "reflog", "--format=%H %s"}
	have := stringslice.SurroundSpacesWith(give, `"`)
	want := []string{"git", "reflog", `"--format=%H %s"`}
	must.Eq(t, want, have)
}
