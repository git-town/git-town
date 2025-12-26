package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

type TreeNode2 struct {
	Branch   gitdomain.LocalBranchName
	Children []TreeNode2
}

// CalculateTree provides the full lineage tree for the given branch,
// from the perennial root to all leafs that have the given branch as a descendent.
func CalculateTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) TreeNode2 {
	result := TreeNode2{
		Branch:   branch,
		Children: []TreeNode2{},
	}
	return result
}
