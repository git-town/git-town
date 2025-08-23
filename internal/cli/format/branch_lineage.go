package format

import (
	"strings"

	"github.com/git-town/git-town/v21/internal/config/configdomain"
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
)

// BranchLineage provides printable formatting of the given branch lineage.
func BranchLineage(lineage configdomain.Lineage) string {
	roots := lineage.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = branchTree(root, lineage)
	}
	return strings.Join(trees, "\n\n")
}

// branchTree provids a printable version of the given branch tree.
func branchTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(branchTree(childBranch, lineage))
	}
	return result
}
