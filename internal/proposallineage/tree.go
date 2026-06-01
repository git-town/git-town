package proposallineage

import (
	"github.com/git-town/git-town/v23/internal/config/configdomain"
	"github.com/git-town/git-town/v23/internal/git/gitdomain"
	. "github.com/git-town/git-town/v23/pkg/prelude"
)

type TreeNode struct {
	Branch        gitdomain.LocalBranchName
	BranchType    configdomain.BranchType
	Children      TreeNodes
	LineageParent Option[gitdomain.LocalBranchName]
}

func (self TreeNode) BranchCount() int {
	result := 1
	for _, child := range self.Children {
		result += child.BranchCount()
	}
	return result
}

// CalculateTree provides the full lineage tree for the given branch,
// from the perennial root to all leafs that have the given branch as a descendent.
func CalculateTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage, order configdomain.Order, branchTypes configdomain.BranchesAndTypes) TreeNode {
	ancestorsAndBranch := lineage.BranchAndAncestors(branch)
	root := ancestorsAndBranch[0]
	descendants := lineage.Descendants(branch, order)
	relevantBranches := append(ancestorsAndBranch, descendants...)
	return buildTree(root, None[gitdomain.LocalBranchName](), lineage, relevantBranches, order, branchTypes)
}

// buildTree provides the TreeNodes for the given branch and all its descendents.
func buildTree(branch gitdomain.LocalBranchName, lineageParent Option[gitdomain.LocalBranchName], lineage configdomain.Lineage, includeBranches gitdomain.LocalBranchNames, order configdomain.Order, branchTypes configdomain.BranchesAndTypes) TreeNode {
	children := make(TreeNodes, 0)
	for _, child := range lineage.Children(branch, order) {
		if includeBranches.Contains(child) {
			childNode := buildTree(child, Some(branch), lineage, includeBranches, order, branchTypes)
			children = append(children, childNode)
		}
	}
	return TreeNode{
		Branch:        branch,
		BranchType:    branchTypeFor(branch, branchTypes),
		Children:      children,
		LineageParent: lineageParent,
	}
}

func branchTypeFor(branch gitdomain.LocalBranchName, branchTypes configdomain.BranchesAndTypes) configdomain.BranchType {
	branchType, hasBranchType := branchTypes[branch]
	if !hasBranchType {
		return ""
	}
	return branchType
}
