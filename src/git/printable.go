package git

import (
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

// GetPrintableHackPushFlag returns a user printable hack push flag
func GetPrintableHackPushFlag() string {
	return strconv.FormatBool(ShouldHackPush())
}

// GetPrintableBranchTree returns a user printable branch tree
func GetPrintableBranchTree(branchName string) (result string) {
	result += branchName
	for _, childBranch := range GetChildBranches(branchName) {
		result += "\n" + util.Indent(GetPrintableBranchTree(childBranch), 1)
	}
	return
}
