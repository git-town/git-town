package cli

import (
	"strings"

	"github.com/git-town/git-town/v9/src/domain"
)

// Lineage defines the configuration values needed by the `cli` package.
type Lineage interface {
	Roots() domain.LocalBranchNames
	Children(domain.LocalBranchName) domain.LocalBranchNames
}

// PrintableBranchLineage provides the branch lineage in CLI printable format.
func PrintableBranchLineage(lineage Lineage) string {
	roots := lineage.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root, lineage)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branch domain.LocalBranchName, lineage Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, lineage))
	}
	return result
}
