package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/internal/gohacks/stringss"
	. "github.com/git-town/git-town/v23/pkg/prelude"
	"github.com/shoenig/test/must"
)

func TestBranchPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Apply", func(t *testing.T) {
		t.Parallel()

		t.Run("empty prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("")
			branch := gitdomain.LocalBranchNameOrPanic("feature")
			have := prefix.Apply(branch)
			want := gitdomain.LocalBranchNameOrPanic("feature")
			must.EqOp(t, want, have)
		})

		t.Run("non-empty prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("prefix-")
			branch := gitdomain.LocalBranchNameOrPanic("feature")
			have := prefix.Apply(branch)
			want := gitdomain.LocalBranchNameOrPanic("prefix-feature")
			must.EqOp(t, want, have)
		})

		t.Run("branch already contains the prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("prefix")
			branch := gitdomain.LocalBranchNameOrPanic("prefix-branch")
			have := prefix.Apply(branch)
			want := gitdomain.LocalBranchNameOrPanic("prefix-branch")
			must.EqOp(t, want, have)
		})
	})

	t.Run("ParseBranchPrefix", func(t *testing.T) {
		t.Parallel()
		tests := map[stringss.TrimmedString]Option[configdomain.BranchPrefix]{
			"":        None[configdomain.BranchPrefix](),
			"prefix-": Some(configdomain.BranchPrefix("prefix-")),
		}
		for give, want := range tests {
			have, err := configdomain.ParseBranchPrefix(give, "test")
			must.NoError(t, err)
			must.Eq(t, want, have)
		}
	})
}
