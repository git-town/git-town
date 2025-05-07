package undobranches

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	"github.com/git-town/git-town/v20/internal/undo/undodomain"
)

type LocalBranchChange map[gitdomain.LocalBranchName]undodomain.Change[gitdomain.SHA]

func (self LocalBranchChange) BranchNames() gitdomain.LocalBranchNames {
	result := make(gitdomain.LocalBranchNames, 0, len(self))
	for branch := range self {
		result = append(result, branch)
	}
	result.Sort()
	return result
}
