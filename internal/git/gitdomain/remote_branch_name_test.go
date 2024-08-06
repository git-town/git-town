package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v14/internal/git/gitdomain"
	"github.com/git-town/git-town/v14/test/asserts"
	"github.com/shoenig/test/must"
)

func TestRemoteBranchName(t *testing.T) {
	t.Parallel()

	t.Run("LocalBranchName", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is at the origin remote", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewRemoteBranchName("origin/branch")
			want := gitdomain.NewLocalBranchName("branch")
			must.EqOp(t, want, branch.LocalBranchName())
		})
		t.Run("branch is at the upstream remote", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewRemoteBranchName("upstream/branch")
			want := gitdomain.NewLocalBranchName("branch")
			must.EqOp(t, want, branch.LocalBranchName())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewRemoteBranchName("origin/branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"origin/branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewRemoteBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("valid remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := gitdomain.NewRemoteBranchName("origin/branch")
			must.EqOp(t, "origin/branch", branch.String())
		})
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewRemoteBranchName("branch")
		})
		t.Run("empty branch name", func(t *testing.T) {
			t.Parallel()
			defer asserts.Paniced(t)
			gitdomain.NewRemoteBranchName("")
		})
	})

	t.Run("Parts", func(t *testing.T) {
		t.Parallel()
		remoteBranch := gitdomain.NewRemoteBranchName("origin/branch")
		remote, localBranch := remoteBranch.Parts()
		must.EqOp(t, gitdomain.RemoteOrigin, remote)
		must.EqOp(t, gitdomain.NewLocalBranchName("branch"), localBranch)
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"origin/branch-1"`
		var have gitdomain.RemoteBranchName
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.NewRemoteBranchName("origin/branch-1")
		must.EqOp(t, want, have)
	})
}
