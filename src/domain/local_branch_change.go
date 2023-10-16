package domain

type Change[T any] struct {
	Before T
	After  T
}

type LocalBranchChange map[LocalBranchName]Change[SHA]

func (self LocalBranchChange) Categorize(branchTypes BranchTypes) (changedPerennials, changedFeatures LocalBranchChange) {
	changedPerennials = LocalBranchChange{}
	changedFeatures = LocalBranchChange{}
	for branch, change := range self {
		if branchTypes.IsFeatureBranch(branch) {
			changedFeatures[branch] = change
		} else {
			changedPerennials[branch] = change
		}
	}
	return
}

func (self LocalBranchChange) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, 0, len(self))
	for branch := range self {
		result = append(result, branch)
	}
	result.Sort()
	return result
}
