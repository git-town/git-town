package undobranches

import (
	"maps"
	"slices"

	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/gohacks/slice"
)

type RemoteBranchesSHAs map[gitdomain.RemoteBranchName]gitdomain.SHA

// BranchNames provides the names of the involved branches as strings.
func (self RemoteBranchesSHAs) BranchNames() gitdomain.RemoteBranchNames {
	result := gitdomain.RemoteBranchNames(slices.Collect(maps.Keys(self)))
	slice.NaturalSort(result)
	return result
}
