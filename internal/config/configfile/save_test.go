package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/pkg/asserts"
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
		t.Run("all options given", func(t *testing.T) {
			t.Parallel()
			have := configfile.RenderTOML(configdomain.PartialConfig{
				DevRemote:                Some(gitdomain.RemoteOrigin),
				ForgeType:                asserts.NoError1(forgedomain.ParseForgeType("github")),
				HostingOriginHostname:    configdomain.ParseHostingOriginHostname("forge"),
				MainBranch:               Some(gitdomain.NewLocalBranchName("main")),
				NewBranchType:            Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
				NoAutoResolve:            Some(configdomain.NoAutoResolve(false)),
				PerennialBranches:        gitdomain.NewLocalBranchNames("qa", "staging"),
				PerennialRegex:           asserts.NoError1(configdomain.ParsePerennialRegex("perennial-")),
				ProposalsShowLineage:     Some(configdomain.ProposalsShowLineageCLI),
				PushHook:                 Some(configdomain.PushHook(true)),
				ShareNewBranches:         Some(configdomain.ShareNewBranchesPropose),
				ShipDeleteTrackingBranch: Some(configdomain.ShipDeleteTrackingBranch(true)),
				ShipStrategy:             Some(configdomain.ShipStrategyAPI),
				SyncFeatureStrategy:      Some(configdomain.SyncFeatureStrategyMerge),
				SyncPerennialStrategy:    Some(configdomain.SyncPerennialStrategyRebase),
				SyncPrototypeStrategy:    Some(configdomain.SyncPrototypeStrategyCompress),
				SyncTags:                 Some(configdomain.SyncTags(true)),
				SyncUpstream:             Some(configdomain.SyncUpstream(true)),
			})
			want := `
# More info around this file at https://www.git-town.com/configuration-file

[branches]
main = "main"
perennials = ["qa", "staging"]
perennial-regex = "perennial-"

[create]
new-branch-type = "prototype"
share-new-branches = "propose"

[hosting]
dev-remote = "origin"
forge-type = "github"
origin-hostname = "forge"

[propose]
lineage = "cli"

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
auto-resolve = false
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "compress"
push-hook = true
tags = true
upstream = true
`[1:]
			must.EqOp(t, want, have)
		})
		t.Run("no options given", func(t *testing.T) {
			t.Parallel()
			have := configfile.RenderTOML(configdomain.PartialConfig{})
			want := `
# More info around this file at https://www.git-town.com/configuration-file
`[1:]
			must.EqOp(t, want, have)
		})
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		config := configdomain.PartialConfig{}
		err := configfile.Save(config)
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := configfile.RenderTOML(config)
		must.EqOp(t, want, have)
	})
}
