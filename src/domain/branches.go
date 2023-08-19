package domain

type Branches struct {
	All       BranchInfos
	Durations BranchDurations
	Initial   LocalBranchName
}

// EmptyBranches provides the zero value for Branches.
func EmptyBranches() Branches {
	return Branches{
		All:       BranchInfos{},
		Durations: EmptyBranchDurations(),
		Initial:   LocalBranchName{},
	}
}
