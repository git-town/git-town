package domain_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/domain"
	"github.com/stretchr/testify/assert"
)

func TestBranchName(t *testing.T) {
	t.Parallel()
	t.Run("NewBranchName", func(t *testing.T) {
		t.Parallel()
		t.Run("empty branch name", func(t *testing.T) {
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
}
