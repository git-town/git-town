package git

import (
	"strconv"
	"strings"

	"github.com/Originate/git-town/lib/util"
)

func GetPrintableMainBranch() string {
	output := GetMainBranch()
	if output == "" {
		return "[none]"
	}
	return output
}

func GetPrintablePerennialBranches() string {
	output := strings.Join(GetPerennialBranches(), "\n")
	if output == "" {
		return "[none]"
	}
	return output
}

func GetPrintableHackPushFlag() string {
	return strconv.FormatBool(ShouldHackPush())
}

func GetPrintableBranchTree(branchName string) (result string) {
	result += branchName
	for _, childBranch := range GetChildBranches(branchName) {
		result += "\n" + util.Indent(GetPrintableBranchTree(childBranch), 1)
	}
	return
}
