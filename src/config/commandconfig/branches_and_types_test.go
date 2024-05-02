package commandconfig_test

import (
	"testing"

	"github.com/git-town/git-town/v14/src/config/commandconfig"
	"github.com/git-town/git-town/v14/src/config/configdomain"
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	. "github.com/git-town/git-town/v14/src/gohacks/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchesAndTypes(t *testing.T) {
	t.Parallel()

	t.Run("Add", func(t *testing.T) {
		t.Parallel()
		have := commandconfig.BranchesAndTypes{}
		unvalidatedConfig := configdomain.UnvalidatedConfig{ //nolint:exhaustruct
			MainBranch: Some(gitdomain.NewLocalBranchName("main")),
		}
		have.Add("main", unvalidatedConfig)
		want := map[gitdomain.LocalBranchName]configdomain.BranchType{
			"main": configdomain.BranchTypeMainBranch,
		}
		must.Eq(t, want, have)
	})

	t.Run("AddMany", func(t *testing.T) {
		t.Parallel()
		have := commandconfig.BranchesAndTypes{}
		unvalidatedConfig := configdomain.UnvalidatedConfig{ //nolint:exhaustruct
			MainBranch:        Some(gitdomain.NewLocalBranchName("main")),
			PerennialBranches: gitdomain.NewLocalBranchNames("perennial"),
		}
		have.AddMany(gitdomain.NewLocalBranchNames("main", "perennial"), unvalidatedConfig)
		want := map[gitdomain.LocalBranchName]configdomain.BranchType{
			"main":      configdomain.BranchTypeMainBranch,
			"perennial": configdomain.BranchTypePerennialBranch,
		}
		must.Eq(t, want, have)
	})

	t.Run("Keys", func(t *testing.T) {
		t.Parallel()
		give := commandconfig.BranchesAndTypes{
			"main":      configdomain.BranchTypeMainBranch,
			"perennial": configdomain.BranchTypePerennialBranch,
		}
		want := gitdomain.NewLocalBranchNames("main", "perennial")
		have := give.Keys()
		must.Eq(t, want, have)
	})
}
