package config_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	"github.com/git-town/git-town/v16/test/testruntime"
	"github.com/shoenig/test/must"
)

func TestUnvalidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("Reload", func(t *testing.T) {
		t.Parallel()
		t.Run("lineage changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.CreateFeatureBranch("branch", "main")
			repo.Config.Reload()
			want := configdomain.NewLineageBuilder()
			want.Add("branch", "main")
			must.Eq(t, want.Lineage(), repo.Config.NormalConfig.Lineage)
		})
		t.Run("contribution branches changed", func(t *testing.T) {
			t.Parallel()
			repo := testruntime.CreateGitTown(t)
			repo.CreateBranch("branch", "main")
			err := repo.Config.NormalConfig.AddToContributionBranches("branch")
			must.NoError(t, err)
			repo.Config.Reload()
			want := gitdomain.NewLocalBranchNames("branch")
			must.Eq(t, want, repo.Config.NormalConfig.ContributionBranches)
		})
	})
}
