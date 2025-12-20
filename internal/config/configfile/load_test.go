package configfile_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/kr/pretty"
	"github.com/shoenig/test/must"
)

func TestConfigfile(t *testing.T) {
	t.Parallel()

	t.Run("parse", func(t *testing.T) {
		t.Parallel()
		t.Run("complete content", func(t *testing.T) {
			t.Parallel()
			giveTOML := `
# See https://www.git-town.com/configuration-file for details

[branches]
contribution-regex = "^gittown-"
display-types = "no main perennial"
feature-regex = "^kg-"
main = "main"
observed-regex = "^dependabot\\/"
order = "desc"
perennials = ["public", "staging"]
perennial-regex = "release-.*"
unknown-type = "prototype"

[create]
branch-prefix = "feature-"
new-branch-type = "prototype"
share-new-branches = "push"
stash = true

[hosting]
browser = "chrome"
dev-remote = "origin"
forge-type = "github"
github-connector = "gh"
gitlab-connector = "glab"
origin-hostname = "github.com"

[propose]
lineage = "cli"

[ship]
delete-tracking-branch = false
ignore-uncommitted = true
strategy = "api"

[sync]
auto-resolve = false
detached = true
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "compress"
push-hook = true
tags = false
upstream = true
`[1:]

			// step 1: decode into low-level data
			haveData, err := configfile.Decode(giveTOML)
			must.NoError(t, err)
			wantData := configfile.Data{
				Branches: &configfile.Branches{
					ContributionRegex: Ptr("^gittown-"),
					DefaultType:       nil,
					DisplayTypes:      Ptr("no main perennial"),
					FeatureRegex:      Ptr("^kg-"),
					Main:              Ptr("main"),
					ObservedRegex:     Ptr(`^dependabot\/`),
					Order:             Ptr("desc"),
					PerennialRegex:    Ptr("release-.*"),
					Perennials:        []string{"public", "staging"},
					UnknownType:       Ptr("prototype"),
				},
				Create: &configfile.Create{
					BranchPrefix:     Ptr("feature-"),
					NewBranchType:    Ptr("prototype"),
					PushNewbranches:  nil,
					ShareNewBranches: Ptr("push"),
					Stash:            Ptr(true),
				},
				Hosting: &configfile.Hosting{
					Browser:             Ptr("chrome"),
					DevRemote:           Ptr("origin"),
					ForgeType:           Ptr("github"),
					GitHubConnectorType: Ptr("gh"),
					GitLabConnectorType: Ptr("glab"),
					OriginHostname:      Ptr("github.com"),
					Platform:            nil,
				},
				Propose: &configfile.Propose{
					Lineage: Ptr("cli"),
				},
				Ship: &configfile.Ship{
					DeleteTrackingBranch: Ptr(false),
					IgnoreUncommitted:    Ptr(true),
					Strategy:             Ptr("api"),
				},
				Sync: &configfile.Sync{
					AutoResolve:       Ptr(false),
					AutoSync:          nil,
					Detached:          Ptr(true),
					FeatureStrategy:   Ptr("merge"),
					PerennialStrategy: Ptr("rebase"),
					PrototypeStrategy: Ptr("compress"),
					PushBranches:      nil,
					PushHook:          Ptr(true),
					Tags:              Ptr(false),
					Upstream:          Ptr(true),
				},
				SyncStrategy:             nil,
				CreatePrototypeBranches:  nil,
				PushHook:                 nil,
				PushNewbranches:          nil,
				ShipDeleteTrackingBranch: nil,
				ShipStrategy:             nil,
				SyncTags:                 nil,
				SyncUpstream:             nil,
			}
			must.Eq(t, wantData, *haveData)

			// step 2: validate into high-level data
			finalMessages := stringslice.NewCollector()
			haveConfig, err := configfile.Validate(*haveData, finalMessages)
			must.NoError(t, err)
			wantConfig := configdomain.PartialConfig{
				Aliases:              configdomain.Aliases{},
				AutoResolve:          Some(configdomain.AutoResolve(false)),
				AutoSync:             None[configdomain.AutoSync](),
				BitbucketAppPassword: None[forgedomain.BitbucketAppPassword](),
				BitbucketUsername:    None[forgedomain.BitbucketUsername](),
				BranchPrefix:         Some(configdomain.BranchPrefix("feature-")),
				BranchTypeOverrides:  configdomain.BranchTypeOverrides{},
				Browser:              Some(configdomain.Browser("chrome")),
				ContributionRegex:    asserts.NoError1(configdomain.ParseContributionRegex("^gittown-", "test")),
				Detached:             Some(configdomain.Detached(true)),
				DevRemote:            Some(gitdomain.Remote("origin")),
				DisplayTypes: Some(configdomain.DisplayTypes{
					BranchTypes: []configdomain.BranchType{configdomain.BranchTypeMainBranch, configdomain.BranchTypePerennialBranch},
					Quantifier:  configdomain.QuantifierNo,
				}),
				DryRun:                   None[configdomain.DryRun](),
				FeatureRegex:             asserts.NoError1(configdomain.ParseFeatureRegex("^kg-", "test")),
				ForgeType:                asserts.NoError1(forgedomain.ParseForgeType("github", "test")),
				ForgejoToken:             None[forgedomain.ForgejoToken](),
				GitHubConnectorType:      Some(forgedomain.GitHubConnectorTypeGh),
				GitHubToken:              None[forgedomain.GitHubToken](),
				GitLabConnectorType:      Some(forgedomain.GitLabConnectorTypeGlab),
				GitLabToken:              None[forgedomain.GitLabToken](),
				GitUserEmail:             None[gitdomain.GitUserEmail](),
				GitUserName:              None[gitdomain.GitUserName](),
				GiteaToken:               None[forgedomain.GiteaToken](),
				HostingOriginHostname:    configdomain.ParseHostingOriginHostname("github.com"),
				IgnoreUncommitted:        Some(configdomain.IgnoreUncommitted(true)),
				Lineage:                  configdomain.NewLineage(),
				MainBranch:               Some(gitdomain.NewLocalBranchName("main")),
				NewBranchType:            Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
				ObservedRegex:            asserts.NoError1(configdomain.ParseObservedRegex("^dependabot\\/", "test")),
				Offline:                  None[configdomain.Offline](),
				Order:                    Some(configdomain.OrderDesc),
				PerennialBranches:        gitdomain.NewLocalBranchNames("public", "staging"),
				PerennialRegex:           asserts.NoError1(configdomain.ParsePerennialRegex("release-.*", "test")),
				ProposalsShowLineage:     Some(forgedomain.ProposalsShowLineageCLI),
				PushBranches:             None[configdomain.PushBranches](),
				PushHook:                 Some(configdomain.PushHook(true)),
				ShareNewBranches:         Some(configdomain.ShareNewBranchesPush),
				ShipDeleteTrackingBranch: Some(configdomain.ShipDeleteTrackingBranch(false)),
				ShipStrategy:             Some(configdomain.ShipStrategyAPI),
				Stash:                    Some(configdomain.Stash(true)),
				SyncFeatureStrategy:      Some(configdomain.SyncFeatureStrategyMerge),
				SyncPerennialStrategy:    Some(configdomain.SyncPerennialStrategyRebase),
				SyncPrototypeStrategy:    Some(configdomain.SyncPrototypeStrategyCompress),
				SyncTags:                 Some(configdomain.SyncTags(false)),
				SyncUpstream:             Some(configdomain.SyncUpstream(true)),
				UnknownBranchType:        Some(configdomain.UnknownBranchType(configdomain.BranchTypePrototypeBranch)),
				Verbose:                  None[configdomain.Verbose](),
			}
			pretty.Ldiff(t, haveConfig, wantConfig)
			must.Eq(t, wantConfig, haveConfig)

			// step 3: serialize back into TOML
			haveTOML := configfile.RenderTOML(haveConfig)
			must.EqOp(t, giveTOML, haveTOML)
		})

		t.Run("incomplete content", func(t *testing.T) {
			t.Parallel()
			give := `
[branches]
main = "main"
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main:           Ptr("main"),
					Perennials:     nil,
					PerennialRegex: nil,
				},
				Create:                   nil,
				Hosting:                  nil,
				SyncStrategy:             nil,
				PushHook:                 nil,
				ShipDeleteTrackingBranch: nil,
				SyncUpstream:             nil,
			}
			must.Eq(t, want, *have)
		})

		t.Run("dotted keys", func(t *testing.T) {
			t.Parallel()
			give := `
branches.main = "main"
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main: Ptr("main"),
				},
			}
			must.Eq(t, want, *have)
		})

		t.Run("multi-line array", func(t *testing.T) {
			t.Parallel()
			give := `
[branches]
perennials = [
	"one",
	"two",
]
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			want := configfile.Data{
				Branches: &configfile.Branches{
					Perennials: []string{"one", "two"},
				},
			}
			must.Eq(t, want, *have)
		})

		t.Run("outdated entries", func(t *testing.T) {
			t.Parallel()
			giveTOML := `
create-prototype-branches = true
push-hook = true
push-new-branches = true
ship-delete-tracking-branch = false
ship-strategy = "api"
sync-tags = false
sync-upstream = true

[branches]
main = "main"
default-type = "contribution"

[create]
push-new-branches = true

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
prototype-branches = "compress"
`[1:]

			// step 1: decode into low-level data
			haveData, err := configfile.Decode(giveTOML)
			must.NoError(t, err)
			wantData := configfile.Data{
				Branches: &configfile.Branches{
					DefaultType: Ptr("contribution"),
					Main:        Ptr("main"),
				},
				Create: &configfile.Create{
					PushNewbranches: Ptr(true),
				},
				SyncStrategy: &configfile.SyncStrategy{
					FeatureBranches:   Ptr("merge"),
					PerennialBranches: Ptr("rebase"),
					PrototypeBranches: Ptr("compress"),
				},
				CreatePrototypeBranches:  Ptr(true),
				PushHook:                 Ptr(true),
				PushNewbranches:          Ptr(true),
				ShipDeleteTrackingBranch: Ptr(false),
				ShipStrategy:             Ptr("api"),
				SyncTags:                 Ptr(false),
				SyncUpstream:             Ptr(true),
			}
			must.Eq(t, wantData, *haveData)

			// step 2: validate into high-level data
			finalMessages := stringslice.NewCollector()
			haveConfig, err := configfile.Validate(*haveData, finalMessages)
			must.NoError(t, err)
			wantConfig := configdomain.PartialConfig{
				Aliases:                  configdomain.Aliases{},
				BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
				Lineage:                  configdomain.NewLineage(),
				MainBranch:               Some(gitdomain.NewLocalBranchName("main")),
				NewBranchType:            Some(configdomain.NewBranchType(configdomain.BranchTypePrototypeBranch)),
				PerennialBranches:        gitdomain.LocalBranchNames{},
				PushHook:                 Some(configdomain.PushHook(true)),
				ShareNewBranches:         Some(configdomain.ShareNewBranchesPush),
				ShipDeleteTrackingBranch: Some(configdomain.ShipDeleteTrackingBranch(false)),
				ShipStrategy:             Some(configdomain.ShipStrategyAPI),
				SyncFeatureStrategy:      Some(configdomain.SyncFeatureStrategyMerge),
				SyncPerennialStrategy:    Some(configdomain.SyncPerennialStrategyRebase),
				SyncPrototypeStrategy:    Some(configdomain.SyncPrototypeStrategyCompress),
				SyncTags:                 Some(configdomain.SyncTags(false)),
				SyncUpstream:             Some(configdomain.SyncUpstream(true)),
				UnknownBranchType:        Some(configdomain.UnknownBranchType(configdomain.BranchTypeContributionBranch)),
			}
			haveJSON := asserts.NoError1(json.MarshalIndent(haveConfig, "", "  "))
			wantJSON := asserts.NoError1(json.MarshalIndent(wantConfig, "", "  "))
			must.EqOp(t, string(wantJSON), string(haveJSON))
		})
	})
}
