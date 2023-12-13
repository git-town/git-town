package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
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
			have, err := configfile.Parse(give)
			must.NoError(t, err)
			github := "github"
			githubCom := "github.com"
			merge := "merge"
			rebase := "rebase"
			boolFalse := false
			boolTrue := true
			want := configfile.ConfigFile{
				Branches: configfile.Branches{
					Main:       "main",
					Perennials: []string{"public", "release"},
				},
				CodeHosting: &configfile.CodeHosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configfile.SyncStrategy{
					FeatureBranches:   &merge,
					PerennialBranches: &rebase,
				},
				PushNewbranches:        &boolTrue,
				ShipDeleteRemoteBranch: &boolFalse,
				SyncUpstream:           &boolTrue,
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
			boolFalse := false
			boolTrue := true
			give := configfile.ConfigFile{
				Branches: configfile.Branches{
					Main:       "main",
					Perennials: []string{"public", "qa"},
				},
				CodeHosting: &configfile.CodeHosting{
					Platform:       &github,
					OriginHostname: &githubCom,
				},
				SyncStrategy: &configfile.SyncStrategy{
					FeatureBranches:   &configdomain.SyncFeatureStrategyMerge.Name,
					PerennialBranches: &configdomain.SyncPerennialStrategyRebase.Name,
				},
				PushNewbranches:        &boolFalse,
				ShipDeleteRemoteBranch: &boolFalse,
				SyncUpstream:           &boolTrue,
			}
			have := configfile.Encode(give)
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
			give := configfile.ConfigFile{
				Branches: configfile.Branches{
					Main:       "main",
					Perennials: []string{"public", "qa"},
				},
			}
			have := configfile.Encode(give)
			want := `
[branches]
  main = "main"
  perennials = ["public", "qa"]
`[1:]
			must.EqOp(t, want, have)
		})

	})
}
