package format_test

import (
	"testing"

	"github.com/git-town/git-town/v21/internal/cli/format"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("empty", func(t *testing.T) {
		t.Parallel()
		have := format.BranchNames(gitdomain.LocalBranchNames{})
		want := "(none)"
		must.EqOp(t, want, have)
	})

	t.Run("normal", func(t *testing.T) {
		t.Parallel()
		have := format.BranchNames(gitdomain.LocalBranchNames{"one", "two"})
		want := "one, two"
		must.EqOp(t, want, have)
	})
}
