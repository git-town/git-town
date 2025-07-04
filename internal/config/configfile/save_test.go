package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v21/internal/config"
	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/config/gitconfig"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("RenderPerennialBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("no perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.LocalBranchNames{}
			have := configfile.RenderPerennialBranches(give)
			want := "[]"
			must.EqOp(t, want, have)
		})
		t.Run("one perennial branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one")
			have := configfile.RenderPerennialBranches(give)
			want := `["one"]`
			must.EqOp(t, want, have)
		})
		t.Run("multiple perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := configfile.RenderPerennialBranches(give)
			want := `["one", "two"]`
			must.EqOp(t, want, have)
		})
	})

	t.Run("RenderTOML", func(t *testing.T) {
		t.Parallel()
		give := config.UnvalidatedConfig{
			UnvalidatedConfig: configdomain.UnvalidatedConfigData{
				MainBranch: gitdomain.NewLocalBranchNameOption("main"),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					DevRemote:                "fork",
					FeatureRegex:             None[configdomain.FeatureRegex](),
					ForgeType:                None[forgedomain.ForgeType](),
					HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
					Lineage:                  configdomain.NewLineage(),
					NewBranchType:            Some(configdomain.BranchTypePrototypeBranch),
					Offline:                  false,
					PerennialBranches:        gitdomain.NewLocalBranchNames("one", "two"),
					PerennialRegex:           None[configdomain.PerennialRegex](),
					PushHook:                 true,
					ShareNewBranches:         configdomain.ShareNewBranchesPush,
					ShipStrategy:             configdomain.ShipStrategySquashMerge,
					ShipDeleteTrackingBranch: true,
					SyncFeatureStrategy:      configdomain.SyncFeatureStrategyMerge,
					SyncPerennialStrategy:    configdomain.SyncPerennialStrategyRebase,
					SyncTags:                 true,
					SyncUpstream:             true,
					UnknownBranchType:        configdomain.BranchTypeFeatureBranch,
				},
			},
		}
		have := configfile.RenderTOML(&give)
		want := `
# More info around this file at https://www.git-town.com/configuration-file

[branches]
main = "main"
perennials = ["one", "two"]
perennial-regex = ""

[create]
new-branch-type = "prototype"
share-new-branches = "push"

[hosting]
dev-remote = "fork"

[ship]
delete-tracking-branch = true
strategy = "squash-merge"

[sync]
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = ""
push-hook = true
tags = true
upstream = true
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		var gitIO gitconfig.IO
		config := config.DefaultUnvalidatedConfig(gitIO, git.EmptyVersion())
		config.UnvalidatedConfig.MainBranch = gitdomain.NewLocalBranchNameOption("main")
		err := configfile.Save(&config)
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := `
# More info around this file at https://www.git-town.com/configuration-file

[branches]
main = "main"
perennials = []
perennial-regex = ""

[create]
new-branch-type = ""
share-new-branches = "no"

[hosting]
dev-remote = "origin"

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "rebase"
push-hook = true
tags = true
upstream = true
`[1:]
		must.EqOp(t, want, have)
	})
}
