package undobranches

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
	"github.com/git-town/git-town/v22/internal/undo/undodomain"
)

type LocalBranchChange map[gitdomain.LocalBranchName]undodomain.Change[gitdomain.SHA]

func (self LocalBranchChange) BranchNames() gitdomain.LocalBranchNames {
	if len(self) == 0 {
		return gitdomain.LocalBranchNames{}
	}
	result := make(gitdomain.LocalBranchNames, 0, len(self))
	for branch := range self { // okay to iterate the map in random order because we sort the result below
		result = append(result, branch)
	}
	slice.NaturalSort(result)
	return result
}
