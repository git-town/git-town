package domain_test

import (
	"encoding/json"
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocalBranchName(t *testing.T) {
	t.Parallel()

	t.Run("IsEmpty", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.LocalBranchName{}
			assert.True(t, branch.IsEmpty())
		})
		t.Run("branch is not empty", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewLocalBranchName("branch")
			assert.False(t, branch.IsEmpty())
		})
	})

	t.Run("MarshalJSON", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		have, err := json.MarshalIndent(branch, "", "  ")
		assert.Nil(t, err)
		want := `"branch-1"`
		assert.Equal(t, want, string(have))
	})

	t.Run("NewLocalBranchName and String", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch-1")
		assert.Equal(t, "branch-1", branch.String())
	})

	t.Run("RemoteBranch", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		want := domain.NewRemoteBranchName("origin/branch")
		assert.Equal(t, want, branch.RemoteBranch())
	})

	t.Run("UnmarshalJSON", func(t *testing.T) {
		t.Parallel()
		give := `"branch-1"`
		have := domain.LocalBranchName{}
		json.Unmarshal([]byte(give), &have)
		want := domain.NewLocalBranchName("branch-1")
		assert.Equal(t, want, have)
	})
}
