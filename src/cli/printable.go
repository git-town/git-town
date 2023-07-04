package cli

import (
	"sort"
	"strings"

	"github.com/git-town/git-town/v9/src/config"
)

// BranchAncestryConfig defines the configuration values needed by the `cli` package.
type BranchAncestryConfig interface {
	BranchAncestryRoots(config.BranchParents) []string
	ChildBranches(string) []string
	ParentBranchMap() config.BranchParents
}

// PrintableBranchAncestry provides the branch ancestry in CLI printable format.
func PrintableBranchAncestry(config BranchAncestryConfig) string {
	roots := config.BranchAncestryRoots(config.ParentBranchMap())
	trees := make([]string, len(roots))
	for r, root := range roots {
		trees[r] = PrintableBranchTree(root, config)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branch string, config BranchAncestryConfig) string {
	result := branch
	childBranches := config.ChildBranches(branch)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(PrintableBranchTree(childBranch, config))
	}
	return result
}
