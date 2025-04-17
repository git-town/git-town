package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		have := gitdomain.RemoteBranchNames{
			"origin/branch-3",
			"origin/branch-2",
			"origin/branch-1",
		}
		have.Sort()
		want := gitdomain.RemoteBranchNames{
			"origin/branch-1",
			"origin/branch-2",
			"origin/branch-3",
		}
		must.Eq(t, want, have)
	})

	t.Run("Strings", func(t *testing.T) {
		t.Parallel()
		give := gitdomain.RemoteBranchNames{
			"origin/branch-1",
			"origin/branch-2",
			"origin/branch-3",
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
