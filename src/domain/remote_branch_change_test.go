package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchChange(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		give := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := give.Categorize(branchTypes)
		wantPerennials := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.RemoteBranchChange{
			gitdomain.NewRemoteBranchName("origin/branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
