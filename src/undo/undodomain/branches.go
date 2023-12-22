package undodomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
	"github.com/git-town/git-town/v11/src/sync/syncdomain"
)

// TODO: this struct doesn't really belong here. Do we even need it in the new world where most config information is available via a single variable?
type Branches struct {
	All     BranchInfos
	Types   syncdomain.BranchTypes
	Initial gitdomain.LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:     BranchInfos{},
		Types:   syncdomain.EmptyBranchTypes(),
		Initial: gitdomain.EmptyLocalBranchName(),
	}
}
