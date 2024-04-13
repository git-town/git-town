package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"github.com/git-town/git-town/v14/src/undo/undodomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchChange map[gitdomain.RemoteBranchName]undodomain.Change[gitdomain.SHA]

func (self RemoteBranchChange) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
