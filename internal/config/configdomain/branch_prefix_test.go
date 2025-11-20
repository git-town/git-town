package configdomain_test

import (
	"testing"

	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/shoenig/test/must"
)

func TestBranchPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Apply", func(t *testing.T) {
		t.Parallel()

		t.Run("empty prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("")
			branch := gitdomain.NewLocalBranchName("feature")
			have := prefix.Apply(branch)
			want := gitdomain.NewLocalBranchName("feature")
			must.EqOp(t, want, have)
		})

		t.Run("non-empty prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("prefix-")
			branch := gitdomain.NewLocalBranchName("feature")
			have := prefix.Apply(branch)
			want := gitdomain.NewLocalBranchName("prefix-feature")
			must.EqOp(t, want, have)
		})

		t.Run("branch already contains the prefix", func(t *testing.T) {
			t.Parallel()
			prefix := configdomain.BranchPrefix("prefix")
			branch := gitdomain.NewLocalBranchName("prefix-branch")
			have := prefix.Apply(branch)
			want := gitdomain.NewLocalBranchName("prefix-branch")
			must.EqOp(t, want, have)
		})
	})
}
