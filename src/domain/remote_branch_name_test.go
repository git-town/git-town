package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/git-town/git-town/v9/test/asserts"
	"github.com/stretchr/testify/assert"
)

func TestRemoteBranchName(t *testing.T) {
	t.Parallel()
	t.Run("NewBranchName and String", func(t *testing.T) {
		t.Parallel()
		t.Run("valid remote branch name", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("origin/branch")
			assert.Equal(t, "origin/branch", branch.String())
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

	t.Run("LocalBranchName", func(t *testing.T) {
		t.Parallel()
		t.Run("branch is at the origin remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("origin/branch")
			want := domain.NewLocalBranchName("branch")
			assert.Equal(t, want, branch.LocalBranchName())
		})
		t.Run("branch is at the upstream remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("upstream/branch")
			want := domain.NewLocalBranchName("branch")
			assert.Equal(t, want, branch.LocalBranchName())
		})
	})

	t.Run("Parts", func(t *testing.T) {
		t.Parallel()
		remoteBranch := domain.NewRemoteBranchName("origin/branch")
		remote, localBranch := remoteBranch.Parts()
		assert.Equal(t, domain.OriginRemote, remote)
		assert.Equal(t, domain.NewLocalBranchName("branch"), localBranch)
	})
}
