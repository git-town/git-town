package domain

type InconsistentChanges []InconsistentChange

func (ics InconsistentChanges) Categorize(branchTypes BranchTypes) (perennials, features InconsistentChanges) {
	for _, change := range ics {
		if branchTypes.IsFeatureBranch(change.Before.LocalName) {
			features = append(features, change)
		} else {
			perennials = append(perennials, change)
		}
	}
	return
}

func (ics InconsistentChanges) BranchNames() LocalBranchNames {
	result := make(LocalBranchNames, len(ics))
	for i, change := range ics {
		result[i] = change.Before.LocalName
	}
	return result
}
