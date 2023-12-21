package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/shoenig/test/must"
)

func TestConfigfile(t *testing.T) {
	t.Parallel()

	t.Run("parse", func(t *testing.T) {
		t.Parallel()
		t.Run("complete content", func(t *testing.T) {
			t.Parallel()
			give := `
push-new-branches = true
ship-delete-remote-branch = false
sync-upstream = true

[branches]
main = "main"
perennials = [ "public", "release" ]

[code-hosting]
platform = "github"
origin-hostname = "github.com"

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
`[1:]
			have, err := configdomain.ParseTOML(give)
			must.NoError(t, err)
			github := "github"
			githubCom := "github.com"
			main := "main"
			merge := "merge"
			newBranchPush := true
			rebase := "rebase"
			shipDeleteTrackingBranch := false
			syncUpstream := true
			want := configdomain.ConfigFileData{
				Branches: configdomain.Branches{
					Main:       &main,
					Perennials: []string{"public", "release"},
				},
				CodeHosting: &configdomain.CodeHosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configdomain.SyncStrategy{
					FeatureBranches:   &merge,
					PerennialBranches: &rebase,
				},
				PushNewbranches:          &newBranchPush,
				ShipDeleteTrackingBranch: &shipDeleteTrackingBranch,
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
			have, err := configdomain.ParseTOML(give)
			must.NoError(t, err)
			main := "main"
			want := configdomain.ConfigFileData{
				Branches: configdomain.Branches{
					Main:       &main,
					Perennials: nil,
				},
				CodeHosting:              nil,
				SyncStrategy:             nil,
				PushNewbranches:          nil,
				ShipDeleteTrackingBranch: nil,
				SyncUpstream:             nil,
			}
			must.Eq(t, want, *have)
		})
	})
}
