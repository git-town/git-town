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
		t.Run("valid content", func(t *testing.T) {
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
	})

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()
		t.Run("fully configured", func(t *testing.T) {
			t.Parallel()
			github := "github"
			githubCom := "github.com"
			mainBranch := "main"
			merge := "merge"
			newBranchPush := false
			rebase := "rebase"
			shipDeleteTrackingBranch := false
			syncUpstream := true
			give := configdomain.ConfigFileData{
				Branches: configdomain.Branches{
					Main:       &mainBranch,
					Perennials: []string{"public", "qa"},
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
			have := configdomain.EncodeConfigFile(give)
			want := `
push-new-branches = false
ship-delete-remote-branch = false
sync-upstream = true

[branches]
  main = "main"
  perennials = ["public", "qa"]

[code-hosting]
  platform = "github"
  origin-hostname = "github.com"

[sync-strategy]
  feature-branches = "merge"
  perennial-branches = "rebase"
`[1:]
			must.EqOp(t, want, have)
		})

		t.Run("partially configured", func(t *testing.T) {
			t.Parallel()
			mainBranch := "main"
			give := configdomain.ConfigFileData{
				Branches: configdomain.Branches{
					Main:       &mainBranch,
					Perennials: []string{"public", "qa"},
				},
			}
			have := configdomain.EncodeConfigFile(give)
			want := `
[branches]
  main = "main"
  perennials = ["public", "qa"]
`[1:]
			must.EqOp(t, want, have)
		})
	})
}
