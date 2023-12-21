package domain

import "github.com/git-town/git-town/v11/src/git/gitdomain"

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
