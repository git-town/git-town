package hosting

import "github.com/git-town/git-town/v9/src/domain"

// Proposal contains information about a change request on a code hosting platform.
// Alternative names are "pull request" or "merge request".
type Proposal interface {
	GetNumber() int

	// name of the target branch ("base") of this proposal
	GetTarget() domain.LocalBranchName

	// textual title of the proposal
	GetTitle() string

	// whether this proposal can be merged via the API
	CanMergeWithAPI() bool
}
