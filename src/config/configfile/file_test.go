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
		t.Run("valid content", func(t *testing.T) {
			t.Parallel()
			give := `
push-new-branches = true
ship-delete-remote-branch = false
syncUpstream = true

[branches]
main = "main"
perennials = [ "public", "release" ]

[code-hosting]
platform = "github"
originHostname = "github.com"

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"
`[1:]
			have, err := configfile.Parse(give)
			must.NoError(t, err)
			want := configfile.ConfigFile{
				Branches: configfile.Branches{
					Main:       "main",
					Perennials: []string{"public", "release"},
				},
				CodeHosting: configfile.CodeHosting{
					Platform:       "github",
					OriginHostname: "github.com",
				},
				SyncStrategy: configfile.SyncStrategy{
					FeatureBranches:   "merge",
					PerennialBranches: "rebase",
				},
				PushNewbranches:        true,
				ShipDeleteRemoteBranch: false,
				SyncUpstream:           true,
			}
			must.Eq(t, want, *have)
		})
	})

	t.Run("Encode", func(t *testing.T) {
		t.Parallel()
		give := configfile.ConfigFile{
			Branches: configfile.Branches{
				Main: "main",
			},
			CodeHosting:            configfile.CodeHosting{},
			SyncStrategy:           configfile.SyncStrategy{},
			PushNewbranches:        false,
			ShipDeleteRemoteBranch: false,
			SyncUpstream:           true,
		}
		have := configfile.Encode(give)
		want := `
push-new-branches = false
ship-delete-remote-branch = false
SyncUpstream = true

[Branches]
  Main = "main"

[code-hosting]
  Platform = ""
  OriginHostname = ""

[sync-strategy]
  feature-branches = ""
  perennial-branches = ""
`[1:]
		must.EqOp(t, want, have)
	})
}
