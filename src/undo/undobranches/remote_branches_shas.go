package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type RemoteBranchesSHAs map[gitdomain.RemoteBranchName]gitdomain.SHA

// BranchNames provides the names of the involved branches as strings.
func (self RemoteBranchesSHAs) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
