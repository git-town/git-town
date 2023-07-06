package cli

import (
	"sort"
	"strings"
)

// BranchAncestryConfig defines the configuration values needed by the `cli` package.
type Ancestry interface {
	Roots() []string
	Children(string) []string
}

// PrintableBranchAncestry provides the branch ancestry in CLI printable format.
func PrintableBranchAncestry(ancestry Ancestry) string {
	roots := ancestry.Roots()
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root, ancestry)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branch string, ancestry Ancestry) string {
	result := branch
	childBranches := ancestry.Children(branch)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, ancestry))
	}
	return result
}
