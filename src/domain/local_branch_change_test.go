package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestLocalBranchChange(t *testing.T) {
	t.Parallel()

	t.Run("Categorize", func(t *testing.T) {
		t.Parallel()
		give := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		branchTypes := domain.BranchTypes{
			MainBranch:        gitdomain.NewLocalBranchName("main"),
			PerennialBranches: gitdomain.NewLocalBranchNames("dev"),
		}
		havePerennials, haveFeatures := give.Categorize(branchTypes)
		wantPerennials := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("dev"): {
				Before: gitdomain.NewSHA("333333"),
				After:  gitdomain.NewSHA("444444"),
			},
		}
		must.Eq(t, wantPerennials, havePerennials)
		wantFeatures := domain.LocalBranchChange{
			gitdomain.NewLocalBranchName("branch-1"): {
				Before: gitdomain.NewSHA("111111"),
				After:  gitdomain.NewSHA("222222"),
			},
		}
		must.Eq(t, wantFeatures, haveFeatures)
	})
}
