package gitdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchNames(t *testing.T) {
	t.Parallel()

	t.Run("LocalBranchNames", func(t *testing.T) {
		t.Run("empty", func(t *testing.T) {
			t.Parallel()
			branchNames := gitdomain.BranchNames{}
			have := branchNames.LocalBranchNames()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})

		t.Run("local branches", func(t *testing.T) {
			t.Parallel()
			branchNames := gitdomain.BranchNames{
				gitdomain.NewBranchName("main"),
				gitdomain.NewBranchName("feature"),
				gitdomain.NewBranchName("develop"),
			}
			have := branchNames.LocalBranchNames()
			want := gitdomain.NewLocalBranchNames("main", "feature", "develop")
			must.Eq(t, want, have)
		})

		t.Run("remote branches", func(t *testing.T) {
			t.Parallel()
			branchNames := gitdomain.BranchNames{
				gitdomain.NewBranchName("origin/main"),
				gitdomain.NewBranchName("origin/feature"),
				gitdomain.NewBranchName("origin/develop"),
			}
			have := branchNames.LocalBranchNames()
			want := gitdomain.LocalBranchNames{}
			must.Eq(t, want, have)
		})

		t.Run("mixed local and remote branches", func(t *testing.T) {
			t.Parallel()
			branchNames := gitdomain.BranchNames{
				gitdomain.NewBranchName("main"),
				gitdomain.NewBranchName("origin/feature"),
				gitdomain.NewBranchName("develop"),
				gitdomain.NewBranchName("origin/hotfix"),
			}
			have := branchNames.LocalBranchNames()
			want := gitdomain.NewLocalBranchNames("main", "develop")
			must.Eq(t, want, have)
		})
	})
}
