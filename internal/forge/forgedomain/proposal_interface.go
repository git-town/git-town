package forgedomain

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// ProposalInterface provides information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type ProposalInterface interface {
	// text of the body of the proposal
	// if Some, the string is guaranteed to be non-empty
	GetBody() Option[string]

	// whether this proposal can be merged via the API
	GetMergeWithAPI() bool

	// the number used to identify the proposal on the forge
	GetNumber() int

	// name of the source branch ("head") of this proposal
	GetSource() gitdomain.LocalBranchName

	// name of the target branch ("base") of this proposal
	GetTarget() gitdomain.LocalBranchName

	// text of the title of the proposal
	GetTitle() string

	// the URL of this proposal
	GetURL() string
}

func CommitBody(proposal ProposalInterface, title string) string {
	result := title
	if body, has := proposal.GetBody().Get(); has {
		result += "\n\n"
		result += body
	}
	return result
}
