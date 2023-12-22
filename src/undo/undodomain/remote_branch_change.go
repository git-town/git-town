package undodomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchChange map[gitdomain.RemoteBranchName]Change[gitdomain.SHA]

func (self RemoteBranchChange) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
