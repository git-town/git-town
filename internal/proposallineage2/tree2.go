package proposallineage2

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
func CalculateTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, order configdomain.Order) TreeNode2 {
	ancestorsAndBranch := lineage.BranchAndAncestors(branch)
	root := ancestorsAndBranch[0]
	descendants := lineage.Descendants(branch, order)
	relevantBranches := append(ancestorsAndBranch, descendants...)
	return buildTree2(root, lineage, relevantBranches, order)
}

// buildTree2 provides the Tree2 for the given branch and all its descendents.
func buildTree2(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, includeBranches gitdomain.LocalBranchNames, order configdomain.Order) TreeNode2 {
	children := []TreeNode2{}
	for _, child := range lineage.Children(branch, order) {
		if includeBranches.Contains(child) {
			childNode := buildTree2(child, lineage, includeBranches, order)
			children = append(children, childNode)
		}
	}
	return TreeNode2{
		Branch:   branch,
		Children: children,
	}
}
