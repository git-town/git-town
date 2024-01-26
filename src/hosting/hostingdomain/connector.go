package hostingdomain

import "github.com/git-town/git-town/v11/src/git/gitdomain"

// Connector describes the activities that Git Town can perform on code hosting platforms.
// Individual implementations exist to talk to specific hosting platforms.
type Connector interface {
	// DefaultProposalMessage provides the text that the form for creating new proposals
	// on the respective hosting platform is prepopulated with.
	DefaultProposalMessage(proposal Proposal) string

	// FindProposal provides details about the proposal for the given branch into the given target branch.
	// Returns nil if no proposal exists.
	FindProposal(branch, target gitdomain.LocalBranchName) (*Proposal, error)

	// HostingPlatformName provides the name of the code hosting platform
	// supported by the respective connector implementation.
	HostingPlatformName() string

	// SquashMergeProposal squash-merges the proposal with the given number
	// using the given commit message.
	SquashMergeProposal(number int, message string) (mergeSHA gitdomain.SHA, err error)

	// NewProposalURL provides the URL of the page
	// to create a new proposal online.
	NewProposalURL(branch, parentBranch gitdomain.LocalBranchName) (string, error)

	// RepositoryURL provides the URL where the current repository can be found online.
	RepositoryURL() string

	// UpdateProposalTarget updates the target branch of the given proposal.
	UpdateProposalTarget(number int, target gitdomain.LocalBranchName) error
}
