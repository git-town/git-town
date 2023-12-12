package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/git-town/git-town/v11/test/asserts"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchName(t *testing.T) {
	t.Parallel()

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("is empty", func(t *testing.T) {
			give := domain.EmptyRemoteBranchName()
			must.True(t, give.IsEmpty())
		})
		t.Run("is not empty", func(t *testing.T) {
			give := domain.NewRemoteBranchName("origin/branch-1")
			must.False(t, give.IsEmpty())
		})
	})

	t.Run("LocalBranchName", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is at the origin remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("origin/branch")
			want := domain.NewLocalBranchName("branch")
			must.EqOp(t, want, branch.LocalBranchName())
		})
		t.Run("branch is at the upstream remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("upstream/branch")
			want := domain.NewLocalBranchName("branch")
			must.EqOp(t, want, branch.LocalBranchName())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewRemoteBranchName("origin/branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"origin/branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("valid remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("origin/branch")
			must.EqOp(t, "origin/branch", branch.String())
		})
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewRemoteBranchName("branch")
		})
		t.Run("empty branch name", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			domain.NewRemoteBranchName("")
		})
	})

	t.Run("Parts", func(t *testing.T) {
		t.Parallel()
		remoteBranch := domain.NewRemoteBranchName("origin/branch")
		remote, localBranch := remoteBranch.Parts()
		must.EqOp(t, domain.OriginRemote, remote)
		must.EqOp(t, domain.NewLocalBranchName("branch"), localBranch)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"origin/branch-1"`
		have := domain.EmptyRemoteBranchName()
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := domain.NewRemoteBranchName("origin/branch-1")
		must.EqOp(t, want, have)
	})
}
