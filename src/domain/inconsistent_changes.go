package domain

import "github.com/git-town/git-town/v11/src/git/gitdomain"

type InconsistentChanges []InconsistentChange

func (self InconsistentChanges) BranchNames() gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, len(self))
	for i, change := range self {
		result[i] = change.Before.LocalName
	}
	return result
}

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
