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
	// Find the root ancestor of the given branch
	root := lineage.Root(branch)

	// Get all branches that should be in the tree:
	// - All ancestors from root to branch
	// - The branch itself
	// - All descendants of the branch
	ancestorsAndBranch := lineage.BranchAndAncestors(branch)
	descendants := lineage.Descendants(branch, configdomain.OrderAsc)
	relevantBranches := append(ancestorsAndBranch, descendants...)

	// Build the tree recursively from the root
	return buildTree2(root, lineage, relevantBranches)
}

// buildTree2 recursively builds a tree node for the given branch and its children,
// filtering to only include branches in the relevantBranches set.
func buildTree2(current gitdomain.LocalBranchName, lineage configdomain.Lineage, relevantBranches gitdomain.LocalBranchNames) TreeNode2 {
	node := TreeNode2{
		Branch:   current,
		Children: []TreeNode2{},
	}

	// Get all children and filter to only relevant ones
	allChildren := lineage.Children(current, configdomain.OrderAsc)
	for _, child := range allChildren {
		if relevantBranches.Contains(child) {
			childNode := buildTree2(child, lineage, relevantBranches)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}
