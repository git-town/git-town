package domain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchChange map[gitdomain.RemoteBranchName]Change[gitdomain.SHA]

func (self RemoteBranchChange) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}

func (self RemoteBranchChange) Categorize(branchTypes BranchTypes) (perennialChanges, featureChanges RemoteBranchChange) {
	perennialChanges = RemoteBranchChange{}
	featureChanges = RemoteBranchChange{}
	for branch, change := range self {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			featureChanges[branch] = change
		} else {
			perennialChanges[branch] = change
		}
	}
	return
}
