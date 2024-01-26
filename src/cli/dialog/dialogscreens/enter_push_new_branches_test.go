package dialogscreens_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog/dialogscreens"
	"github.com/shoenig/test/must"
)

func TestEnterPushNewBranches(t *testing.T) {
	t.Parallel()

	t.Run("pushNewBranchesEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogscreens.PushNewBranchesEntryYes.Short())
			must.Eq(t, "no", dialogscreens.PushNewBranchesEntryNo.Short())
		})
	})
}
