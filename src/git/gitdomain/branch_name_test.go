package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/shoenig/test/must"
)

func TestBranchName(t *testing.T) {
	t.Parallel()

	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("normal branch name", func(t *testing.T) {
			branchName := gitdomain.NewBranchName("branch-1")
			must.EqOp(t, "branch-1", branchName.String())
		})
		t.Run("does not allow empty branch names", func(t *testing.T) {
			defer asserts.Paniced(t)
			gitdomain.NewBranchName("")
		})
	})

	t.Run("IsLocal", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("main")
			must.True(t, branch.IsLocal())
		})
		t.Run("remote branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("origin/main")
			must.False(t, branch.IsLocal())
		})
	})

	t.Run("LocalName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("branch-1")
			want := gitdomain.NewLocalBranchName("branch-1")
			must.EqOp(t, want, branch.LocalName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("origin/branch-1")
			want := gitdomain.NewLocalBranchName("branch-1")
			must.EqOp(t, want, branch.LocalName())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("RemoteName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("branch-1")
			want := gitdomain.NewRemoteBranchName("origin/branch-1")
			must.EqOp(t, want, branch.RemoteName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewBranchName("origin/branch-1")
			want := gitdomain.NewRemoteBranchName("origin/branch-1")
			must.EqOp(t, want, branch.RemoteName())
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := gitdomain.NewBranchName("placeholder")
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.NewBranchName("branch-1")
		must.EqOp(t, want, have)
	})
}
