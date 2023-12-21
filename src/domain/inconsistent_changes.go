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
