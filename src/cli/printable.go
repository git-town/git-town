package cli

import (
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/src/util"
)

// PrintableBranchAncestry provides the branch ancestry in CLI printable format.
func PrintableBranchAncestry(config Config) string {
	roots := config.GetBranchAncestryRoots()
	trees := make([]string, len(roots))
	for r := range roots {
		trees[r] = PrintableBranchTree(roots[r], config)
	}
	return strings.Join(trees, "\n\n")
}

// PrintableBranchTree returns a user printable branch tree.
func PrintableBranchTree(branchName string, config Config) (result string) {
	result += branchName
	childBranches := config.GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(PrintableBranchTree(childBranch, config))
	}
	return
}

// PrintableMainBranch returns a user printable main branch.
func PrintableMainBranch(mainBranch string) string {
	if mainBranch == "" {
		return "[none]"
	}
	return mainBranch
}

// PrintableNewBranchPushFlag returns a user printable new branch push flag.
func PrintableNewBranchPushFlag(flag bool) string {
	return strconv.FormatBool(flag)
}

// PrintableOfflineFlag provides a printable version of the given offline flag.
func PrintableOfflineFlag(flag bool) string {
	return strconv.FormatBool(flag)
}

// PrintablePerennialBranches returns a user printable list of perennial branches.
func PrintablePerennialBranches(perennialBranches []string) string {
	if len(perennialBranches) == 0 {
		return "[none]"
	}
	return strings.Join(perennialBranches, "\n")
}
