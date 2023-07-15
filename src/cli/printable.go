package cli

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
)

// PrintableBranchLineage provides the branch lineage in CLI printable format.
func PrintableBranchLineage(lineage config.Lineage) string {
	roots := lineage.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root.Name, lineage)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branch string, lineage config.Lineage) string {
	result := branch
	childBranches := lineage.Children(branch).BranchNames()
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, lineage))
	}
	return result
}
