package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configfile"
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
default-type = "prototype"
feature-regex = "^kg-"
observed-regex = "^dependabot\\/"
perennials = [ "public", "staging" ]
perennial-regex = "release-.*"

[create]
new-branch-type = "prototype"
push-new-branches = true

[hosting]
platform = "github"
origin-hostname = "github.com"

[ship]
delete-tracking-branch = false
strategy = "api"

[sync]
feature-strategy = "merge"
perennial-strategy = "rebase"
push-hook = true
tags = false
upstream = true

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
prototype-branches = "compress"
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			contributionRegex := "^gittown-"
			createPrototypeBranches := true
			createDefaultType := "prototype"
			featureRegex := "^kg-"
			github := "github"
			githubCom := "github.com"
			main := "main"
			merge := "merge"
			newBranchType := "prototype"
			observedRegex := `^dependabot\/`
			pushNewBranches := true
			pushHook := true
			rebase := "rebase"
			compress := "compress"
			releaseRegex := "release-.*"
			shipDeleteTrackingBranch := false
			shipStrategy := "api"
			syncTags := false
			syncUpstream := true
			want := configfile.Data{
				Branches: &configfile.Branches{
					ContributionRegex: &contributionRegex,
					DefaultType:       &createDefaultType,
					FeatureRegex:      &featureRegex,
					Main:              &main,
					ObservedRegex:     &observedRegex,
					PerennialRegex:    &releaseRegex,
					Perennials:        []string{"public", "staging"},
				},
				Create: &configfile.Create{
					NewBranchType:   &newBranchType,
					PushNewbranches: &pushNewBranches,
				},
				Hosting: &configfile.Hosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				Ship: &configfile.Ship{
					DeleteTrackingBranch: &shipDeleteTrackingBranch,
					Strategy:             &shipStrategy,
				},
				Sync: &configfile.Sync{
					FeatureStrategy: &merge,
					PushHook:        &pushHook,
					Tags:            &syncTags,
					Upstream:        &syncUpstream,
				},
				SyncStrategy: &configfile.SyncStrategy{
					FeatureBranches:   &merge,
					PerennialBranches: &rebase,
					PrototypeBranches: &compress,
				},
				CreatePrototypeBranches:  &createPrototypeBranches,
				PushHook:                 &pushHook,
				PushNewbranches:          &pushNewBranches,
				ShipDeleteTrackingBranch: &shipDeleteTrackingBranch,
				ShipStrategy:             &shipStrategy,
				SyncTags:                 &syncTags,
				SyncUpstream:             &syncUpstream,
			}
			must.Eq(t, want, *have)
		})

		t.Run("incomplete content", func(t *testing.T) {
			t.Parallel()
			give := `
[branches]
main = "main"
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			main := "main"
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main:           &main,
					Perennials:     nil,
					PerennialRegex: nil,
				},
				Hosting:                  nil,
				SyncStrategy:             nil,
				PushNewbranches:          nil,
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
			main := "main"
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main: &main,
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
