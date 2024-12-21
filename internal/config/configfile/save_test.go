package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v17/internal/config"
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/config/configfile"
	"github.com/git-town/git-town/v17/internal/config/gitconfig"
	"github.com/git-town/git-town/v17/internal/git"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
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
				MainBranch: Some(gitdomain.NewLocalBranchName("main")),
			},
			NormalConfig: config.NormalConfig{
				NormalConfigData: configdomain.NormalConfigData{
					DefaultBranchType:        configdomain.BranchTypeFeatureBranch,
					DevRemote:                "fork",
					FeatureRegex:             None[configdomain.FeatureRegex](),
					HostingOriginHostname:    None[configdomain.HostingOriginHostname](),
					HostingPlatform:          None[configdomain.HostingPlatform](),
					Lineage:                  configdomain.NewLineage(),
					NewBranchType:            configdomain.BranchTypePrototypeBranch,
					ObservedBranches:         gitdomain.LocalBranchNames{},
					Offline:                  false,
					ParkedBranches:           gitdomain.LocalBranchNames{},
					PerennialBranches:        gitdomain.NewLocalBranchNames("one", "two"),
					PerennialRegex:           None[configdomain.PerennialRegex](),
					PushHook:                 true,
					PushNewBranches:          false,
					ShipStrategy:             configdomain.ShipStrategySquashMerge,
					ShipDeleteTrackingBranch: true,
					SyncFeatureStrategy:      configdomain.SyncFeatureStrategyMerge,
					SyncPerennialStrategy:    configdomain.SyncPerennialStrategyRebase,
					SyncTags:                 true,
					SyncUpstream:             true,
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
push-new-branches = false

[hosting]
dev-remote = "fork"
# platform = ""
# origin-hostname = ""

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
		var gitAccess gitconfig.Access
		config := config.DefaultUnvalidatedConfig(gitAccess, git.EmptyVersion())
		config.UnvalidatedConfig.MainBranch = Some(gitdomain.NewLocalBranchName("main"))
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
new-branch-type = "feature"
push-new-branches = false

[hosting]
dev-remote = "origin"
# platform = ""
# origin-hostname = ""

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
