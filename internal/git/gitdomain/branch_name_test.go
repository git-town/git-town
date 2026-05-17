package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	"github.com/git-town/git-town/v23/pkg/asserts"
	"github.com/shoenig/test/must"
)

func TestBranchName(t *testing.T) {
	t.Parallel()

	t.Run("IsLocal", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("main")
			must.True(t, branch.IsLocal())
		})
		t.Run("remote branch", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("origin/main")
			must.False(t, branch.IsLocal())
		})
	})

	t.Run("LocalName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("branch-1")
			want := gitdomain.LocalBranchName("branch-1")
			must.EqOp(t, want, branch.LocalName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("origin/branch-1")
			want := gitdomain.LocalBranchName("branch-1")
			must.EqOp(t, want, branch.LocalName())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.BranchNameOrPanic("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("normal branch name", func(t *testing.T) {
			t.Parallel()
			branchName := gitdomain.BranchNameOrPanic("branch-1")
			must.EqOp(t, "branch-1", branchName.String())
		})
		t.Run("does not allow empty branch names", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.BranchNameOrPanic("")
		})
	})

	t.Run("RefName", func(t *testing.T) {
		t.Parallel()
		tests := map[gitdomain.BranchName]string{
			"main":        "refs/heads/main",
			"origin/main": "origin/main",
		}
		for give, want := range tests {
			must.EqOp(t, want, give.RefName())
		}
	})

	t.Run("RemoteName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("branch-1")
			want := gitdomain.RemoteBranchNameOrPanic("origin/branch-1")
			must.EqOp(t, want, branch.RemoteName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.BranchNameOrPanic("origin/branch-1")
			want := gitdomain.RemoteBranchNameOrPanic("origin/branch-1")
			must.EqOp(t, want, branch.RemoteName())
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := gitdomain.BranchNameOrPanic("placeholder")
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.BranchNameOrPanic("branch-1")
		must.EqOp(t, want, have)
	})
}
