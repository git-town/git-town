package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configdomain"
	"github.com/git-town/git-town/v16/internal/git/gitdomain"
	. "github.com/git-town/git-town/v16/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestValidatedConfig(t *testing.T) {
	t.Parallel()

	t.Run("Author", func(t *testing.T) {
		t.Parallel()
		conf := configdomain.ValidatedConfig{
			GitUserName:  configdomain.GitUserName("name"),
			GitUserEmail: configdomain.GitUserEmail("email"),
		}
		have := conf.Author()
		want := gitdomain.Author("name <email>")
		must.EqOp(t, want, have)
	})

	t.Run("IsMainBranch", func(t *testing.T) {
		t.Parallel()
		config := configdomain.ValidatedConfig{
			MainBranch: gitdomain.NewLocalBranchName("main"),
		}
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("feature")))
		must.True(t, config.IsMainBranch(gitdomain.NewLocalBranchName("main")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren1")))
		must.False(t, config.IsMainBranch(gitdomain.NewLocalBranchName("peren2")))
	})

	t.Run("MainAndPerennials", func(t *testing.T) {
		t.Parallel()
		config := configdomain.UnvalidatedConfig{
			MainBranch: Some(gitdomain.NewLocalBranchName("main")),
			NormalConfig: &configdomain.NormalConfig{
				PerennialBranches: gitdomain.NewLocalBranchNames("perennial-1", "perennial-2"),
			},
		}
		have := config.MainAndPerennials()
		want := gitdomain.NewLocalBranchNames("main", "perennial-1", "perennial-2")
		must.Eq(t, want, have)
	})
}
