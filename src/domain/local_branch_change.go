package domain

type LocalBranchChange map[LocalBranchName]Change[SHA]

func (lbc LocalBranchChange) Categorize(branchTypes BranchTypes) (changedPerennials, changedFeatures LocalBranchChange) {
	changedPerennials = LocalBranchChange{}
	changedFeatures = LocalBranchChange{}
	for branch, change := range lbc {
		if branchTypes.IsFeatureBranch(branch) {
			changedFeatures[branch] = change
		} else {
			changedPerennials[branch] = change
		}
	}
	return
}

func (lbc LocalBranchChange) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(lbc))
	for branch := range lbc {
		result = append(result, branch)
	}
	result.Sort()
	return result
}

type Change[T any] struct {
	Before T
	After  T
}
