package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchesSHAs(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"):   gitdomain.NewSHA("111111"),
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial-branch"),
		}
		havePerennials, haveFeatures := give.Categorize(branchTypes)
		wantPerennials := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/perennial-branch"): gitdomain.NewSHA("222222"),
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchesSHAs{
			gitdomain.NewRemoteBranchName("origin/feature-branch"): gitdomain.NewSHA("111111"),
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
