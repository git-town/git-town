package git

import (
	"github.com/Originate/git-town/lib/util"
)

// GetBranchSha returns the SHA1 of the latest commit
// on the branch with the given name.
func GetBranchSha(branchName string) string {
	return util.GetCommandOutput("git", "rev-parse", branchName)
}

// GetCurrentSha returns the SHA of the currently checked out commit.
func GetCurrentSha() string {
	return GetBranchSha("HEAD")
}
