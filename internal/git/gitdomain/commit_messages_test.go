package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v24/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommitMessages(t *testing.T) {
	t.Parallel()

	t.Run("Join", func(t *testing.T) {
		t.Parallel()
		t.Run("contains elements", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewCommitMessages("commit 1", "commit 2\n\nbody line 1", "commit 3")
			have := give.Join("\n\n")
			want := gitdomain.CommitMessage("commit 1\n\ncommit 2\n\nbody line 1\n\ncommit 3")
			must.EqOp(t, want, have)
		})
		t.Run("empty list", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.CommitMessages{}
			have := give.Join("\n\n")
			want := gitdomain.CommitMessage("")
			must.EqOp(t, want, have)
		})
		t.Run("trims whitespace", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.CommitMessages{"", "commit"}
			have := give.Join("\n\n")
			want := gitdomain.CommitMessage("commit")
			must.EqOp(t, want, have)
		})
	})
}
