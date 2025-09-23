package undobranches

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
)

type LocalBranchesSHAs map[gitdomain.LocalBranchName]gitdomain.SHA

func (self LocalBranchesSHAs) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(slices.Collect(maps.Keys(self)))
	slice.NaturalSort(result)
	return result
}
