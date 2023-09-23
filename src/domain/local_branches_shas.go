package domain

type LocalBranchesSHAs map[LocalBranchName]SHA

func (lbs LocalBranchesSHAs) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(lbs))
	for branch := range lbs {
		result = append(result, branch)
	}
	return result
}
