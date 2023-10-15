package domain

type InconsistentChanges []InconsistentChange

func (self InconsistentChanges) Categorize(branchTypes BranchTypes) (perennials, features InconsistentChanges) {
	for _, change := range self {
		if branchTypes.IsFeatureBranch(change.Before.LocalName) {
			features = append(features, change)
		} else {
			perennials = append(perennials, change)
		}
	}
	return
}

func (self InconsistentChanges) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, len(self))
	for i, change := range self {
		result[i] = change.Before.LocalName
	}
	return result
}
