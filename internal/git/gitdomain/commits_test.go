package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestCommits(t *testing.T) {
	t.Parallel()

	t.Run("FindByCommitMessage", func(t *testing.T) {
		t.Parallel()
		t.Run("contains a commit with the message", func(t *testing.T) {
			t.Parallel()
			commits := gitdomain.Commits{
				{
					Message: "commit 1",
					SHA:     "111111",
				},
				{
					Message: "commit 2",
					SHA:     "222222",
				},
			}
			have, has := commits.FindByCommitMessage("commit 2").Get()
			must.True(t, has)
			must.EqOp(t, have.Message, "commit 2")
			must.EqOp(t, have.SHA, "222222")
		})
		t.Run("does not contain the commit", func(t *testing.T) {
			t.Parallel()
			commits := gitdomain.Commits{}
			have := commits.FindByCommitMessage("zonk")
			must.True(t, have.IsNone())
		})
	})
}
