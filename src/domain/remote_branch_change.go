package domain

import "golang.org/x/exp/maps"

type RemoteBranchChange map[RemoteBranchName]Change[SHA]

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

func (self RemoteBranchChange) BranchNames() RemoteBranchNames {
	result := RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
