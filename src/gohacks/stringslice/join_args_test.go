package stringslice_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/gohacks/stringslice"
	"github.com/shoenig/test/must"
)

func TestJoinArgs(t *testing.T) {
	t.Parallel()
	args := []string{"one", "", "two", "the args"}
	have := stringslice.JoinArgs(args)
	want := `one "" two "the args"`
	must.EqOp(t, want, have)
}
