package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v13/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestSurroundEmptyWith(t *testing.T) {
	t.Parallel()
	give := []string{"git", "config", "perennial-branches", ""}
	have := stringslice.SurroundEmptyWith(give, `"`)
	want := []string{"git", "config", "perennial-branches", `""`}
	must.Eq(t, want, have)
}
