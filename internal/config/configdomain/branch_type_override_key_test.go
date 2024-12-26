package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v17/internal/config/configdomain"
	. "github.com/git-town/git-town/v17/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchTypeOverrideKey(t *testing.T) {
	t.Parallel()

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
