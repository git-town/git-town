package format

import (
	"github.com/git-town/git-town/v11/src/config/configdomain"
	"github.com/git-town/git-town/v11/src/domain"
)

// BranchTree provids a printable version of the given branch tree.
func BranchTree(branch domain.LocalBranchName, lineage configdomain.Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(BranchTree(childBranch, lineage))
	}
	return result
}
