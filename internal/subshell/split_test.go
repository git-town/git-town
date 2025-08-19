package subshell_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/subshell"
	"github.com/shoenig/test/must"
)

func TestSplit(t *testing.T) {
	t.Parallel()
	t.Run("empty text", func(t *testing.T) {
		t.Parallel()
		_, _, err := subshell.Split("")
		must.Error(t, err)
		must.EqOp(t, "empty", err.Error())
	})
	t.Run("only executable", func(t *testing.T) {
		t.Parallel()
		executable, args, err := subshell.Split("op")
		must.NoError(t, err)
		must.EqOp(t, "op", executable)
		must.Len(t, 0, args)
	})
	t.Run("executable and arguments", func(t *testing.T) {
		t.Parallel()
		executable, args, err := subshell.Split("op read op://vault/github/token")
		must.NoError(t, err)
		must.EqOp(t, "op", executable)
		must.Eq(t, []string{"read", "op://vault/github/token"}, args)
	})
	t.Run("broken script", func(t *testing.T) {
		t.Parallel()
		_, _, err := subshell.Split("op 'unclosed")
		must.Error(t, err)
		must.EqOp(t, "xxx", err.Error())
	})
}
