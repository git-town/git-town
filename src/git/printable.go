package git

import (
	"sort"
	"strconv"
	"strings"

	"github.com/git-town/git-town/src/util"
)

var noneString = "[none]"

// getBranchAncestryRoots returns the branches with children and no parents
func getBranchAncestryRoots() []string {
	parentMap := Config().GetParentBranchMap()
	children := make([]string, len(parentMap))
	parents := make([]string, len(parentMap))
	i := 0
	for child, parent := range parentMap {
		children[i] = child
		parents[i] = parent
		i++
	}
	roots := make([]string, len(parentMap))
	i = 0
	for _, parent := range parents {
		if !util.DoesStringArrayContain(children, parent) && !util.DoesStringArrayContain(roots, parent) {
			roots[i] = parent
			i++
		}
	}
	roots = roots[0:i]
	sort.Strings(roots)
	return roots
}

// GetPrintableMainBranch returns a user printable main branch
func GetPrintableMainBranch() string {
	output := Config().GetMainBranch()
	if output == "" {
		return noneString
	}
	return output
}

// GetPrintablePerennialBranches returns a user printable list of perennial branches
func GetPrintablePerennialBranches() string {
	output := strings.Join(Config().GetPerennialBranches(), "\n")
	if output == "" {
		return noneString
	}
	return output
}

// GetPrintableNewBranchPushFlag returns a user printable new branch push flag
func GetPrintableNewBranchPushFlag() string {
	return strconv.FormatBool(Config().ShouldNewBranchPush())
}

// getPrintableBranchTree returns a user printable branch tree
func getPrintableBranchTree(branchName string) (result string) {
	result += branchName
	childBranches := Config().GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(getPrintableBranchTree(childBranch), 1)
	}
	return
}

// GetPrintableBranchAncestry returns a user printable branch ancestry
func GetPrintableBranchAncestry() string {
	roots := getBranchAncestryRoots()
	trees := make([]string, len(roots))
	i := 0
	for _, root := range roots {
		trees[i] = getPrintableBranchTree(root)
		i++
	}
	return strings.Join(trees, "\n\n")
}

// GetPrintableOfflineFlag returns a user printable offline flag
func GetPrintableOfflineFlag() string {
	return strconv.FormatBool(Config().IsOffline())
}
