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

// GetPrintableNewBranchPushFlag returns a user printable new branch push flag
func GetPrintableNewBranchPushFlag() string {
	return strconv.FormatBool(ShouldNewBranchPush())
}

// GetPrintableBranchTree returns a user printable branch tree
func GetPrintableBranchTree(branchName string, parentBranchMap map[string]string) (result string) {
	result += branchName
	for child, parent := range parentBranchMap {
		if parent == branchName {
			result += "\n" + util.Indent(GetPrintableBranchTree(child, parentBranchMap), 1)
		}
	}
	return
}

// GetPrintableOfflineFlag returns a user printable offline flag
func GetPrintableOfflineFlag() string {
	return strconv.FormatBool(IsOffline())
}
