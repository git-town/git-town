package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
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
				AutoResolve:  Some(configdomain.AutoResolve(false)),
				BranchPrefix: Some(configdomain.BranchPrefix("feature-")),
				DevRemote:    Some(gitdomain.RemoteOrigin),
				DisplayTypes: Some(configdomain.DisplayTypes{
					BranchTypes: []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch},
					Quantifier:  configdomain.QuantifierNo,
				}),
				ForgeType:                asserts.NoError1(forgedomain.ParseForgeType("github", "test")),
				HostingOriginHostname:    configdomain.ParseHostingOriginHostname("forge"),
				MainBranch:               Some(gitdomain.NewLocalBranchName("main")),
				NewBranchType:            Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
				Order:                    Some(configdomain.OrderDesc),
				PerennialBranches:        gitdomain.NewLocalBranchNames("qa", "staging"),
				PerennialRegex:           asserts.NoError1(configdomain.ParsePerennialRegex("perennial-", "test")),
				ProposalsShowLineage:     Some(forgedomain.ProposalsShowLineageCLI),
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
# See https://www.git-town.com/configuration-file for details

[branches]
main = "main"
order = "desc"
perennials = ["qa", "staging"]
perennial-regex = "perennial-"
display-types = "no main perennial"

[create]
branch-prefix = "feature-"
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
# See https://www.git-town.com/configuration-file for details
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
