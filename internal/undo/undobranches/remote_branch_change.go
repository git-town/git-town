package undobranches

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v19/internal/git/gitdomain"
	"github.com/git-town/git-town/v19/internal/undo/undodomain"
)

type RemoteBranchChange map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]

func (self RemoteBranchChange) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(slices.Collect(maps.Keys(self)))
	result.Sort()
	return result
}
