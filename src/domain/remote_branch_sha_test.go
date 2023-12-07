package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchesSHAs(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchesSHAs{
			domain.NewRemoteBranchName("origin/feature-branch"):   domain.NewSHA("111111"),
			domain.NewRemoteBranchName("origin/perennial-branch"): domain.NewSHA("222222"),
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("perennial-branch"),
		}
		havePerennials, haveFeatures := give.Categorize(branchTypes)
		wantPerennials := domain.RemoteBranchesSHAs{
			domain.NewRemoteBranchName("origin/perennial-branch"): domain.NewSHA("222222"),
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchesSHAs{
			domain.NewRemoteBranchName("origin/feature-branch"): domain.NewSHA("111111"),
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
