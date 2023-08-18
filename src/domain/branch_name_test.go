package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/stretchr/testify/assert"
)

func TestBranchName(t *testing.T) {
	t.Parallel()
	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("normal branch name", func(t *testing.T) {
			branchName := domain.NewBranchName("branch-1")
			assert.Equal(t, "branch-1", branchName.String())
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
			assert.True(t, branch.IsLocal())
		})
		t.Run("remote branch", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/main")
			assert.False(t, branch.IsLocal())
		})
	})

	t.Run("LocalName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("branch-1")
			want := domain.NewLocalBranchName("branch-1")
			assert.Equal(t, want, branch.LocalName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/branch-1")
			want := domain.NewLocalBranchName("branch-1")
			assert.Equal(t, want, branch.LocalName())
		})
	})

	t.Run("RemoteName", func(t *testing.T) {
		t.Run("local branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("branch-1")
			want := domain.NewRemoteBranchName("origin/branch-1")
			assert.Equal(t, want, branch.RemoteName())
		})
		t.Run("remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewBranchName("origin/branch-1")
			want := domain.NewRemoteBranchName("origin/branch-1")
			assert.Equal(t, want, branch.RemoteName())
		})
	})
}
