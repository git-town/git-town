package git

import (
	"sort"

	"github.com/git-town/git-town/src/util"
)

var noneString = "[none]"

// getBranchAncestryRoots returns the branches with children and no parents.
func getBranchAncestryRoots() []string {
	parentMap := Config().GetParentBranchMap()
	roots := []string{}
	for _, parent := range parentMap {
		if _, ok := parentMap[parent]; !ok && !util.DoesStringArrayContain(roots, parent) {
			roots = append(roots, parent)
		}
	}
	sort.Strings(roots)
	return roots
}

// getPrintableBranchTree returns a user printable branch tree.
func getPrintableBranchTree(branchName string) (result string) {
	result += branchName
	childBranches := Config().GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(getPrintableBranchTree(childBranch))
	}
	return
}
