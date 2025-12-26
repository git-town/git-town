package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

type TreeNode2 struct {
	Branch   gitdomain.LocalBranchName
	Children []TreeNode2
}

// NewTree2 provides the full lineage tree for the given branch,
// from the perennial root to all leafs that have the given branch as a descendent.
func NewTree2(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) TreeNode2 {
	// add ancestorBranches
	ancestorBranches := lineage.Ancestors(branch)
	rootBranch, ancestorsWithoutRoot := splitAncestorBranches(ancestorBranches)
	rootNode := TreeNode2{
		Branch:   rootBranch,
		Children: []TreeNode2{},
	}

	// add descendents
}

func splitAncestorBranches(ancestorBranches gitdomain.LocalBranchNames) (gitdomain.LocalBranchName, gitdomain.LocalBranchNames) {
	return ancestorBranches[0], ancestorBranches[1:]
}

func AncestorNodes(branch gitdomain.LocalBranchName, lineage configdomain.Lineage)

// Descendents provides the tree
func DescendentTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) TreeNode2 {
	//
}
