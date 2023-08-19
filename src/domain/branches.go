package domain

type Branches struct {
	All         BranchInfos
	BranchTypes BranchTypes
	Initial     LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:         BranchInfos{},
		BranchTypes: EmptyBranchTypes(),
		Initial:     LocalBranchName{},
	}
}
