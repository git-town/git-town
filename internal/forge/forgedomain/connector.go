package forgedomain

import (
	"github.com/git-town/git-town/v21/internal/git/gitdomain"
	"github.com/git-town/git-town/v21/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v21/pkg/prelude"
)

// This file defines an interface for every capability that a forge and connector implementation can support.

// Connector describes the capabilities that all connectors can perform on forges.
type Connector interface {
	// CreateProposal creates a proposal at the forge.
	CreateProposal(CreateProposalArgs) error

	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective forge type is prepopulated with.
	DefaultProposalMessage(proposal ProposalData) string

	// OpenRepository opens this repository in the associated application, typically the browser.
	OpenRepository(runner subshelldomain.Runner) error
}

type CreateProposalArgs struct {
	Branch         gitdomain.LocalBranchName
	FrontendRunner subshelldomain.Runner
	MainBranch     gitdomain.LocalBranchName
	ParentBranch   gitdomain.LocalBranchName
	ProposalBody   Option[gitdomain.ProposalBody]
	ProposalTitle  Option[gitdomain.ProposalTitle]
}

// CredentialVerifier describes capabilities to verify credentials.
type CredentialVerifier interface {
	// VerifyConnection checks whether this connector can make successful requests to the forge.
	VerifyCredentials() VerifyCredentialsResult
}

type VerifyCredentialsResult struct {
	AuthenticatedUser   Option[string] // the authenticated username
	AuthenticationError error          // error while verifying to verify authentication
	AuthorizationError  error          // error while verifying authorization, nil == user is authenticated
}

// ProposalFinder describes methods that connectors need to implement
// to enable Git Town to find proposals at the active forge.
type ProposalFinder interface {
	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to load details about the proposal for the given branch into the given target branch.
	// A None return value indicates that this connector does not support this feature (yet).
	FindProposal(branch, target gitdomain.LocalBranchName) (Option[Proposal], error)
}

// ProposalSearcher describes methods that connectors need to implement
// to enable Git Town to search for proposals at the active forge.
type ProposalSearcher interface {
	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to search for a proposal that has the given branch as its source branch.
	// A None return value indicates that this connector does not support this feature (yet).
	SearchProposal(branch gitdomain.LocalBranchName) (Option[Proposal], error)
}

// ProposalMerger describes methods that connectors need to implement
// to enable Git Town to merge for proposals at the active forge.
type ProposalMerger interface {
	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to merge the proposal with the given number using the given message.
	// A None return value indicates that this connector does not support this feature (yet).
	SquashMergeProposal(number int, message gitdomain.CommitMessage) error
}

// ProposalUpdater describes methods that connectors need to implement
// to enable Git Town to update proposals at the active forge.
type ProposalBodyUpdater interface {
	// If the connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the body (description) of the given proposal to the given value.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalBody(proposal ProposalInterface, newBody string) error
}

// ProposalSourceUpdater describes methods that connectors need to implement
// to enable Git Town to update the source branch of proposals at the active forge.
type ProposalSourceUpdater interface {
	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the source branch of the proposal with the given number.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalSource(proposal ProposalInterface, newSource gitdomain.LocalBranchName) error
}

type ProposalTargetUpdater interface {
	// If this connector instance supports loading proposals via the API,
	// calling this function returns a function that you can call
	// to update the target branch of the proposal with the given number.
	// A None return value indicates that this connector does not support this feature (yet).
	UpdateProposalTarget(proposal ProposalInterface, newTarget gitdomain.LocalBranchName) error
}
