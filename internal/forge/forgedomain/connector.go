package forgedomain

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// Connector describes the activities that Git Town can perform on forges.
// Individual implementations exist to talk to specific forge types.
// Functions that might or might not be supported by a connector are implemented as higher-level functions,
// i.e. they return an option of the function to call.
// A `None` value implies that the respective functionality isn't supported by this connector implementation.
type Connector interface {
	// CreateProposal creates a proposal at the forge.
	CreateProposal(CreateProposalArgs) error

	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective forge type is prepopulated with.
	DefaultProposalMessage(proposal ProposalData) string

	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to load details about the proposal for the given branch into the given target branch.
	// A None return value indicates that this connector does not support this feature (yet).
	FindProposalFn() Option[func(branch, target gitdomain.LocalBranchName) (Option[Proposal], error)]

	// OpenRepository opens this repository in the associated application, typically the browser.
	OpenRepository(runner subshelldomain.Runner) error

	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to search for a proposal that has the given branch as its source branch.
	// A None return value indicates that this connector does not support this feature (yet).
	SearchProposalFn() Option[func(branch gitdomain.LocalBranchName) (Option[Proposal], error)]

	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to merge the proposal with the given number using the given message.
	// A None return value indicates that this connector does not support this feature (yet).
	SquashMergeProposalFn() Option[func(number int, message gitdomain.CommitMessage) error]

	// If the connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the body (description) of the given proposal to the given value.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalBodyFn() Option[func(proposal ProposalInterface, newBody string) error]

	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the source branch of the proposal with the given number.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalSourceFn() Option[func(proposal ProposalInterface, newSource gitdomain.LocalBranchName) error]

	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the target branch of the proposal with the given number.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalTargetFn() Option[func(proposal ProposalInterface, newTarget gitdomain.LocalBranchName) error]

	// VerifyConnection checks whether this connector can make successful requests to the forge.
	VerifyConnection() VerifyConnectionResult
}

type VerifyConnectionResult struct {
	AuthenticatedUser   Option[string] // the authenticated username
	AuthenticationError error          // error while verifying to verify authentication
	AuthorizationError  error          // error while verifying authorization, nil == user is authenticated
}

type CreateProposalArgs struct {
	Branch         gitdomain.LocalBranchName
	FrontendRunner subshelldomain.Runner
	MainBranch     gitdomain.LocalBranchName
	ParentBranch   gitdomain.LocalBranchName
	ProposalBody   Option[gitdomain.ProposalBody]
	ProposalTitle  Option[gitdomain.ProposalTitle]
}
