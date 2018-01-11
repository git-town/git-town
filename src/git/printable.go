package git

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Originate/git-town/src/util"
)

// GetPrintableMainBranch returns a user printable main branch
func GetPrintableMainBranch() string {
	output := GetMainBranch()
	if output == "" {
		return "[none]"
	}
	return output
}

// GetPrintablePerennialBranches returns a user printable list of perennial branches
func GetPrintablePerennialBranches() string {
	output := strings.Join(GetPerennialBranches(), "\n")
	if output == "" {
		return "[none]"
	}
	return output
}

// GetPrintableNewBranchPushFlag returns a user printable new branch push flag
func GetPrintableNewBranchPushFlag() string {
	return strconv.FormatBool(ShouldNewBranchPush())
}

// GetPrintableBranchTree returns a user printable branch tree
func GetPrintableBranchTree(branchName string) (result string) {
	result += branchName
	childBranches := GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(GetPrintableBranchTree(childBranch), 1)
	}
	return
}

// GetPrintableOfflineFlag returns a user printable offline flag
func GetPrintableOfflineFlag() string {
	return strconv.FormatBool(IsOffline())
}
