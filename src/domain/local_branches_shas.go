package domain

import "golang.org/x/exp/maps"

type LocalBranchesSHAs map[LocalBranchName]SHA

func (self LocalBranchesSHAs) BranchNames() LocalBranchNames {
	result := LocalBranchNames(maps.Keys(self))
	result.Sort()
	return result
}
