package format

import (
	"github.com/git-town/git-town/v10/src/config"
	"github.com/git-town/git-town/v10/src/domain"
)

// BranchTree provids a printable version of the given branch tree.
func BranchTree(branch domain.LocalBranchName, lineage config.Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(BranchTree(childBranch, lineage))
	}
	return result
}
