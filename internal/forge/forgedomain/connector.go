package forgedomain

import (
	"github.com/git-town/git-town/v22/internal/git/gitdomain"
	"github.com/git-town/git-town/v22/internal/subshell/subshelldomain"
	. "github.com/git-town/git-town/v22/pkg/prelude"
)

// This file defines an interface for every capability that a forge and connector implementation can support.

// Connector describes the capabilities that all connectors can perform on forges.
type Connector interface {
	// BrowseRepository opens this repository in the associated application, typically the browser.
	BrowseRepository(runner subshelldomain.Runner) error

	// CreateProposal creates a proposal at the forge.
	CreateProposal(CreateProposalArgs) error

	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective forge type is prepopulated with.
	DefaultProposalMessage(proposal ProposalData) string
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

// ProposalBodyUpdater describes methods that connectors need to implement
// to enable Git Town to update proposals at the active forge.
type ProposalBodyUpdater interface {
	// Updates the body (description) of the given proposal to the given value.
	UpdateProposalBody(proposal ProposalInterface, newBody string) error
}

// ProposalFinder describes methods that connectors need to implement
// to enable Git Town to find proposals at the active forge.
type ProposalFinder interface {
	// Loads details about the proposal for the given branch into the given target branch.
	FindProposal(branch, target gitdomain.LocalBranchName) (Option[Proposal], error)
}

// ProposalMerger describes methods that connectors need to implement
// to enable Git Town to merge for proposals at the active forge.
type ProposalMerger interface {
	// Merges the proposal with the given number using the given message.
	SquashMergeProposal(number int, message gitdomain.CommitMessage) error
}

// ProposalSearcher describes methods that connectors need to implement
// to enable Git Town to search for proposals at the active forge.
type ProposalSearcher interface {
	// SearchProposals finds all active proposals that have the given branch as its source branch.
	SearchProposals(branch gitdomain.LocalBranchName) ([]Proposal, error)
}

// ProposalSourceUpdater describes methods that connectors need to implement
// to enable Git Town to update the source branch of proposals at the active forge.
type ProposalSourceUpdater interface {
	// Updates the source branch of the proposal with the given number.
	UpdateProposalSource(proposal ProposalInterface, newSource gitdomain.LocalBranchName) error
}

type ProposalTargetUpdater interface {
	// Updates the target branch of the proposal with the given number.
	UpdateProposalTarget(proposal ProposalInterface, newTarget gitdomain.LocalBranchName) error
}
