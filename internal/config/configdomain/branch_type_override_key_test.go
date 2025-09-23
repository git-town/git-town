package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchTypeOverrideKey(t *testing.T) {
	t.Parallel()

	t.Run("Branch", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.LocalBranchName("my-branch")
		key := configdomain.NewBranchTypeOverrideKeyForBranch(branch)
		have := key.Branch()
		must.EqOp(t, branch, have)
	})

	t.Run("IsBranchTypeOverrideKey", func(t *testing.T) {
		t.Parallel()
		tests := map[string]bool{
			"git-town-branch.foo.branchtype": true,
			"git-town-branch.foo.parent":     false,
			"git-town.prototype-branches":    false,
		}
		for give, want := range tests {
			have := configdomain.IsBranchTypeOverrideKey(give)
			must.EqOp(t, want, have)
		}
	})

	t.Run("NewBranchTypeOverrideKeyForBranch", func(t *testing.T) {
		t.Parallel()
		have := configdomain.NewBranchTypeOverrideKeyForBranch("my-branch")
		want := configdomain.BranchTypeOverrideKey{
			BranchSpecificKey: configdomain.BranchSpecificKey{
				Key: "git-town-branch.my-branch.branchtype",
			},
		}
		must.EqOp(t, want, have)
	})

	t.Run("ParseBranchTypeOverrideKey", func(t *testing.T) {
		t.Parallel()
		t.Run("is branch type override", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseBranchTypeOverrideKey("git-town-branch.my-branch.branchtype")
			want := Some(configdomain.BranchTypeOverrideKey{
				BranchSpecificKey: configdomain.BranchSpecificKey{
					Key: "git-town-branch.my-branch.branchtype",
				},
			})
			must.Eq(t, want, have)
		})
		t.Run("is parent entry", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseBranchTypeOverrideKey("git-town-branch.my-branch.parent")
			want := None[configdomain.BranchTypeOverrideKey]()
			must.Eq(t, want, have)
		})
		t.Run("is something else", func(t *testing.T) {
			t.Parallel()
			have := configdomain.ParseBranchTypeOverrideKey("git-town.feature-regex")
			want := None[configdomain.BranchTypeOverrideKey]()
			must.Eq(t, want, have)
		})
	})
}
