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
	branchTypes configdomain.BranchesAndTypes,
	excluded set.Set[configdomain.BranchType],
) Forest {
	return filterTreeNode(tree, branchTypes, excluded)
}

func filterTreeNode(
	node TreeNode,
	branchTypes configdomain.BranchesAndTypes,
	excluded set.Set[configdomain.BranchType],
) Forest {
	filteredChildren := make(Forest, 0, len(node.Children))
	for _, child := range node.Children {
		filteredChildren = append(filteredChildren, filterTreeNode(child, branchTypes, excluded)...)
	}

	if branchType, hasBranchType := branchTypes[node.Branch]; hasBranchType && excluded.Contains(branchType) {
		return filteredChildren
	}

	return Forest{{
		Branch:        node.Branch,
		LineageParent: node.LineageParent,
		Children:      filteredChildren,
	}}
}
