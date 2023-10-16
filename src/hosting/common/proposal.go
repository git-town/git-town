package common

import "github.com/git-town/git-town/v9/src/domain"

// Proposal contains information about a change request
// on a code hosting platform.
// Alternative names are "pull request" or "merge request".
type Proposal struct {
	// the number used to identify the proposal on the hosting platform
	Number int

	// name of the target branch ("base") of this proposal
	Target domain.LocalBranchName

	// textual title of the proposal
	Title string

	// whether this proposal can be merged via the API
	CanMergeWithAPI bool
}
