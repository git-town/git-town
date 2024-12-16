package gitdomain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/test/git"
	"github.com/shoenig/test/must"
)

func TestLocalBranchName(t *testing.T) {
	t.Parallel()

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewLocalBranchName and String", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch-1")
		must.EqOp(t, "branch-1", branch.String())
	})

	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		branch := gitdomain.NewLocalBranchName("branch")
		want := gitdomain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, branch.TrackingBranch(git.RemoteOrigin))
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		var have gitdomain.LocalBranchName
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := gitdomain.NewLocalBranchName("branch-1")
		must.EqOp(t, want, have)
	})
}
