package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v16/internal/config/configfile"
	. "github.com/git-town/git-town/v16/pkg/prelude"
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
prototype-strategy = "compress"
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
			want := configfile.Data{
				Branches: &configfile.Branches{
					ContributionRegex: Ptr("^gittown-"),
					DefaultType:       Ptr("prototype"),
					FeatureRegex:      Ptr("^kg-"),
					Main:              Ptr("main"),
					ObservedRegex:     Ptr(`^dependabot\/`),
					PerennialRegex:    Ptr("release-.*"),
					Perennials:        []string{"public", "staging"},
				},
				Create: &configfile.Create{
					NewBranchType:   Ptr("prototype"),
					PushNewbranches: Ptr(true),
				},
				Hosting: &configfile.Hosting{
					Platform:       Ptr("github"),
					OriginHostname: Ptr("github.com"),
				},
				Ship: &configfile.Ship{
					DeleteTrackingBranch: Ptr(false),
					Strategy:             Ptr("api"),
				},
				Sync: &configfile.Sync{
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
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main:           Ptr("main"),
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
