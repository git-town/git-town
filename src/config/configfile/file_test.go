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
pushNewBranches = true
shipDeleteRemoteBranch = false
syncUpstream = true

[branches]
main = "main"
perennial = [ "public", "release" ]

[code-hosting]
platform = "github"
origin-hostname = "github.com"

[sync-strategy]
feature-branches = "merge"
perennial-branches = "rebase"

`[1:]
			have, err := configfile.Parse(give)
			must.NoError(t, err)
			want := configfile.Config{
				Branches: configfile.Branches{
					Main:       "main",
					Perennials: []string{"public", "release"},
				},
				CodeHosting: configfile.CodeHosting{
					Platform:       "github",
					OriginHostname: "github.com",
				},
				SyncStrategy: configfile.SyncStrategy{
					FeatureBranches:   configdomain.SyncFeatureStrategyMerge,
					PerennialBranches: configdomain.SyncPerennialStrategyRebase,
				},
				PushNewbranches:        true,
				ShipDeleteRemoteBranch: false,
				SyncUpstream:           true,
			}
			must.Eq(t, want, *have)
		})
	})
}
