package cli

import (
	"sort"
	"strings"
)

// BranchLineageConfig defines the configuration values needed by the `cli` package.
type BranchLineageConfig interface {
	BranchLineageRoots() []string
	ChildBranches(string) []string
}

// PrintableBranchLineage provides the branch lineage in CLI printable format.
func PrintableBranchLineage(config BranchLineageConfig) string {
	roots := config.BranchLineageRoots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root, config)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branch string, config BranchLineageConfig) string {
	result := branch
	childBranches := config.ChildBranches(branch)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, config))
	}
	return result
}
