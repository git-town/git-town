package cli

import (
	"sort"
	"strings"
)

// Lineage defines the configuration values needed by the `cli` package.
type Lineage interface {
	Roots() []string
	Children(string) []string
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
func PrintableBranchTree(branch string, lineage Lineage) string {
	result := branch
	childBranches := lineage.Children(branch)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, lineage))
	}
	return result
}
