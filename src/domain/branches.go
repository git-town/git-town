package domain

type Branches struct {
	All        BranchInfos
	Perennials BranchTypes
	Initial    LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:        BranchInfos{},
		Perennials: EmptyBranchTypes(),
		Initial:    LocalBranchName{},
	}
}
