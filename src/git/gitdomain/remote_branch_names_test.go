package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		have := gitdomain.RemoteBranchNames{
			gitdomain.NewRemoteBranchName("origin/branch-3"),
			gitdomain.NewRemoteBranchName("origin/branch-2"),
			gitdomain.NewRemoteBranchName("origin/branch-1"),
		}
		have.Sort()
		want := gitdomain.RemoteBranchNames{
			gitdomain.NewRemoteBranchName("origin/branch-1"),
			gitdomain.NewRemoteBranchName("origin/branch-2"),
			gitdomain.NewRemoteBranchName("origin/branch-3"),
		}
		must.Eq(t, want, have)
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		give := gitdomain.RemoteBranchNames{
			gitdomain.NewRemoteBranchName("origin/branch-1"),
			gitdomain.NewRemoteBranchName("origin/branch-2"),
			gitdomain.NewRemoteBranchName("origin/branch-3"),
		}
		have := give.Strings()
		want := []string{
			"origin/branch-1",
			"origin/branch-2",
			"origin/branch-3",
		}
		must.Eq(t, want, have)
	})
}
