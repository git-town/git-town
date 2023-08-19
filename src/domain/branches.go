package domain

type Branches struct {
	All     BranchInfos
	Types   BranchTypes
	Initial LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:     BranchInfos{},
		Types:   EmptyBranchTypes(),
		Initial: LocalBranchName{},
	}
}
