package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchChange(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchChange{
			domain.NewRemoteBranchName("origin/branch-1"): {
				Before: domain.NewSHA("111111"),
				After:  domain.NewSHA("222222"),
			},
			domain.NewRemoteBranchName("origin/dev"): {
				Before: domain.NewSHA("333333"),
				After:  domain.NewSHA("444444"),
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        domain.NewLocalBranchName("main"),
			PerennialBranches: domain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := give.Categorize(branchTypes)
		wantPerennials := domain.RemoteBranchChange{
			domain.NewRemoteBranchName("origin/dev"): {
				Before: domain.NewSHA("333333"),
				After:  domain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchChange{
			domain.NewRemoteBranchName("origin/branch-1"): {
				Before: domain.NewSHA("111111"),
				After:  domain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
