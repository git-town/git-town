package undobranches

import (
	"github.com/git-town/git-town/v17/internal/git/gitdomain"
	"github.com/git-town/git-town/v17/internal/undo/undodomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchChange map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]

func (self RemoteBranchChange) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
