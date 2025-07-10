package configfile_test

import (
	"os"
	"testing"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/config/configfile"
	"github.com/git-town/git-town/v21/internal/forge/forgedomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/pkg/asserts"
	. "github.com/git-town/git-town/v21/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()

	t.Run("RenderPerennialBranches", func(t *testing.T) {
		t.Parallel()
		t.Run("no perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.LocalBranchNames{}
			have := configfile.RenderPerennialBranches(give)
			want := "[]"
			must.EqOp(t, want, have)
		})
		t.Run("one perennial branch", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one")
			have := configfile.RenderPerennialBranches(give)
			want := `["one"]`
			must.EqOp(t, want, have)
		})
		t.Run("multiple perennial branches", func(t *testing.T) {
			t.Parallel()
			give := gitdomain.NewLocalBranchNames("one", "two")
			have := configfile.RenderPerennialBranches(give)
			want := `["one", "two"]`
			must.EqOp(t, want, have)
		})
	})

	t.Run("RenderTOML", func(t *testing.T) {
		t.Parallel()
		give := exampleConfig()
		have := configfile.RenderTOML(give, "main")
		want := `
# More info around this file at https://www.git-town.com/configuration-file

[branches]
main = "main"
perennials = ["qa", "staging"]
perennial-regex = "perennial-"

[create]
new-branch-type = "prototype"
share-new-branches = "propose"

[hosting]
dev-remote = "origin"
forge-type = "github"
origin-hostname = "forge"

[ship]
delete-tracking-branch = true
strategy = "api"

[sync]
feature-strategy = "merge"
perennial-strategy = "rebase"
prototype-strategy = "compress"
push-hook = true
tags = true
upstream = true
`[1:]
		must.EqOp(t, want, have)
	})

	t.Run("Save", func(t *testing.T) {
		t.Parallel()
		config := exampleConfig()
		err := configfile.Save(config, "main")
		defer os.Remove(configfile.FileName)
		must.NoError(t, err)
		bytes, err := os.ReadFile(configfile.FileName)
		must.NoError(t, err)
		have := string(bytes)
		want := configfile.RenderTOML(config, "main")
		must.EqOp(t, want, have)
	})
}

func exampleConfig() configdomain.NormalConfigData {
	return configdomain.NormalConfigData{
		DevRemote:                "origin",
		ForgeType:                asserts.NoError1(forgedomain.ParseForgeType("github")),
		HostingOriginHostname:    configdomain.ParseHostingOriginHostname("forge"),
		NewBranchType:            Some(configdomain.BranchTypePrototypeBranch),
		PerennialBranches:        gitdomain.NewLocalBranchNames("qa", "staging"),
		PerennialRegex:           asserts.NoError1(configdomain.ParsePerennialRegex("perennial-")),
		PushHook:                 true,
		ShareNewBranches:         configdomain.ShareNewBranchesPropose,
		ShipDeleteTrackingBranch: true,
		ShipStrategy:             configdomain.ShipStrategyAPI,
		SyncFeatureStrategy:      configdomain.SyncFeatureStrategyMerge,
		SyncPerennialStrategy:    configdomain.SyncPerennialStrategyRebase,
		SyncPrototypeStrategy:    configdomain.SyncPrototypeStrategyCompress,
		SyncTags:                 true,
		SyncUpstream:             true,
	}
}
