package domain

type InconsistentChanges []InconsistentChange

// InconsistentChange describes a change where both local and remote branch exist before and after,
// but it's not an OmniChange, i.e. the SHA are different.
type InconsistentChange struct {
	Before BranchInfo
	After  BranchInfo
}

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
