package common

import "github.com/git-town/git-town/v9/src/domain"

type Proposal struct {
	// the number used to identify the proposal on the hosting platform
	Number int

	// name of the target branch ("base") of this proposal
	Target domain.LocalBranchName

	// textual title of the proposal
	Title string

	// whether this proposal can be merged via the API
	MergeWithAPI bool
}

func (self Proposal) GetNumber() int {
	return self.Number
}

func (self Proposal) GetTarget() domain.LocalBranchName {
	return self.Target
}

func (self Proposal) GetTitle() string {
	return self.Title
}

func (self Proposal) CanMergeWithAPI() bool {
	return self.MergeWithAPI
}
