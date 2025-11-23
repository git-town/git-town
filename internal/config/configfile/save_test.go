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
			contributionRegex := asserts.NoError1(configdomain.ParseContributionRegex("contribution-", "test"))
			featureRegex := asserts.NoError1(configdomain.ParseFeatureRegex("feature-", "test"))
			observedRegex := asserts.NoError1(configdomain.ParseObservedRegex("observed-", "test"))
			perennialRegex := asserts.NoError1(configdomain.ParsePerennialRegex("perennial-", "test"))
			have := configfile.RenderTOML(configdomain.PartialConfig{
				AutoResolve:       Some(configdomain.AutoResolve(false)),
				BranchPrefix:      Some(configdomain.BranchPrefix("feature-")),
				ContributionRegex: contributionRegex,
				Detached:          Some(configdomain.Detached(true)),
				DevRemote:         Some(gitdomain.RemoteOrigin),
				DisplayTypes: Some(configdomain.DisplayTypes{
					BranchTypes: []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch},
					Quantifier:  configdomain.QuantifierNo,
				}),
				FeatureRegex:             featureRegex,
				ForgeType:                asserts.NoError1(forgedomain.ParseForgeType("github", "test")),
				GitHubConnectorType:      Some(forgedomain.GitHubConnectorTypeGh),
				GitLabConnectorType:      Some(forgedomain.GitLabConnectorTypeGlab),
				HostingOriginHostname:    configdomain.ParseHostingOriginHostname("forge"),
				MainBranch:               Some(gitdomain.NewLocalBranchName("main")),
				NewBranchType:            Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
				ObservedRegex:            observedRegex,
				Order:                    Some(configdomain.OrderDesc),
				PerennialBranches:        gitdomain.NewLocalBranchNames("qa", "staging"),
				PerennialRegex:           perennialRegex,
				ProposalsShowLineage:     Some(forgedomain.ProposalsShowLineageCLI),
				PushBranches:             Some(configdomain.PushBranches(true)),
				PushHook:                 Some(configdomain.PushHook(true)),
				ShareNewBranches:         Some(configdomain.ShareNewBranchesPropose),
				ShipDeleteTrackingBranch: Some(configdomain.ShipDeleteTrackingBranch(true)),
				ShipStrategy:             Some(configdomain.ShipStrategyAPI),
				Stash:                    Some(configdomain.Stash(true)),
				SyncFeatureStrategy:      Some(configdomain.SyncFeatureStrategyMerge),
				SyncPerennialStrategy:    Some(configdomain.SyncPerennialStrategyRebase),
				SyncPrototypeStrategy:    Some(configdomain.SyncPrototypeStrategyCompress),
				SyncTags:                 Some(configdomain.SyncTags(true)),
				SyncUpstream:             Some(configdomain.SyncUpstream(true)),
				UnknownBranchType:        Some(configdomain.UnknownBranchType(configdomain.BranchTypePrototypeBranch)),
			})
			want := `
# See https://www.git-town.com/configuration-file for details

[branches]
contribution-regex = "contribution-"
display-types = "no main perennial"
feature-regex = "feature-"
main = "main"
observed-regex = "observed-"
order = "desc"
perennials = ["qa", "staging"]
perennial-regex = "perennial-"
unknown-branch-type = "prototype"

[create]
branch-prefix = "feature-"
new-branch-type = "prototype"
share-new-branches = "propose"
stash = true

[hosting]
dev-remote = "origin"
forge-type = "github"
github-connector-type = "gh"
gitlab-connector-type = "glab"
origin-hostname = "forge"

[propose]
lineage = "cli"

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
auto-resolve = false
detached = true
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "compress"
push-branches = true
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
