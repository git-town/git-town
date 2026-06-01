package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/pkg/set"
)

// FilterTree removes branches whose type is excluded.
// When it removes a branch in the middle of the tree,
// it keeps the visible descendants and moves them up.
// This can turn one tree into multiple visible roots,
// so the result is a forest data structure.
func FilterTree(
	tree TreeNode,
	excluded set.Set[configdomain.BranchType],
) TreeNodes {
	return filterTreeNode(tree, excluded)
}

func filterTreeNode(
	node TreeNode,
	excluded set.Set[configdomain.BranchType],
) TreeNodes {
	filteredChildren := make(TreeNodes, 0, len(node.Children))
	for _, child := range node.Children {
		filteredChildren = append(filteredChildren, filterTreeNode(child, excluded)...)
	}

	if excluded.Contains(node.BranchType) {
		return filteredChildren
	}

	return TreeNodes{{
		Branch:        node.Branch,
		BranchType:    node.BranchType,
		LineageParent: node.LineageParent,
		Children:      filteredChildren,
	}}
}
