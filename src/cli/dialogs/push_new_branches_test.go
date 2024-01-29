package dialogs_test

import (
	"testing"

	"github.com/shoenig/test/must"
)

func TestPushNewBranches(t *testing.T) {
	t.Parallel()

	t.Run("pushNewBranchesEntry", func(t *testing.T) {
		t.Parallel()
		t.Run("Short", func(t *testing.T) {
			t.Parallel()
			must.Eq(t, "yes", dialogs.PushNewBranchesEntryYes.Short())
			must.Eq(t, "no", dialogs.PushNewBranchesEntryNo.Short())
		})
	})
}
