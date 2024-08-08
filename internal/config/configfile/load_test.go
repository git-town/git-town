package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v15/internal/config/configfile"
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
sync-tags = false
sync-upstream = true

[branches]
main = "main"
perennials = [ "public", "staging" ]
perennial-regex = "release-.*"

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
			pushNewBranches := true
			pushHook := true
			rebase := "rebase"
			releaseRegex := "release-.*"
			shipDeleteTrackingBranch := false
			syncTags := false
			syncUpstream := true
			want := configfile.Data{
				Branches: &configfile.Branches{
					Main:           &main,
					Perennials:     []string{"public", "staging"},
					PerennialRegex: &releaseRegex,
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
				PushNewbranches:          &pushNewBranches,
				ShipDeleteTrackingBranch: &shipDeleteTrackingBranch,
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
