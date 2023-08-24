package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestLocalBranchNames(t *testing.T) {
	t.Parallel()
	t.Run("AtRemote", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		have := branch.AtRemote(domain.OriginRemote)
		want := domain.NewRemoteBranchName("origin/branch")
		assert.Equal(t, want, have)
	})

	t.Run("NewLocalBranchNames and Strings", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "two", "three"}
		assert.Equal(t, want, branches.Strings())
	})

	t.Run("RemoteName", func(t *testing.T) {
		t.Parallel()
		branch := domain.NewLocalBranchName("branch")
		have := branch.RemoteBranch()
		want := domain.NewRemoteBranchName("origin/branch")
		assert.Equal(t, want, have)
	})

	t.Run("Sort", func(t *testing.T) {
		t.Parallel()
		branches := domain.NewLocalBranchNames("one", "two", "three")
		want := []string{"one", "three", "two"}
		branches.Sort()
		assert.Equal(t, want, branches.Strings())
	})
}
