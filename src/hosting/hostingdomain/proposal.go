package hostingdomain

import "github.com/git-town/git-town/v14/src/git/gitdomain"

// Proposal contains information about a change request on a code hosting platform.
// Alternative names are "pull request" or "merge request".
type Proposal struct {
	// whether this proposal can be merged via the API
	MergeWithAPI bool

	// the number used to identify the proposal on the hosting platform
	Number int

	// name of the target branch ("base") of this proposal
	Target gitdomain.LocalBranchName

	// textual title of the proposal
	Title string
}
