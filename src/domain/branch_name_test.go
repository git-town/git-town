package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/shoenig/test"
)

func TestBranchName(t *testing.T) {
	t.Parallel()

	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("normal branch name", func(t *testing.T) {
			branchName := domain.NewBranchName("branch-1")
			test.EqOp(t, "branch-1", branchName.String())
		})
		t.Run("does not allow empty branch names", func(t *testing.T) {
			defer asserts.Paniced(t)
			domain.NewBranchName("")
		})
	})

	t.Run("IsLocal", func(t *testing.T) {
		t.Parallel()
		t.Run("local branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("main")
			test.True(t, branch.IsLocal())
		})
		t.Run("remote branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/main")
			test.False(t, branch.IsLocal())
		})
	})

	t.Run("LocalName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("branch-1")
			want := domain.NewLocalBranchName("branch-1")
			test.EqOp(t, want, branch.LocalName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/branch-1")
			want := domain.NewLocalBranchName("branch-1")
			test.EqOp(t, want, branch.LocalName())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		test.NoError(t, err)
		want := `"branch-1"`
		test.EqOp(t, want, string(have))
	})

	t.Run("RemoteName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("branch-1")
			want := domain.NewRemoteBranchName("origin/branch-1")
			test.EqOp(t, want, branch.RemoteName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/branch-1")
			want := domain.NewRemoteBranchName("origin/branch-1")
			test.EqOp(t, want, branch.RemoteName())
		})
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.BranchName{}
		err := json.Unmarshal([]byte(give), &have)
		test.NoError(t, err)
		want := domain.NewBranchName("branch-1")
		test.EqOp(t, want, have)
	})
}
