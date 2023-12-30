package configdomain

import (
	"github.com/git-town/git-town/v11/src/git/gitdomain"
)

// TODO: this struct doesn't really belong here. Do we even need it in the new world where most config information is available via a single variable?
type Branches struct {
	All     gitdomain.BranchInfos
	Initial gitdomain.LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:     gitdomain.BranchInfos{},
		Initial: gitdomain.EmptyLocalBranchName(),
	}
}
