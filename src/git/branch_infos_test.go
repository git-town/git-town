package git_test

import (
	"testing"

	"github.com/git-town/git-town/v9/src/git"
	"github.com/stretchr/testify/assert"
)

func TestBranchInfos(t *testing.T) {
	t.Parallel()
	t.Run("IndexOfBranch", func(t *testing.T) {
		t.Parallel()
		t.Run("branch_exists", func(t *testing.T) {
			t.Parallel()
			bi := git.BranchInfos{
				{Name: "branch-1"},
				{Name: "branch-2"},
				{Name: "branch-3"},
			}
			have, found := bi.IndexOfBranch("branch-2")
			assert.True(t, found)
			assert.Equal(t, 1, have)
		})
		t.Run("branch does not exist", func(t *testing.T) {
			t.Parallel()
			bi := git.BranchInfos{
				git.BranchInfo{Name: "branch-1"},
			}
			_, found := bi.IndexOfBranch("branch-2")
			assert.False(t, found)
		})
	})

	t.Run("OrderedHierarchically", func(t *testing.T) {
		t.Parallel()
		bi := git.BranchInfos{
			{Name: "branch-1", Parent: "main"},
			{Name: "branch-1a", Parent: "branch-1"},
			{Name: "branch-1b", Parent: "branch-1"},
			{Name: "branch-1a1", Parent: "branch-1a"},
			{Name: "branch-1a1a", Parent: "branch-1a1"},
			{Name: "branch-2", Parent: "main"},
			{Name: "main", Parent: ""},
		}
		have := bi.OrderedHierarchically()
		want := git.BranchInfos{
			{Name: "main", Parent: ""},
			{Name: "branch-1", Parent: "main"},
			{Name: "branch-2", Parent: "main"},
			{Name: "branch-1a", Parent: "branch-1"},
			{Name: "branch-1b", Parent: "branch-1"},
			{Name: "branch-1a1", Parent: "branch-1a"},
			{Name: "branch-1a1a", Parent: "branch-1a1"},
		}
		assert.Equal(t, want, have)
	})
}
