package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/shoenig/test"
	"github.com/stretchr/testify/assert"
)

func TestLocalBranchName(t *testing.T) {
	t.Parallel()

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.LocalBranchName{}
			test.True(t, branch.IsEmpty())
		})
		t.Run("branch is not empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewLocalBranchName("branch")
			test.False(t, branch.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		test.NoError(t, err)
		want := `"branch-1"`
		assert.Equal(t, want, string(have))
	})

	t.Run("NewLocalBranchName and String", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		test.EqOp(t, "branch-1", branch.String())
	})

	t.Run("TrackingBranch", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		want := domain.NewRemoteBranchName("origin/branch")
		assert.Equal(t, want, branch.TrackingBranch())
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.LocalBranchName{}
		err := json.Unmarshal([]byte(give), &have)
		test.NoError(t, err)
		want := domain.NewLocalBranchName("branch-1")
		assert.Equal(t, want, have)
	})
}
