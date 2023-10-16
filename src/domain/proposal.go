package domain

type Proposal struct {
	// the number used to identify the proposal on the hosting platform
	Number int

	// name of the target branch ("base") of this proposal
	Target LocalBranchName

	// textual title of the proposal
	Title string

	// whether this proposal can be merged via the API
	MergeWithAPI bool
}
