package domain

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
	result := make(RemoteBranchNames, 0, len(rbc))
	for branch := range rbc {
		result = append(result, branch)
	}
	return result
}
