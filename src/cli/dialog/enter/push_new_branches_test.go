package enter_test

import (
	"testing"

	"github.com/git-town/git-town/v11/src/cli/dialog/enter"
	"github.com/shoenig/test/must"
)

func TestPushNewBranches(t *testing.T) {
	t.Parallel()

	t.Run("pushNewBranchesEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", enter.PushNewBranchesEntryYes.Short())
			must.Eq(t, "no", enter.PushNewBranchesEntryNo.Short())
		})
	})
}
