package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v11/src/domain"
	"github.com/shoenig/test/must"
)

func TestLocalBranchName(t *testing.T) {
	t.Parallel()

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.EmptyLocalBranchName()
			must.True(t, branch.IsEmpty())
		})
		t.Run("branch is not empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewLocalBranchName("branch")
			must.False(t, branch.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		must.NoError(t, err)
		want := `"branch-1"`
		must.EqOp(t, want, string(have))
	})

	t.Run("NewLocalBranchName and String", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		must.EqOp(t, "branch-1", branch.String())
	})

	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		want := domain.NewRemoteBranchName("origin/branch")
		must.EqOp(t, want, branch.TrackingBranch())
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.EmptyLocalBranchName()
		err := json.Unmarshal([]byte(give), &have)
		must.NoError(t, err)
		want := domain.NewLocalBranchName("branch-1")
		must.EqOp(t, want, have)
	})
}
