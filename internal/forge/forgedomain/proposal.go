package forgedomain

import (
	"github.com/git-town/git-town/v20/internal/git/gitdomain"
	. "github.com/git-town/git-town/v20/pkg/prelude"
)

// Proposal provides information about a change request on a forge.
// Alternative names are "pull request" or "merge request".
type Proposal interface {
	// text of the body of the proposal
	// if Some, the string is guaranteed to be non-empty
	Body() Option[string]

	// whether this proposal can be merged via the API
	MergeWithAPI() bool

	// the number used to identify the proposal on the forge
	Number() int

	// name of the source branch ("head") of this proposal
	Source() gitdomain.LocalBranchName

	// name of the target branch ("base") of this proposal
	Target() gitdomain.LocalBranchName

	// text of the title of the proposal
	Title() string

	// the URL of this proposal
	URL() string
}

func CommitBody(proposal Proposal, title string) string {
	result := title
	if body, has := proposal.Body().Get(); has {
		result += "\n\n"
		result += body
	}
	return result
}

type ProposalWrapper struct {
	Type    ProposalType
	Content interface{}
}

type ProposalType string

const (
	ProposalTypeGitHub ProposalType = "github"
)

func (self ProposalType) MarshalJSON() ([]byte, error) {
	//
}
