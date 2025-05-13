package forgedomain

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Proposal contains information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type Proposal struct {
	// text of the body of the proposal
	// if Some, the string is guaranteed to be non-empty
	Body Option[string]

	// whether this proposal can be merged via the API
	MergeWithAPI bool

	// the number used to identify the proposal on the forge
	Number int

	// name of the source branch ("head") of this proposal
	Source gitdomain.LocalBranchName

	// name of the target branch ("base") of this proposal
	Target gitdomain.LocalBranchName

	// text of the title of the proposal
	Title string

	// the URL of this proposal
	URL string
}
