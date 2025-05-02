package undobranches

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v20/internal/git/gitdomain"
)

type LocalBranchesSHAs map[gitdomain.LocalBranchName]gitdomain.SHA

func (self LocalBranchesSHAs) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(slices.Collect(maps.Keys(self)))
	result.Sort()
	return result
}
