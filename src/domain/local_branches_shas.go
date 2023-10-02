package domain

import "golang.org/x/exp/maps"

type LocalBranchesSHAs map[LocalBranchName]SHA

func (lbss LocalBranchesSHAs) BranchNames() LocalBranchNames {
	result := LocalBranchNames(maps.Keys(lbss))
	result.Sort()
	return result
}
