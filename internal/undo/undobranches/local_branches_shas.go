package undobranches

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/gohacks/slice"
)

type LocalBranchesSHAs map[gitdomain.LocalBranchName]gitdomain.SHA

func (self LocalBranchesSHAs) BranchNames() gitdomain.LocalBranchNames {
	result := gitdomain.LocalBranchNames(slices.Collect(maps.Keys(self)))
	return slice.NaturalSort(result)
}
