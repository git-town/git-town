package format

import (
	"github.com/git-town/git-town/v9/src/config"
	"github.com/git-town/git-town/v9/src/domain"
)

// BranchTree returns a user printable branch tree.
func BranchTree(branch domain.LocalBranchName, lineage config.Lineage) string {
	result := branch.String()
	childBranches := lineage.Children(branch)
	for _, childBranch := range childBranches {
		result += "\n" + Indent(BranchTree(childBranch, lineage))
	}
	return result
}
