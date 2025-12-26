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
	root := lineage.Root(branch)
	ancestorsAndBranch := lineage.BranchAndAncestors(branch)
	descendants := lineage.Descendants(branch, configdomain.OrderAsc)
	relevantBranches := append(ancestorsAndBranch, descendants...)
	return buildTree2(root, lineage, relevantBranches)
}

// builds a tree2 from the given root that contains only the given relevant branches
func buildTree2(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, includeBranches gitdomain.LocalBranchNames) TreeNode2 {
	node := TreeNode2{
		Branch:   branch,
		Children: []TreeNode2{},
	}
	for _, child := range lineage.Children(branch, configdomain.OrderAsc) {
		if includeBranches.Contains(child) {
			childNode := buildTree2(child, lineage, includeBranches)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}
