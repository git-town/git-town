package undobranches

import (
	"github.com/git-town/git-town/v14/src/git/gitdomain"
	"golang.org/x/exp/maps"
)

type LocalBranchesSHAs map[gitdomain.LocalBranchName]gitdomain.SHA

func (self LocalBranchesSHAs) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
