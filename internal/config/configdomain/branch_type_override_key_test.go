package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
)

func TestBranchTypeOverrideKey(t *testing.T) {
	t.Parallel()
	t.Run("BranchType", func(t *testing.T) {
		t.Parallel()
		key := configdomain.BranchTypeOverrideKey("git-town-branch.foo.branchtype")
		have := key.BranchType()
		want := configdomain.
	})
}
