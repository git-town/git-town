package syncdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// TODO: this struct doesn't really belong here. Do we even need it in the new world where most config information is available via a single variable?
type Branches struct {
	All     BranchInfos
	Types   BranchTypes
	Initial gitdomain.LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:     BranchInfos{},
		Types:   EmptyBranchTypes(),
		Initial: gitdomain.EmptyLocalBranchName(),
	}
}
