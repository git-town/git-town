package domain

import "golang.org/x/exp/maps"

type RemoteBranchesSHAs map[RemoteBranchName]SHA

func (rbss RemoteBranchesSHAs) Categorize(branchTypes BranchTypes) (perennials, features RemoteBranchesSHAs) {
	perennials = RemoteBranchesSHAs{}
	features = RemoteBranchesSHAs{}
	for branch, sha := range rbss {
		if branchTypes.IsFeatureBranch(branch.LocalBranchName()) {
			features[branch] = sha
		} else {
			perennials[branch] = sha
		}
	}
	return
}

// BranchNames provides the names of the involved branches as strings.
func (rbss RemoteBranchesSHAs) BranchNames() RemoteBranchNames {
	result := RemoteBranchNames(maps.Keys(rbss))
	result.Sort()
	return result
}
