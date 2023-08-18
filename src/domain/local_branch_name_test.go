package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocalBranchName(t *testing.T) {
	t.Run("NewLocalBranchName and String", func(t *testing.T) {
		branch := domain.NewLocalBranchName("branch-1")
		assert.Equal(t, "branch-1", branch.String())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		t.Run("branch is empty", func(t *testing.T) {
			branch := domain.LocalBranchName{}
			assert.True(t, branch.IsEmpty())
		})
		t.Run("branch is not empty", func(t *testing.T) {
			branch := domain.NewLocalBranchName("branch")
			assert.False(t, branch.IsEmpty())
		})
	})

	t.Run("RemoteName", func(t *testing.T) {
		branch := domain.NewLocalBranchName("branch")
		want := domain.NewRemoteBranchName("origin/branch")
		assert.Equal(t, want, branch.RemoteName())
	})
}

func TestLocalBranchNames(t *testing.T) {
	t.Run("NewLocalBranchNames and Strings", func(t *testing.T) {
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "two", "three"}
		assert.Equal(t, want, branches.Strings())
	})

	t.Run("Sort", func(t *testing.T) {
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "three", "two"}
		branches.Sort()
		assert.Equal(t, want, branches.Strings())
	})
}
