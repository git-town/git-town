package domain

import "golang.org/x/exp/maps"

type RemoteBranchesSHAs map[RemoteBranchName]SHA

func (rbs RemoteBranchesSHAs) Categorize(branchTypes BranchTypes) (perennials, features RemoteBranchesSHAs) {
	perennials = RemoteBranchesSHAs{}
	features = RemoteBranchesSHAs{}
	for branch, sha := range rbs {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}

// BranchNames provides the names of the involved branches as strings.
func (rbs RemoteBranchesSHAs) BranchNames() RemoteBranchNames {
	return maps.Keys(rbs)
}
