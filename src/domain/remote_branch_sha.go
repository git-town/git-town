package domain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchesSHAs map[gitdomain.RemoteBranchName]gitdomain.SHA

// BranchNames provides the names of the involved branches as strings.
func (self RemoteBranchesSHAs) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}

func (self RemoteBranchesSHAs) Categorize(branchTypes BranchTypes) (perennials, features RemoteBranchesSHAs) {
	perennials = RemoteBranchesSHAs{}
	features = RemoteBranchesSHAs{}
	for branch, sha := range self {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}
