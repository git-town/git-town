package proposallineage

import (
	"github.com/git-town/git-town/v22/internal/config/configdomain"
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
)

type TreeNode struct {
	Branch   gitdomain.LocalBranchName
	Children []TreeNode
}

func (self TreeNode) BranchCount() int {
	result := len(self.Children)
	for _, child := range self.Children {
		result += child.BranchCount()
	}
	return result
}

// CalculateTree provides the full lineage tree for the given branch,
// from the perennial root to all leafs that have the given branch as a descendent.
func CalculateTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, order configdomain.Order) TreeNode {
	ancestorsAndBranch := lineage.BranchAndAncestors(branch)
	root := ancestorsAndBranch[0]
	descendants := lineage.Descendants(branch, order)
	relevantBranches := append(ancestorsAndBranch, descendants...)
	return buildTree(root, lineage, relevantBranches, order)
}

// buildTree provides the TreeNodes for the given branch and all its descendents.
func buildTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, includeBranches gitdomain.LocalBranchNames, order configdomain.Order) TreeNode {
	children := []TreeNode{}
	for _, child := range lineage.Children(branch, order) {
		if includeBranches.Contains(child) {
			childNode := buildTree(child, lineage, includeBranches, order)
			children = append(children, childNode)
		}
	}
	return TreeNode{
		Branch:   branch,
		Children: children,
	}
}
