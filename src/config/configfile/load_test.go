package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/shoenig/test/must"
)

func TestConfigfile(t *testing.T) {
	t.Parallel()

	t.Run("parse", func(t *testing.T) {
		t.Parallel()
		t.Run("complete content", func(t *testing.T) {
			t.Parallel()
			give := `
push-hook = true
push-new-branches = true
ship-delete-tracking-branch = false
sync-before-ship = false
sync-upstream = true

[branches]
main = "main"
perennials = [ "public", "release" ]

[hosting]
platform = "github"
origin-hostname = "github.com"

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
`[1:]
			have, err := configfile.Decode(give)
			must.NoError(t, err)
			github := "github"
			githubCom := "github.com"
			main := "main"
			merge := "merge"
			newBranchPush := true
			pushHook := true
			rebase := "rebase"
			shipDeleteTrackingBranch := false
			syncBeforeShip := false
			syncUpstream := true
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main:       &main,
					Perennials: []string{"public", "release"},
				},
				Hosting: &configfile.Hosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configfile.SyncStrategy{
					FeatureBranches:   &merge,
					PerennialBranches: &rebase,
				},
				PushHook:                 &pushHook,
				PushNewbranches:          &newBranchPush,
				ShipDeleteTrackingBranch: &shipDeleteTrackingBranch,
				SyncBeforeShip:           &syncBeforeShip,
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
					Main:       &main,
					Perennials: nil,
				},
				Hosting:                  nil,
				SyncStrategy:             nil,
				PushNewbranches:          nil,
				PushHook:                 nil,
				ShipDeleteTrackingBranch: nil,
				SyncBeforeShip:           nil,
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
			want := configfile.Data{ //nolint:exhaustruct
				Branches: &configfile.Branches{ //nolint:exhaustruct
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
			want := configfile.Data{ //nolint:exhaustruct
				Branches: &configfile.Branches{ //nolint:exhaustruct
					Perennials: []string{"one", "two"},
				},
			}
			must.Eq(t, want, *have)
		})
	})
}
