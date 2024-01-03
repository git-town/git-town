package configfile_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/config/configfile"
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestSave(t *testing.T) {
	t.Parallel()
	t.Run("Encode", func(t *testing.T) {
		t.Parallel()
		mainBranch := gitdomain.LocalBranchName("main")
		newBranchPush := configdomain.NewBranchPush(false)
		perennialBranches := gitdomain.NewLocalBranchNames("one", "two")
		pushHook := configdomain.PushHook(true)
		shipDeleteTrackingBranch := configdomain.ShipDeleteTrackingBranch(false)
		syncBeforeShip := configdomain.SyncBeforeShip(false)
		syncFeatureStrategy := configdomain.SyncFeatureStrategyMerge
		syncPerennialStrategy := configdomain.SyncPerennialStrategyRebase
		syncUpstream := configdomain.SyncUpstream(true)
		config := configdomain.PartialConfig{
			Aliases:                  map[configdomain.AliasableCommand]string{},
			MainBranch:               &mainBranch,
			NewBranchPush:            &newBranchPush,
			PerennialBranches:        &perennialBranches,
			PushHook:                 &pushHook,
			ShipDeleteTrackingBranch: &shipDeleteTrackingBranch,
			SyncBeforeShip:           &syncBeforeShip,
			SyncFeatureStrategy:      &syncFeatureStrategy,
			SyncPerennialStrategy:    &syncPerennialStrategy,
			SyncUpstream:             &syncUpstream,
		}
		have := configfile.Encode(&config)
		want := `
push-new-branches = false
ship-delete-remote-branch = false
sync-upstream = true

[branches]
  main = "main"
  perennials = ["one", "two"]

[sync-strategy]
  feature-branches = "merge"
  perennial-branches = "rebase"
`[1:]
		must.EqOp(t, want, have)
	})
}
