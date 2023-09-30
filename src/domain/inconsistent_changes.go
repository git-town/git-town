package domain

type InconsistentChanges []InconsistentChange

func (ic InconsistentChanges) Categorize(branchTypes BranchTypes) (perennials, features InconsistentChanges) {
	for _, change := range ic {
		if branchTypes.IsFeatureBranch(change.Before.LocalName) {
			features = append(features, change)
		} else {
			perennials = append(perennials, change)
		}
	}
	return
}

func (ic InconsistentChanges) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, len(ic))
	for i, change := range ic {
		result[i] = change.Before.LocalName
	}
	return result
}
