package domain

import "golang.org/x/exp/maps"

type RemoteBranchChange map[RemoteBranchName]Change[SHA]

func (rbc RemoteBranchChange) Categorize(branchTypes BranchTypes) (perennialChanges, featureChanges RemoteBranchChange) {
	perennialChanges = RemoteBranchChange{}
	featureChanges = RemoteBranchChange{}
	for branch, change := range rbc {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			featureChanges[branch] = change
		} else {
			perennialChanges[branch] = change
		}
	}
	return
}

func (rbc RemoteBranchChange) BranchNames() RemoteBranchNames {
	result := RemoteBranchNames(maps.Keys(rbc))
	result.Sort()
	return result
}
