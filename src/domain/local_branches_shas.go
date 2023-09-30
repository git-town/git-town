package domain

import "golang.org/x/exp/maps"

type LocalBranchesSHAs map[LocalBranchName]SHA

func (lbs LocalBranchesSHAs) BranchNames() LocalBranchNames {
	return maps.Keys(lbs)
}
