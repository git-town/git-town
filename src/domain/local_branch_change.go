package domain

import "github.com/git-town/git-town/v11/src/git/gitdomain"

type Change[T any] struct {
	Before T
	After  T
}

type LocalBranchChange map[gitdomain.LocalBranchName]Change[gitdomain.SHA]

func (self LocalBranchChange) BranchNames() gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, 0, len(self))
	for branch := range self {
		result = append(result, branch)
	}
	result.Sort()
	return result
}

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
