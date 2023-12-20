package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
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
			have, err := configdomain.ParseConfigFileData(give)
			must.NoError(t, err)
			github := configdomain.CodeHostingPlatformName("github")
			githubCom := configdomain.CodeHostingOriginHostname("github.com")
			merge := "merge"
			newBranchPush := configdomain.NewBranchPush(true)
			shipDeleteTrackingBranch := configdomain.ShipDeleteTrackingBranch(false)
			syncUpstream := configdomain.SyncUpstream(true)
			mainBranch := domain.NewLocalBranchName("main")
			want := configdomain.ConfigFile{
				Branches: configdomain.Branches{
					Main:       &mainBranch,
					Perennials: domain.NewLocalBranchNames("public", "release"),
				},
				CodeHosting: &configdomain.CodeHosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configdomain.SyncStrategy{
					FeatureBranches:   &merge,
					PerennialBranches: &configdomain.SyncPerennialStrategyRebase,
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
			github := configdomain.CodeHostingPlatformName("github")
			githubCom := configdomain.CodeHostingOriginHostname("github.com")
			newBranchPush := configdomain.NewBranchPush(false)
			shipDeleteTrackingBranch := configdomain.ShipDeleteTrackingBranch(false)
			syncFeatureStrategy := "merge"
			syncUpstream := configdomain.SyncUpstream(true)
			mainBranch := domain.NewLocalBranchName("main")
			give := configdomain.ConfigFile{
				Branches: configdomain.Branches{
					Main:       &mainBranch,
					Perennials: domain.NewLocalBranchNames("public", "qa"),
				},
				CodeHosting: &configdomain.CodeHosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configdomain.SyncStrategy{
					FeatureBranches:   &syncFeatureStrategy,
					PerennialBranches: &configdomain.SyncPerennialStrategyRebase,
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
			mainBranch := domain.NewLocalBranchName("main")
			give := configdomain.ConfigFile{
				Branches: configdomain.Branches{
					Main:       &mainBranch,
					Perennials: domain.NewLocalBranchNames("public", "qa"),
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
