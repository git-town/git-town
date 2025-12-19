package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/config/configfile"
	"github.com/git-town/git-town/v22/internal/forge/forgedomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/stringslice"
	"github.com/git-town/git-town/v22/pkg/asserts"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestConfigfile(t *testing.T) {
	t.Parallel()

	t.Run("parse", func(t *testing.T) {
		t.Parallel()
		t.Run("complete content", func(t *testing.T) {
			t.Parallel()
			give := `
create-prototype-branches = true
push-hook = true
push-new-branches = true
ship-delete-tracking-branch = false
ship-strategy = "api"
sync-tags = false
sync-upstream = true

[branches]
main = "main"
contribution-regex = "^gittown-"
default-type = "contribution"
feature-regex = "^kg-"
observed-regex = "^dependabot\\/"
order = "desc"
perennials = [ "public", "staging" ]
perennial-regex = "release-.*"
unknown-type = "prototype"

[create]
branch-prefix = "feature-"
new-branch-type = "prototype"
push-new-branches = true
share-new-branches = "push"

[hosting]
browser = "chrome"
forge-type = "github"
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

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
prototype-branches = "compress"
`[1:]
			haveData, err := configfile.Decode(give)
			must.NoError(t, err)
			wantData := configfile.Data{
				Branches: &configfile.Branches{
					ContributionRegex: Ptr("^gittown-"),
					DefaultType:       Ptr("contribution"),
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
					PushNewbranches:  Ptr(true),
					ShareNewBranches: Ptr("push"),
				},
				Hosting: &configfile.Hosting{
					Browser:        Ptr("chrome"),
					ForgeType:      Ptr("github"),
					OriginHostname: Ptr("github.com"),
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
					Detached:          Ptr(true),
					FeatureStrategy:   Ptr("merge"),
					PerennialStrategy: Ptr("rebase"),
					PrototypeStrategy: Ptr("compress"),
					PushHook:          Ptr(true),
					Tags:              Ptr(false),
					Upstream:          Ptr(true),
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

			finalMessages := stringslice.Collector{}
			haveConfig, err := configfile.Validate(*haveData, finalMessages)
			must.NoError(t, err)
			contributionRegex := asserts.NoError1(configdomain.ParseContributionRegex("^gittown-", "test"))
			featureRegex := asserts.NoError1(configdomain.ParseFeatureRegex("feature-", "test"))
			observedRegex := asserts.NoError1(configdomain.ParseObservedRegex("observed-", "test"))
			perennialRegex := asserts.NoError1(configdomain.ParsePerennialRegex("perennial-", "test"))
			wantConfig := configdomain.PartialConfig{
				Aliases:                  configdomain.Aliases{},
				AutoResolve:              Some(configdomain.AutoResolve(false)),
				AutoSync:                 Some(configdomain.AutoSync(false)),
				BitbucketAppPassword:     None[forgedomain.BitbucketAppPassword](),
				BitbucketUsername:        None[forgedomain.BitbucketUsername](),
				BranchPrefix:             Some(configdomain.BranchPrefix("feature-")),
				BranchTypeOverrides:      configdomain.BranchTypeOverrides{},
				Browser:                  Some(configdomain.Browser("chrome")),
				ContributionRegex:        contributionRegex,
				Detached:                 None[configdomain.Detached](),
				DevRemote:                None[gitdomain.Remote](),
				DisplayTypes:             None[configdomain.DisplayTypes](),
				DryRun:                   None[configdomain.DryRun](),
				FeatureRegex:             featureRegex,
				ForgeType:                None[forgedomain.ForgeType](),
				ForgejoToken:             None[forgedomain.ForgejoToken](),
				GitHubConnectorType:      None[forgedomain.GitHubConnectorType](),
				GitHubToken:              None[forgedomain.GitHubToken](),
				GitLabConnectorType:      None[forgedomain.GitLabConnectorType](),
				GitLabToken:              None[forgedomain.GitLabToken](),
				GitUserEmail:             Option[gitdomain.GitUserEmail]{},
				GitUserName:              None[gitdomain.GitUserName](),
				GiteaToken:               Option[forgedomain.GiteaToken]{},
				HostingOriginHostname:    Option[configdomain.HostingOriginHostname]{},
				IgnoreUncommitted:        Option[configdomain.IgnoreUncommitted]{},
				Lineage:                  configdomain.Lineage{},
				MainBranch:               Option[gitdomain.LocalBranchName]{},
				NewBranchType:            Option[configdomain.NewBranchType]{},
				ObservedRegex:            observedRegex,
				Offline:                  Option[configdomain.Offline]{},
				Order:                    Option[configdomain.Order]{},
				PerennialBranches:        gitdomain.LocalBranchNames{},
				PerennialRegex:           perennialRegex,
				ProposalsShowLineage:     Option[forgedomain.ProposalsShowLineage]{},
				PushBranches:             Option[configdomain.PushBranches]{},
				PushHook:                 Option[configdomain.PushHook]{},
				ShareNewBranches:         Option[configdomain.ShareNewBranches]{},
				ShipDeleteTrackingBranch: Option[configdomain.ShipDeleteTrackingBranch]{},
				ShipStrategy:             Option[configdomain.ShipStrategy]{},
				Stash:                    Option[configdomain.Stash]{},
				SyncFeatureStrategy:      Option[configdomain.SyncFeatureStrategy]{},
				SyncPerennialStrategy:    Option[configdomain.SyncPerennialStrategy]{},
				SyncPrototypeStrategy:    Option[configdomain.SyncPrototypeStrategy]{},
				SyncTags:                 Option[configdomain.SyncTags]{},
				SyncUpstream:             Option[configdomain.SyncUpstream]{},
				UnknownBranchType:        Option[configdomain.UnknownBranchType]{},
				Verbose:                  Option[configdomain.Verbose]{},
			}
			must.Eq(t, wantConfig, haveConfig)
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
	})
}
