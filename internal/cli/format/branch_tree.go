package format

import (
	"github.com/git-town/git-town/v17/internal/config/configdomain"
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
)

// BranchTree provids a printable version of the given branch tree.
func BranchTree(branch gitdomain.LocalBranchName, lineage configdomain.Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(BranchTree(childBranch, lineage))
	}
	return result
}
