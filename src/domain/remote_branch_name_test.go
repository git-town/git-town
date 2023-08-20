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

	t.Run("Parts", func(t *testing.T) {
		t.Parallel()
		t.Run("branch at the origin remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("origin/branch")
			remote, localBranch := branch.Parts()
			assert.Equal(t, "origin", remote)
			assert.Equal(t, domain.NewLocalBranchName("branch"), localBranch)
		})
		t.Run("branch at the upstream remote", func(t *testing.T) {
			t.Parallel()
			branch := domain.NewRemoteBranchName("upstream/branch")
			remote, localBranch := branch.Parts()
			assert.Equal(t, "upstream", remote)
			assert.Equal(t, domain.NewLocalBranchName("branch"), localBranch)
		})
	})
}
